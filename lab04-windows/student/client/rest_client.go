package main

// ============================================================
// Lab 04 — RPC and Web Services
// File: rest_client.go  (REST HTTP client)
// Role: Call REST API endpoints using net/http
//
// TASK IN THIS FILE:
//   Task 13 — Implement RESTClient
// ============================================================

// ── HOW REST CLIENT WORKS ─────────────────────────────────
//
// REST client uses standard HTTP requests:
//
//   PUT    http://server:8080/keys/mykey
//   Body:  {"value": "myvalue"}
//
//   GET    http://server:8080/keys/mykey
//   DELETE http://server:8080/keys/mykey
//   GET    http://server:8080/keys
//
// Unlike net/rpc and gRPC, REST needs no special client library.
// You can use curl, a browser, Postman, or net/http.
// This is why REST is the dominant style for public APIs.
//
// Test with curl (from your laptop terminal):
//   curl -X PUT http://localhost:8080/keys/city \
//        -H "Content-Type: application/json" \
//        -d '{"value":"London"}'
//
//   curl http://localhost:8080/keys/city
//   curl -X DELETE http://localhost:8080/keys/city
//   curl http://localhost:8080/keys
// ──────────────────────────────────────────────────────────

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RESTClient makes HTTP calls to the REST server
type RESTClient struct {
	baseURL string
	client  *http.Client
}

// ============================================================
// TASK 13 — Implement RESTClient
// ============================================================
//
// ── Put ───────────────────────────────────────────────────
// Make a PUT request to r.baseURL+"/keys/"+key
// Request body (JSON): {"value": value}
// Check response status == 201 Created
// Print: [REST] PUT key="..." value="..."
//
// HINT:
//
//	body, _ := json.Marshal(map[string]string{"value": value})
//	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
//	req.Header.Set("Content-Type", "application/json")
//	resp, err := r.client.Do(req)
//
// TODO: implement
func (r *RESTClient) Put(key, value string) error {
	url := r.baseURL + "/keys/" + key
	body, _ := json.Marshal(map[string]string{"value": value})

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("put failed: status %d", resp.StatusCode)
	}
	fmt.Printf("[REST] PUT key=%q value=%q\n", key, value)

	return nil
}

// ── Get ───────────────────────────────────────────────────
// Make a GET request to r.baseURL+"/keys/"+key
// If 200: decode JSON {"value": "...", "found": true}
// If 404: return "", false, nil
// Print: [REST] GET key="..."
//
// HINT:
//
//	resp, err := r.client.Get(url)
//	body, _ := io.ReadAll(resp.Body)
//	json.Unmarshal(body, &result)
//
// TODO: implement
func (r *RESTClient) Get(key string) (string, bool, error) {
	url := r.baseURL + "/keys/" + key

	resp, err := r.client.Get(url)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", false, nil
	}
	var result struct {
		Value string `json:"value"`
		Found bool   `json:"found"`
	}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	fmt.Printf("[REST] GET key=%v\n", key)
	return result.Value, result.Found, nil
}

// ── Delete ────────────────────────────────────────────────
// Make a DELETE request to r.baseURL+"/keys/"+key
// Decode JSON {"deleted": true/false}
// Print: [REST] DELETE key="..."
//
// TODO: implement
func (r *RESTClient) Delete(key string) (bool, error) {
	url := r.baseURL + "/keys/" + key

	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := r.client.Do(req)

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var result struct {
		Deleted bool `json:"deleted"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	fmt.Printf("[REST] DELETE key=%v\n", key)
	return result.Deleted, nil
}

// ── List ──────────────────────────────────────────────────
// Make a GET request to r.baseURL+"/keys"
// Decode JSON {"keys": [...], "count": N}
// Return the keys slice
// Print: [REST] LIST → N keys
//
// TODO: implement
func (r *RESTClient) List() ([]string, error) {
	url := r.baseURL + "/keys"
	resp, err := r.client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result struct {
		Keys []string `json:"keys"`
	}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	fmt.Printf("[REST] LIST -> %v keys\n", len(result.Keys))
	return result.Keys, nil
}

// NewRESTClient creates a REST client for the given base URL
func NewRESTClient(baseURL string) *RESTClient {
	return &RESTClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}
