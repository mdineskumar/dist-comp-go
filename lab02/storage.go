package main

// ============================================================
// Lab 02 — Chord DHT
// File: storage.go
// Role: Distributed key-value storage on top of the ring
//
// TASKS IN THIS FILE:
//   Task 9  — Put()
//   Task 10 — Get()
//   Task 11 — Delete()
// ============================================================

import "fmt"

// ============================================================
// TASK 9 — Put
// ============================================================
// Store a key-value pair on the CORRECT node in the ring.
//
// Algorithm:
//   1. Hash the key:     id := hashKey(key)
//   2. Find the node responsible: succ, err := n.findSuccessor(id)
//   3. If succ is US:    call n.localPut(key, value)
//   4. If succ is REMOTE: forward via RPC
//        callRPC(succ.Addr, "ChordRPC.Put", &PutArgs{...}, &PutReply{})
//
// HINT: compare succ.Addr == n.info.Addr to check if it is us
//
// TODO: implement this function
func (n *Node) Put(key, value string) error {
	id := hashKey(key)
	succ, err := n.findSuccessor(id)

	if err != nil {
		return err
	}

	if succ.Addr == n.info.Addr {
		n.localPut(key, value)
	} else {
		callRPC(succ.Addr, "ChordRPC.Put", &PutArgs{Key: key, Value: value}, &PutReply{})
	}

	return nil
}

// ============================================================
// TASK 10 — Get
// ============================================================
// Retrieve a value from the CORRECT node in the ring.
//
// Algorithm:
//   1. Hash the key:     id := hashKey(key)
//   2. Find the node responsible: succ, err := n.findSuccessor(id)
//   3. If succ is US:    call n.localGet(key)
//                        return error if key not found
//   4. If succ is REMOTE: forward via RPC
//        callRPC(succ.Addr, "ChordRPC.Get", &GetArgs{...}, &GetReply{})
//        return reply.Value, or error if !reply.Found
//
// TODO: implement this function
func (n *Node) Get(key string) (string, error) {
	id := hashKey(key)

	succ, err := n.findSuccessor(id)

	if err != nil {
		return "", err
	}
	if succ.Addr == n.info.Addr {
		val, found := n.localGet(key)
		if !found {
			return "", fmt.Errorf("key %q not found", key)
		}
		return val, nil
	}

	var reply GetReply
	err = callRPC(succ.Addr, "ChordRPC.Get", &GetArgs{Key: key}, &reply)
	if err != nil {
		return "", err
	}

	if !reply.Found {
		return "", fmt.Errorf("key %q not found", key)
	}
	return reply.Value, nil

}

// ============================================================
// TASK 11 — Delete
// ============================================================
// Delete a key from the CORRECT node in the ring.
//
// Same routing logic as Put and Get:
//   1. Hash key → find successor
//   2. If us → localDelete
//   3. If remote → forward RPC
//
// TODO: implement this function
func (n *Node) Delete(key string) error {

	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// localPut stores a key-value pair on THIS node only
func (n *Node) localPut(key, value string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.storage[key] = value
	fmt.Printf("[STORAGE] Stored  key=%-20q  value=%q  (hash=%d  local)\n",
		key, value, hashKey(key))
}

// localGet retrieves a value from THIS node only
func (n *Node) localGet(key string) (string, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	val, ok := n.storage[key]
	return val, ok
}

// localDelete removes a key from THIS node only
func (n *Node) localDelete(key string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, ok := n.storage[key]
	if ok {
		delete(n.storage, key)
		fmt.Printf("[STORAGE] Deleted key=%-20q  (hash=%d  local)\n", key, hashKey(key))
	}
	return ok
}

// localKeys returns all keys stored on this node
func (n *Node) localKeys() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	keys := make([]string, 0, len(n.storage))
	for k := range n.storage {
		keys = append(keys, k)
	}
	return keys
}

// printStorage prints all locally stored key-value pairs
func (n *Node) printStorage() {
	n.mu.RLock()
	defer n.mu.RUnlock()
	fmt.Printf("\n── Storage on Node %d (%s) ──────────────────────\n",
		n.info.ID, n.info.Addr)
	if len(n.storage) == 0 {
		fmt.Println("  (empty)")
	} else {
		fmt.Printf("  %d keys:\n", len(n.storage))
		for k, v := range n.storage {
			fmt.Printf("  key=%-20q  value=%-20q  hash=%d\n", k, v, hashKey(k))
		}
	}
	fmt.Println()
}
