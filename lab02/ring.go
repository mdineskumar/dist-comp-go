package main

// ============================================================
// Lab 02 — Chord DHT
// File: ring.go
// Role: Ring operations — lookup, join, maintenance
//
// TASKS IN THIS FILE:
//   Task 3 — findSuccessor()
//   Task 4 — closestPrecedingFinger()
//   Task 5 — join()
//   Task 6 — stabilize()
//   Task 7 — fixFingers()
// ============================================================

import (
	"fmt"
	"log"
	"time"
)

// ============================================================
// TASK 3 — findSuccessor
// ============================================================
// Find which node is responsible for the given id.
// This is the core Chord lookup algorithm.
//
// Algorithm:
//  1. Check if id falls in (n.info.ID, n.successor().ID]
//     → If yes: our successor owns it, return successor
//  2. Otherwise: forward to the closest preceding finger
//     → Call FindSuccessor RPC on that finger node
//     → Return its answer
//
// Special case: if closestPrecedingFinger returns this node itself,
// return this node (avoids infinite loop)
//
// HINT: use inRange(id, n.info.ID, n.successor().ID)
//
//	use closestPrecedingFinger(id)
//	use callRPC(addr, "ChordRPC.FindSuccessor", args, reply)
//
// TODO: implement this function
func (n *Node) findSuccessor(id uint8) (NodeInfo, error) {
	succ := n.successor()

	// 1. Check if id falls in (n.info.ID, n.successor().ID]
	// inRange(id, start, end) exactly maps to the half-open arc (start, end]
	if inRange(id, n.info.ID, n.successor().ID) {
		return succ, nil
	}

	cpNode := n.closestPrecedingFinger(id)

	//there is only one node
	if cpNode.ID == n.info.ID {
		return n.info, nil
	}

	var reply FindSuccessorReply
	err := callRPC(cpNode.Addr, "ChordRPC.FindSuccessor", &FindSuccessorArgs{ID: id}, &reply)

	if err != nil {
		return NodeInfo{}, err
	}

	//return reply node
	return reply.Node, nil
}

// ============================================================
// TASK 4 — closestPrecedingFinger
// ============================================================
// Search the finger table for the closest node preceding id.
//
// Algorithm:
//
//	Walk fingers from M-1 DOWN to 0.
//	Return the first finger whose ID falls in (n.info.ID, id).
//	If no finger qualifies, return this node itself.
//
// HINT: use inRange(finger.ID, n.info.ID, id)
//
//	remember to acquire read lock: n.mu.RLock() / n.mu.RUnlock()
//
// TODO: implement this function
func (n *Node) closestPrecedingFinger(id uint8) NodeInfo {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for i := M - 1; i >= 0; i-- {
		f := n.fingers[i]
		//return first finger whose ID falss in (n.info.ID, id)
		//inRange is (start, end], strictly find preceding finger
		if inRange(f.ID, n.info.ID, id) && f.ID != id { //&& f.ID != id && f.Addr != n.info.Addr {
			return f
		}
	}
	// if no finder qualifies, return this node itself
	return n.info
}

// ============================================================
// TASK 5 — join
// ============================================================
// Join the Chord ring.
//
// Case A — existingAddr is empty: CREATE a new ring
//
//	→ Set predecessor to nil
//	→ Set all fingers to point to self (already done by newNode)
//	→ Done — this node is alone on the ring
//
// Case B — existingAddr is provided: JOIN an existing ring
//
//	→ Call FindSuccessor RPC on existingAddr to find OUR successor
//	   (pass n.info.ID as the ID to look up)
//	→ Set fingers[0] to the returned successor
//	→ Set predecessor to nil (stabilize will discover it later)
//
// HINT: use callRPC(existingAddr, "ChordRPC.FindSuccessor", args, reply)
//
//	use FindSuccessorArgs{ID: n.info.ID} as the args
//	use n.mu.Lock() / n.mu.Unlock() when updating fingers
//
// TODO: implement this function
func (n *Node) join(existingAddr string) error {
	//Case A — existingAddr is empty: CREATE a new ring
	if existingAddr == "" {
		n.mu.Lock()
		n.predecessor = nil
		for i := range n.fingers {
			n.fingers[i] = n.info
		}

		n.mu.Unlock()

		return nil
	}

	//Case B: join existing ring
	var reply FindSuccessorReply
	if err := callRPC(existingAddr, "ChordRPC.FindSuccessor", &FindSuccessorArgs{ID: n.info.ID}, &reply); err != nil {
		return err
	}

	n.mu.Lock()
	n.predecessor = nil
	n.fingers[0] = reply.Node
	n.mu.Unlock()

	return nil
}

