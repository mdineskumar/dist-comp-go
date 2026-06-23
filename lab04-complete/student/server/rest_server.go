package main

// ============================================================
// Lab 04 — RPC and Web Services
// File: rest_server.go  (REST)
// Role: REST HTTP server — implements 4 endpoints
//
// TASKS IN THIS FILE:
//   Task 9  — Set up HTTP router with all routes
//   Task 10 — PUT /keys/{key}
//   Task 11 — GET /keys/{key} and GET /keys
//   Task 12 — DELETE /keys/{key}
// ============================================================

// ── HOW REST WORKS ────────────────────────────────────────
//
// REST uses HTTP methods and URLs to define operations:
//
//   PUT    /keys/{key}    body: {"value":"..."}  → store a key
//   GET    /keys/{key}                           → get a key
//   GET    /keys                                 → list all keys
//   DELETE /keys/{key}                           → delete a key
//
// Unlike RPC (which hides the protocol), REST is visible:
// You can call it with curl, a browser, Postman, or any HTTP client
// in any language. This is why REST is the dominant style for
// public APIs and web services.
// ──────────────────────────────────────────────────────────

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RESTServer handles HTTP requests for the key-value store
type RESTServer struct {
	store *Store
}

// ============================================================
// TASK 9 — Set Up HTTP Router
// ============================================================
// Register URL patterns with their handler functions using
// http.HandleFunc on the given mux.
//
// Routes to register:
//   "/keys/"  → s.handleKey   (handles GET/PUT/DELETE /keys/{key})
//   "/keys"   → s.handleList  (handles GET /keys)
//
// HINT: Use mux.HandleFunc(pattern, handler)
//       The pattern "/keys/" (with trailing slash) matches
//       any path starting with /keys/
//
// TODO: implement
func (s *RESTServer) SetupRoutes(mux *http.ServeMux) {
	// YOUR CODE HERE
}

// ============================================================
// TASK 10 — PUT /keys/{key}
// ============================================================
// Handle PUT requests to store a key-value pair.
//
// Steps:
//   1. Extract key from URL: strings.TrimPrefix(r.URL.Path, "/keys/")
//   2. Decode JSON body into struct: { "value": "..." }
//      Use json.NewDecoder(r.Body).Decode(&body)
//   3. Call s.store.Put(key, body.Value)
//   4. Write response: w.WriteHeader(http.StatusCreated) (201)
//      json.NewEncoder(w).Encode(map[string]bool{"success": true})
//   5. Print: [REST] PUT key="..." value="..."
//
// TODO: implement in handleKey (see routing below)

// ============================================================
// TASK 11 — GET /keys/{key} and GET /keys
// ============================================================
// handleKey: GET /keys/{key} — retrieve one value
//   1. Extract key from URL
//   2. Call s.store.Get(key)
//   3. If found: 200 OK + json {"value": "...", "found": true}
//   4. If not found: 404 + json {"found": false}
//
// handleList: GET /keys — list all keys
//   1. Call s.store.List()
//   2. Return 200 OK + json {"keys": [...], "count": N}
//
// TODO: implement both

// ============================================================
// TASK 12 — DELETE /keys/{key}
// ============================================================
// Handle DELETE requests to remove a key.
//
// Steps:
//   1. Extract key from URL
//   2. Call s.store.Delete(key)
//   3. If deleted: 200 OK + json {"deleted": true}
//   4. If not found: 404 + json {"deleted": false}
//   5. Print: [REST] DELETE key="..."
//
// TODO: implement in handleKey

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// handleKey routes GET/PUT/DELETE /keys/{key} to the right handler
func (s *RESTServer) handleKey(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/keys/")
	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPut:
		s.handlePut(w, r, key)
	case http.MethodGet:
		s.handleGet(w, r, key)
	case http.MethodDelete:
		s.handleDelete(w, r, key)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleList handles GET /keys
func (s *RESTServer) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// TODO: implement in Task 11
	_ = fmt.Sprintf // remove when implementing
	json.NewEncoder(w).Encode(map[string]interface{}{
		"keys":  []string{},
		"count": 0,
	})
}

// handlePut stores a key — called from handleKey
func (s *RESTServer) handlePut(w http.ResponseWriter, r *http.Request, key string) {
	// TODO: implement in Task 10
}

// handleGet retrieves a key — called from handleKey
func (s *RESTServer) handleGet(w http.ResponseWriter, r *http.Request, key string) {
	// TODO: implement in Task 11
}

// handleDelete removes a key — called from handleKey
func (s *RESTServer) handleDelete(w http.ResponseWriter, r *http.Request, key string) {
	// TODO: implement in Task 12
}

// writeJSON is a helper to write a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
