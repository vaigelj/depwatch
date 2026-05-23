package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry holds a cached value with an expiry timestamp.
type Entry struct {
	Value     string    `json:"value"`
	FetchedAt time.Time `json:"fetched_at"`
	TTL       time.Duration `json:"ttl"`
}

// IsExpired returns true when the entry is older than its TTL.
func (e Entry) IsExpired() bool {
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache is a simple file-backed key/value store for version lookups.
type Cache struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
	defaultTTL time.Duration
}

// New loads (or creates) a cache file at the given path.
func New(path string, defaultTTL time.Duration) (*Cache, error) {
	c := &Cache{
		path:       path,
		entries:    make(map[string]Entry),
		defaultTTL: defaultTTL,
	}
	if err := c.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return c, nil
}

// Get returns a cached value and whether it exists and is still valid.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || e.IsExpired() {
		return "", false
	}
	return e.Value, true
}

// Set stores a value under key and persists the cache to disk.
func (c *Cache) Set(key, value string) error {
	c.mu.Lock()
	c.entries[key] = Entry{
		Value:     value,
		FetchedAt: time.Now(),
		TTL:       c.defaultTTL,
	}
	c.mu.Unlock()
	return c.save()
}

func (c *Cache) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Unmarshal(data, &c.entries)
}

func (c *Cache) save() error {
	c.mu.RLock()
	data, err := json.MarshalIndent(c.entries, "", "  ")
	c.mu.RUnlock()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o644)
}