// ============================================================
// TASK 6 — stabilize
// ============================================================
// Periodically called to keep successor pointers correct.
//
// Algorithm:
//  1. Ask our successor for its predecessor
//     → Call GetPredecessor RPC on successor
//  2. If that predecessor (call it x) exists AND
//     x.ID falls in (n.info.ID, n.successor().ID)
//     → x is a better successor for us — update fingers[0] = x
//  3. Notify our (possibly updated) successor that we exist
//     → Call Notify RPC on successor with our NodeInfo
//
// If any RPC fails, just return (network issues are temporary)
//
// HINT: use GetPredecessorArgs{} and GetPredecessorReply{}
//
//	use NotifyArgs{Node: n.info} and NotifyReply{}
//
// TODO: implement this function
func (n *Node) stabilize() {
	succ := n.successor()
	var reply GetPredecessorReply
	err := callRPC(succ.Addr, "ChordRPC.GetPredecessor", &GetPredecessorArgs{}, &reply)
	if err != nil {
		return
	}

	x := reply.Node

	if x != nil && inRange(x.ID, n.info.ID, succ.ID) {
		n.mu.Lock()
		n.fingers[0] = *x
		n.mu.Unlock()
		succ = *x

	}

	var notifyReply NotifyReply

	errNotify := callRPC(succ.Addr, "ChordRPC.Notify", &NotifyArgs{Node: n.info}, &notifyReply)

	if errNotify != nil {
		return
	}
}

// ============================================================
// TASK 7 — fixFingers
// ============================================================
// Periodically called to refresh one finger table entry.
// Rotates through fingers one at a time using n.next.
//
// Algorithm:
//  1. Read n.next (the current finger index to fix)
//  2. Calculate: start = (n.info.ID + 2^next) mod RING_SIZE
//  3. Call findSuccessor(start) to get the correct node for that finger
//  4. Update fingers[next] with the result
//  5. Increment n.next, wrap to 0 after reaching M
//
// HINT: 2^i in Go: uint8(1) << uint(i)
//
//	wrap: n.next = (n.next + 1) % M
//	use n.mu.Lock() when reading/writing n.next and n.fingers
//
// TODO: implement this function
func (n *Node) fixFingers() {
	n.mu.Lock()
	next := n.next
	n.mu.Unlock()

	start := (n.info.ID + uint8(1)<<uint(next))

	fingerNode, err := n.findSuccessor(start)

	n.mu.Lock()
	defer n.mu.Unlock()

	if err != nil {
		return
	}

	n.fingers[next] = fingerNode
	n.next = (n.next + 1) % M

}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// checkPredecessor verifies our predecessor is still alive.
func (n *Node) checkPredecessor() {
	n.mu.RLock()
	pred := n.predecessor
	n.mu.RUnlock()

	if pred == nil {
		return
	}
	err := callRPC(pred.Addr, "ChordRPC.Ping", &PingArgs{}, &PingReply{})
	if err != nil {
		log.Printf("[RING] Predecessor ID:%d appears to have failed — clearing", pred.ID)
		n.mu.Lock()
		n.predecessor = nil
		n.mu.Unlock()
	}
}

// startBackgroundTasks launches periodic maintenance goroutines.
func (n *Node) startBackgroundTasks() {
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			n.stabilize()
		}
	}()
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			n.fixFingers()
		}
	}()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			n.checkPredecessor()
		}
	}()
	fmt.Println("[RING] Background maintenance started (stabilize, fixFingers, checkPredecessor)")
}

// printRing shows this node's ring position and finger table
func (n *Node) printRing() {
	n.mu.RLock()
	defer n.mu.RUnlock()

	fmt.Printf("\n── Ring Position ────────────────────────────────\n")
	fmt.Printf("  This node : ID=%-3d  Addr=%s\n", n.info.ID, n.info.Addr)
	fmt.Printf("  Successor : ID=%-3d  Addr=%s\n", n.fingers[0].ID, n.fingers[0].Addr)
	if n.predecessor != nil {
		fmt.Printf("  Predecessor: ID=%-3d  Addr=%s\n", n.predecessor.ID, n.predecessor.Addr)
	} else {
		fmt.Printf("  Predecessor: (unknown — stabilize not run yet)\n")
	}
	fmt.Printf("\n── Finger Table ─────────────────────────────────\n")
	for i, f := range n.fingers {
		start := n.info.ID + (uint8(1) << uint(i))
		fmt.Printf("  finger[%d]  start=%-3d → ID=%-3d  Addr=%s\n",
			i, start, f.ID, f.Addr)
	}
	fmt.Println()
}
