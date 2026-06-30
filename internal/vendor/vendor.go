// Package vendor downloads third-party front-end assets (currently KaTeX) at
// build time and caches them under the user cache directory. Nothing is checked
// into the repository: a default build links these assets from a CDN, and only
// `--vendor` pulls them down for a fully self-contained, offline/intranet site.
package vendor

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// KaTeX pinning. The npm tarball is content-addressed and immutable, so the
// sha256 below fully verifies what we download.
const (
	KaTeXVersion = "0.16.9"
	KaTeXCDN     = "https://cdn.jsdelivr.net/npm/katex@" + KaTeXVersion + "/dist/"

	katexTarballURL    = "https://registry.npmjs.org/katex/-/katex-" + KaTeXVersion + ".tgz"
	katexTarballSHA256 = "a0ccc017a6bff9e9dbc014b324e5203c52df16aeff105d971e916ce65b552733"
)

// File is a vendored asset destined for <out>/static/<RelPath>.
type File struct {
	RelPath string // e.g. "katex/katex.min.css", "katex/fonts/KaTeX_Main-Regular.woff2"
	Content []byte
}

// EnsureKaTeX returns the KaTeX assets, downloading and caching them on first
// use. When offline is true it serves only from cache and errors if absent.
func EnsureKaTeX(offline bool) ([]File, error) {
	dir := cacheDir()
	if dir != "" {
		if files, ok := loadCache(dir); ok {
			return files, nil
		}
	}
	if offline {
		return nil, fmt.Errorf("KaTeX %s is not cached and --offline was set; run once with network access", KaTeXVersion)
	}
	files, err := download()
	if err != nil {
		return nil, err
	}
	if dir != "" {
		_ = saveCache(dir, files) // a cache miss is not fatal
	}
	return files, nil
}

func cacheDir() string {
	base, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(base, "starbase", "vendor", "katex-"+KaTeXVersion)
}

func loadCache(dir string) ([]File, bool) {
	if _, err := os.Stat(filepath.Join(dir, ".complete")); err != nil {
		return nil, false
	}
	var files []File
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		if rel == ".complete" {
			return nil
		}
		b, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		files = append(files, File{RelPath: filepath.ToSlash(rel), Content: b})
		return nil
	})
	if err != nil || len(files) == 0 {
		return nil, false
	}
	return files, true
}

func saveCache(dir string, files []File) error {
	for _, f := range files {
		dst := filepath.Join(dir, filepath.FromSlash(f.RelPath))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, f.Content, 0o644); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(dir, ".complete"), []byte(KaTeXVersion), 0o644)
}

// download fetches and verifies the KaTeX npm tarball, then extracts the CSS,
// JS, woff2 fonts, and license into File entries under "katex/".
func download() ([]File, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	req, _ := http.NewRequest("GET", katexTarballURL, nil)
	req.Header.Set("User-Agent", "starbase-vendor")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading KaTeX: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading KaTeX: HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, fmt.Errorf("reading KaTeX tarball: %w", err)
	}
	if sum := hex.EncodeToString(hashBytes(raw)); sum != katexTarballSHA256 {
		return nil, fmt.Errorf("KaTeX tarball checksum mismatch: got %s, want %s", sum, katexTarballSHA256)
	}

	gz, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	var files []File
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// npm tarballs root everything under "package/".
		name := strings.TrimPrefix(hdr.Name, "package/")
		rel, keep := katexDest(name)
		if !keep {
			continue
		}
		b, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		files = append(files, File{RelPath: rel, Content: b})
	}
	if len(files) < 3 {
		return nil, fmt.Errorf("KaTeX tarball missing expected files (got %d)", len(files))
	}
	return files, nil
}

// katexDest maps a path inside the tarball to its output path under "katex/",
// keeping only the files we actually serve (css, js, woff2 fonts, license).
func katexDest(name string) (string, bool) {
	switch {
	case name == "dist/katex.min.css":
		return "katex/katex.min.css", true
	case name == "dist/katex.min.js":
		return "katex/katex.min.js", true
	case strings.HasPrefix(name, "dist/fonts/") && strings.HasSuffix(name, ".woff2"):
		return "katex/fonts/" + path.Base(name), true
	case name == "LICENSE":
		return "katex/LICENSE", true
	}
	return "", false
}

func hashBytes(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}
