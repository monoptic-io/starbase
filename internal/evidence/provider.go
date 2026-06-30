package evidence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

// inputSpec is one line of a check's `inputs` manifest: a source resolved by a
// provider and staged into the check's working directory under LocalName.
type inputSpec struct {
	Source    string
	LocalName string
}

// resolved is an input materialized to local bytes: a content hash for the cache
// key, and realize() to place the bytes into the staging dir.
type resolved struct {
	LocalName string
	Hash      string
	realize   func(dest string) error
}

// provider resolves a source to local bytes. The file and http(s) providers
// ship; register more by scheme in providerFor (s3, db, …).
type provider interface {
	resolve(spec inputSpec, contentDir string, force bool) (resolved, error)
}

func providerFor(source string) provider {
	if u, err := url.Parse(source); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		return httpProvider{}
	}
	return fileProvider{}
}

func defaultName(source string) string {
	if u, err := url.Parse(source); err == nil && u.Scheme != "" {
		if b := path.Base(u.Path); b != "" && b != "/" && b != "." {
			return b
		}
		return "input"
	}
	return filepath.Base(source)
}

func nameOr(localName, source string) string {
	if localName != "" {
		return localName
	}
	return defaultName(source)
}

// --- file provider ---

type fileProvider struct{}

func (fileProvider) resolve(spec inputSpec, contentDir string, _ bool) (resolved, error) {
	src := spec.Source
	if !filepath.IsAbs(src) {
		src = filepath.Join(contentDir, src)
	}
	h, err := hashFile(src)
	if err != nil {
		return resolved{}, err
	}
	return resolved{
		LocalName: nameOr(spec.LocalName, spec.Source),
		Hash:      h,
		realize:   func(dest string) error { return linkOrCopy(src, dest) },
	}, nil
}

// --- http provider ---
//
// Fetches the URL into a content cache (user cache dir) and treats it as
// immutable between builds: a local verify reuses the cached bytes rather than
// re-fetching. CI starts cold, fetches fresh, and the output comparison — not
// the cache — is what catches a URL whose content has drifted.

type httpProvider struct{}

func (httpProvider) resolve(spec inputSpec, _ string, force bool) (resolved, error) {
	store := inputCachePath(spec.Source)
	if force || !regularFile(store) {
		if err := fetch(spec.Source, store); err != nil {
			return resolved{}, err
		}
	}
	h, err := hashFile(store)
	if err != nil {
		return resolved{}, err
	}
	return resolved{
		LocalName: nameOr(spec.LocalName, spec.Source),
		Hash:      h,
		realize:   func(dest string) error { return linkOrCopy(store, dest) },
	}, nil
}

func fetch(rawurl, dest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawurl, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", rawurl, resp.Status)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp := dest + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, dest)
}

func inputCachePath(source string) string {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}
	sum := sha256.Sum256([]byte(source))
	return filepath.Join(base, "starbase", "inputs", hex.EncodeToString(sum[:])[:24])
}

// --- shared helpers ---

func hashFile(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func regularFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.Mode().IsRegular()
}

func linkOrCopy(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := os.Symlink(src, dest); err == nil {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
