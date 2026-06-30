// Package cache persists per-output fingerprints so a build can skip pages
// whose inputs haven't changed.
package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const version = 2

type Cache struct {
	Version int               `json:"version"`
	Pages   map[string]string `json:"pages"`  // output path -> fingerprint
	Assets  map[string]string `json:"assets"` // asset path -> content hash
}

func New() *Cache {
	return &Cache{Version: version, Pages: map[string]string{}, Assets: map[string]string{}}
}

// Load reads the cache from dir, returning a fresh one if absent or stale.
func Load(dir string) *Cache {
	b, err := os.ReadFile(filepath.Join(dir, ".starbase-cache.json"))
	if err != nil {
		return New()
	}
	var c Cache
	if json.Unmarshal(b, &c) != nil || c.Version != version || c.Pages == nil {
		return New()
	}
	if c.Assets == nil {
		c.Assets = map[string]string{}
	}
	return &c
}

// Save writes the cache to dir.
func (c *Cache) Save(dir string) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ".starbase-cache.json"), b, 0o644)
}

// PageFresh reports whether the page's fingerprint matches the cached one.
func (c *Cache) PageFresh(out, fingerprint string) bool {
	return c.Pages[out] == fingerprint
}

func (c *Cache) PutPage(out, fingerprint string) { c.Pages[out] = fingerprint }

func (c *Cache) AssetFresh(name, hash string) bool { return c.Assets[name] == hash }
func (c *Cache) PutAsset(name, hash string)        { c.Assets[name] = hash }
