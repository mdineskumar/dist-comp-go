package main

// ============================================================
// Lab 06 — Caching and Replication
// File: stats.go
// Role: Cache performance statistics
//
// TASKS IN THIS FILE:
//   Task 12 — Hit(), Miss(), Eviction() counters
//   Task 13 — HitRate() and PrintStats()
// ============================================================

import (
	"fmt"
	"sync/atomic" // used in solution
)

// Stats tracks cache performance counters.
// Uses atomic operations so counters are safe to update
// from multiple goroutines without a mutex.
type Stats struct {
	hits      int64
	misses    int64
	evictions int64
}

// NewStats creates a zeroed Stats struct
func NewStats() *Stats {
	return &Stats{}
}

// ============================================================
// TASK 12 — Hit, Miss, Eviction counters
// ============================================================
// Increment the appropriate counter using atomic operations.
//
// Use: atomic.AddInt64(&s.hits, 1)
//
//	atomic.AddInt64(&s.misses, 1)
//	atomic.AddInt64(&s.evictions, 1)
//
// TODO: implement all three functions
func (s *Stats) Hit() {
	atomic.AddInt64(&s.hits, 1)
}

func (s *Stats) Miss() {
	atomic.AddInt64(&s.misses, 1)
}

func (s *Stats) Eviction() {
	atomic.AddInt64(&s.evictions, 1)
}

// ============================================================
// TASK 13 — HitRate and PrintStats
// ============================================================
//
// ── HitRate ───────────────────────────────────────────────
// Return the cache hit rate as a float64 (0.0 to 1.0).
// Formula: hits / (hits + misses)
// If hits + misses == 0, return 0.0 (avoid division by zero)
//
// Use atomic.LoadInt64(&s.hits) to safely read the counters
//
// TODO: implement
func (s *Stats) HitRate() float64 {
	hits := atomic.LoadInt64(&s.hits)
	misses := atomic.LoadInt64(&s.misses)

	if (hits + misses) == 0 {
		return 0.0
	}
	return float64(hits) / (float64(hits) + float64(misses))
}

// ── PrintStats ────────────────────────────────────────────
// Print a formatted stats summary. Example output:
//
//	── Cache Statistics ─────────────────────────────
//	  Hits:      42
//	  Misses:    8
//	  Evictions: 3
//	  Hit Rate:  84.0%
//
// TODO: implement
func (s *Stats) PrintStats() {
	fmt.Println("---Cache Statistics--------------")
	fmt.Printf("Hits: %v\n", atomic.LoadInt64(&s.hits))
	fmt.Printf("Misses: %v\n", atomic.LoadInt64(&s.misses))
	fmt.Printf("Evictions: %v\n", atomic.LoadInt64(&s.evictions))
	fmt.Printf("Hit Rate: %v\n", s.HitRate())
}
