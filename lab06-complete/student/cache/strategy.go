package main

// ============================================================
// Lab 06 — Caching and Replication
// File: strategy.go
// Role: Write strategies — write-through and write-back
//
// TASKS IN THIS FILE:
//   Task 5 — WriteThrough()
//   Task 6 — WriteBack()
//   Task 7 — flushDirty()
// ============================================================

import (
	"fmt"
	"time"
)

// ============================================================
// TASK 5 — WriteThrough
// ============================================================
// Write a key-value pair to BOTH the cache and the origin
// simultaneously. The write is only confirmed when BOTH succeed.
//
// Steps:
//  1. Write to origin first: originClient.Set(key, value)
//     If this fails, return the error (cache stays unchanged)
//  2. Write to cache: c.Set(key, value, false)
//     dirty=false because cache and origin are now in sync
//  3. Print: [WRITE-THROUGH] key="..." written to origin + cache
//  4. Return nil
//
// ── WRITE-THROUGH TRADE-OFFS ──────────────────────────────
// Pro: Cache and origin always consistent — no stale reads
// Con: Every write waits for the origin (slower)
// Use when: data consistency is critical (e.g. bank balances)
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func WriteThrough(c *Cache, originClient *OriginClient, key, value string) error {
	ok := originClient.Set(key, value)
	if !ok {
		return fmt.Errorf("write to origin failed\n")
	}
	c.Set(key, value, false)

	fmt.Printf("[WRITE-THROUGH] key=%v written to origin + cache", key)

	return nil
}

// ============================================================
// TASK 6 — WriteBack
// ============================================================
// Write a key-value pair to the CACHE ONLY immediately.
// The origin is updated later in the background by flushDirty().
//
// Steps:
//  1. Write to cache: c.Set(key, value, true)
//     dirty=true because origin has NOT been updated yet
//  2. Print: [WRITE-BACK] key="..." written to cache (pending origin flush)
//  3. Return nil (always succeeds immediately)
//
// ── WRITE-BACK TRADE-OFFS ────────────────────────────────
// Pro: Writes are fast — no waiting for origin
// Con: If cache crashes before flush, writes are lost forever
// Use when: speed matters more than perfect durability
//
//	(e.g. analytics counters, session data)
//
// ──────────────────────────────────────────────────────────
//
// TODO: implement this function
func WriteBack(c *Cache, key, value string) error {
	c.Set(key, value, true)
	fmt.Printf("[WRITE-BACK] key=%v written to cache(pending origin flush)\n", key)
	return nil
}

// ============================================================
// TASK 7 — flushDirty
// ============================================================
// Background function that periodically writes dirty cache
// entries to the origin. Called by startFlushWorker() below.
//
// Steps:
//  1. Get all dirty keys: keys := c.DirtyKeys()
//  2. For each dirty key:
//     a. Get the value: value, ok := c.Get(key)
//     b. If not found (was evicted), skip it
//     c. Write to origin: originClient.Set(key, value)
//     d. If write succeeds: c.MarkClean(key)
//     e. Print: [FLUSH] key="..." flushed to origin
//  3. If len(keys) > 0: print total flushed count
//
// TODO: implement this function
func flushDirty(c *Cache, originClient *OriginClient) {
	keys := c.DirtyKeys()
	for _, key := range keys {
		value, ok := c.Get(key)
		if !ok {
			continue
		}
		success := originClient.Set(key, value)
		if success {
			c.MarkClean(key)
		}
		fmt.Printf("[FLUSH] key=%v flushed to origin\n", key)
	}

	if len(keys) > 0 {
		fmt.Printf("[FLUSH] total flushed count: %v\n", len(keys))
	}

}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// startFlushWorker launches a background goroutine that calls
// flushDirty() every flushInterval seconds
func startFlushWorker(c *Cache, originClient *OriginClient, flushInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(flushInterval)
		for range ticker.C {
			flushDirty(c, originClient)
		}
	}()
	fmt.Printf("[FLUSH] Background flush worker started (interval: %v)\n", flushInterval)
}
