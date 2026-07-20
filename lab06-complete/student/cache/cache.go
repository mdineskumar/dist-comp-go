package main

// ============================================================
// Lab 06 — Caching and Replication
// File: cache.go
// Role: Core LRU cache data structure
//
// TASKS IN THIS FILE:
//   Task 1 — NewCache()
//   Task 2 — Get()
//   Task 3 — Set()
//   Task 4 — evictLRU()
// ============================================================

import (
	"fmt"
	"sync"
	"time"
)

// Entry represents one cached key-value pair with metadata
type Entry struct {
	Value        string
	LastAccessed time.Time // updated every time this entry is read or written
	IsDirty      bool      // true = written to cache but NOT yet flushed to origin (write-back)
}

// Cache is a thread-safe LRU cache with a fixed capacity.
// When the cache is full and a new key arrives, the Least
// Recently Used entry is evicted to make room.
type Cache struct {
	mu       sync.Mutex
	items    map[string]Entry // key → entry
	capacity int              // maximum number of keys
	stats    *Stats           // hit/miss/eviction counters
}

// ============================================================
// TASK 1 — NewCache
// ============================================================
// Create and return a new Cache with the given capacity.
//
// Steps:
//  1. Create a Cache with:
//     items:    make(map[string]Entry)
//     capacity: capacity
//     stats:    NewStats()
//  2. Print: [CACHE] Initialised (capacity=N)
//  3. Return a pointer to the cache
//
// TODO: implement this function
func NewCache(capacity int) *Cache {
	cache := Cache{items: make(map[string]Entry), capacity: capacity, stats: NewStats()}
	fmt.Printf("[CACHE] Initialised (capacity=%v)\n", capacity)
	return &cache
}

// ============================================================
// TASK 2 — Get
// ============================================================
// Retrieve a value from the cache.
//
// Steps:
//  1. Lock: c.mu.Lock() / defer c.mu.Unlock()
//  2. Look up the key in c.items
//  3. If NOT found:
//     c.stats.Miss()
//     return "", false
//  4. If found:
//     Update LastAccessed: entry.LastAccessed = time.Now()
//     Store back: c.items[key] = entry
//     c.stats.Hit()
//     return entry.Value, true
//
// ── WHY UPDATE LastAccessed? ──────────────────────────────
// LRU = Least Recently Used. Every time we read a key, we mark
// it as "recently used" so it stays in the cache longer.
// Only keys that nobody reads for a long time get evicted.
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func (c *Cache) Get(key string) (string, bool) {
	// YOUR CODE HERE
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.items[key]
	if !ok {
		c.stats.Miss()
		return "", false
	}
	entry.LastAccessed = time.Now()
	c.items[key] = entry
	c.stats.Hit()
	return entry.Value, true
}

// ============================================================
// TASK 3 — Set
// ============================================================
// Store a key-value pair in the cache.
//
// Steps:
//  1. Lock: c.mu.Lock() / defer c.mu.Unlock()
//  2. If the key does NOT already exist AND cache is full:
//     call c.evictLRU()   (removes the least recently used entry)
//  3. Store the entry:
//     c.items[key] = Entry{
//     Value:        value,
//     LastAccessed: time.Now(),
//     IsDirty:      dirty,  (true for write-back, false for write-through)
//     }
//  4. Print: [CACHE] Set key="..." value="..." dirty=true/false
//
// TODO: implement this function
func (c *Cache) Set(key, value string, dirty bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[key]

	if !found && c.IsFull() {
		c.evictLRU()
	}

	c.items[key] = Entry{
		Value:        value,
		LastAccessed: time.Now(),
		IsDirty:      dirty,
	}
	fmt.Printf("[CACHE] Set key=%v value=%v dirty=%v\n", key, value, dirty)
}

// ============================================================
// TASK 4 — evictLRU
// ============================================================
// Remove the least recently used entry from the cache.
// Called by Set() when the cache is full.
//
// Algorithm:
//  1. Loop through all entries in c.items
//  2. Find the entry with the OLDEST LastAccessed time
//     (smallest time.Time value)
//  3. Delete that key from c.items
//  4. c.stats.Eviction()
//  5. Print: [CACHE] Evicted key="..." (LRU)
//
// HINT: Use a variable to track the oldest key seen so far:
//
//	oldestKey := ""
//	var oldestTime time.Time
//	for key, entry := range c.items {
//	    if oldestKey == "" || entry.LastAccessed.Before(oldestTime) {
//	        oldestKey = key
//	        oldestTime = entry.LastAccessed
//	    }
//	}
//
// NOTE: evictLRU is called from Set() which already holds the lock —
// do NOT lock again here (would deadlock)
//
// TODO: implement this function
func (c *Cache) evictLRU() {
	oldestKey := ""
	var oldestTime time.Time

	for key, entry := range c.items {
		if oldestKey == "" || entry.LastAccessed.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.LastAccessed
		}
	}
	delete(c.items, oldestKey)

	c.stats.Eviction()
	fmt.Printf("[CACHE] Evicted key=%v (LRU)\n", oldestKey)
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// Size returns the number of items currently in the cache
func (c *Cache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

// IsFull returns true if the cache has reached capacity
func (c *Cache) IsFull() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items) >= c.capacity
}

// DirtyKeys returns all keys that are dirty (written to cache but not origin)
func (c *Cache) DirtyKeys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var keys []string
	for k, e := range c.items {
		if e.IsDirty {
			keys = append(keys, k)
		}
	}
	return keys
}

// MarkClean marks a key as clean (synced to origin)
func (c *Cache) MarkClean(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.items[key]; ok {
		entry.IsDirty = false
		c.items[key] = entry
	}
}

// PrintContents shows all cached key-value pairs
func (c *Cache) PrintContents() {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Printf("\n── Cache Contents (%d/%d) ──────────────────────\n",
		len(c.items), c.capacity)
	if len(c.items) == 0 {
		fmt.Println("  (empty)")
	}
	for k, e := range c.items {
		dirty := ""
		if e.IsDirty {
			dirty = " [dirty]"
		}
		fmt.Printf("  key=%-20q  value=%-20q  last_used=%s%s\n",
			k, e.Value, e.LastAccessed.Format("15:04:05.000"), dirty)
	}
	fmt.Println()
}
