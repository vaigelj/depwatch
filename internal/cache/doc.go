// Package cache provides a lightweight, file-backed key/value cache used by
// depwatch to store resolved latest-version lookups between runs.
//
// Entries are stored as JSON on disk and are considered valid until their
// configured TTL elapses. A zero or negative TTL causes entries to expire
// immediately, effectively disabling caching.
//
// Typical usage:
//
//	c, err := cache.New("/tmp/depwatch/cache.json", 24*time.Hour)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if v, ok := c.Get("npm:lodash"); ok {
//		// use cached version v
//	} else {
//		v = fetchLatest("lodash")
//		_ = c.Set("npm:lodash", v)
//	}
package cache
