package main

// ============================================================
// Lab 07 — Consensus (Simplified Raft)
// File: rpc.go
// Role: RPC server and the two core Raft RPC handlers
//
// TASKS IN THIS FILE:
//   Task 6  — RequestVote handler
//   Task 10 — AppendEntries handler
// ============================================================

import (
	"fmt"
	"net"
	"net/rpc"
)

// ── RPC argument / reply types ────────────────────────────

// RequestVote — sent by Candidates to gather votes
type RequestVoteArgs struct {
	Term         int    // candidate's current term
	CandidateID  string // candidate requesting vote
	LastLogIndex int    // index of candidate's last log entry
	LastLogTerm  int    // term of candidate's last log entry
}
type RequestVoteReply struct {
	Term        int  // current term of this node (for candidate to update itself)
	VoteGranted bool // true if candidate received vote
}

// AppendEntries — sent by Leader for heartbeats AND log replication
type AppendEntriesArgs struct {
	Term         int        // leader's current term
	LeaderID     string     // so follower can redirect clients
	PrevLogIndex int        // index of log entry immediately before new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	Entries      []LogEntry // log entries to store (empty for heartbeat)
	LeaderCommit int        // leader's commitIndex
}
type AppendEntriesReply struct {
	Term    int  // current term (for leader to update itself)
	Success bool // true if follower contained entry matching prevLogIndex/prevLogTerm
}

// Client RPC types
type SubmitArgs  struct{ Command string }
type SubmitReply struct{ Success bool; Leader string; Err string }
type GetValueArgs  struct{ Key string }
type GetValueReply struct{ Value string; Found bool }
type StatusArgs  struct{}
type StatusReply struct{ ID, State, Leader string; Term, LogLen, Committed int }

// RaftRPC is the RPC handler
type RaftRPC struct{ node *Node }

// ============================================================
// TASK 6 — RequestVote Handler
// ============================================================
// Called by a Candidate asking for our vote.
// We grant our vote ONLY if ALL conditions are true:
//   1. args.Term >= n.currentTerm  (candidate not behind us)
//   2. We haven't voted yet this term OR we already voted for this candidate
//   3. Candidate's log is at least as up-to-date as ours:
//        args.LastLogTerm > n.lastLogTerm()
//        OR (args.LastLogTerm == n.lastLogTerm() AND args.LastLogIndex >= n.lastLogIndex())
//
// Steps:
//   1. Lock: n.mu.Lock() / defer n.mu.Unlock()
//   2. If args.Term > n.currentTerm → n.becomeFollower(args.Term)
//   3. Set reply.Term = n.currentTerm
//   4. Check the 3 conditions above
//   5. If all conditions met:
//        n.votedFor = args.CandidateID
//        reply.VoteGranted = true
//        n.resetElectionTimer()  (we voted — reset our own timer)
//        Print: [NODE id] Voted for CandidateID in term N
//   6. Otherwise: reply.VoteGranted = false
//
// TODO: implement this function
func (r *RaftRPC) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// TASK 10 — AppendEntries Handler
// ============================================================
// Called by the Leader for heartbeats AND log replication.
// This is the most complex RPC in Raft.
//
// Steps:
//   1. Lock: n.mu.Lock() / defer n.mu.Unlock()
//   2. Set reply.Term = n.currentTerm
//   3. If args.Term < n.currentTerm: reply.Success = false; return
//      (stale leader — reject)
//   4. If args.Term > n.currentTerm: n.becomeFollower(args.Term)
//   5. n.resetElectionTimer()  (we heard from a valid leader — reset)
//   6. Check log consistency:
//        If args.PrevLogIndex > 0:
//          If we don't have that index OR term doesn't match:
//            reply.Success = false; return
//   7. Append any new entries (skip entries we already have):
//        For each entry in args.Entries:
//          If entry.Index > len(n.log): append it
//          Else if n.log[entry.Index-1].Term != entry.Term:
//            truncate log to entry.Index-1, append entry
//   8. Update commitIndex:
//        if args.LeaderCommit > n.commitIndex:
//          n.commitIndex = min(args.LeaderCommit, n.lastLogIndex())
//          go n.applyCommitted()
//   9. reply.Success = true
//   10. Print if entries received: [NODE id] Appended N entries (term T)
//
// TODO: implement this function
func (r *RaftRPC) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	// YOUR CODE HERE
	return nil
}

// ============================================================
// Below this line — already implemented, do not change
// ============================================================

// Submit — client sends a command to this node
// If this node is the leader: append to log and replicate
// If not the leader: tell client who the leader is
func (r *RaftRPC) Submit(args *SubmitArgs, reply *SubmitReply) error {
	r.node.mu.Lock()
	if r.node.state != Leader {
		leader := r.node.votedFor
		r.node.mu.Unlock()
		reply.Success = false
		reply.Leader = leader
		reply.Err = "not the leader"
		return nil
	}
	r.node.mu.Unlock()

	success, err := r.node.AppendEntry(args.Command)
	reply.Success = success
	if err != nil {
		reply.Err = err.Error()
	}
	return nil
}

// GetValue — read from state machine
func (r *RaftRPC) GetValue(args *GetValueArgs, reply *GetValueReply) error {
	reply.Value, reply.Found = r.node.GetValue(args.Key)
	return nil
}

// Status — return node status
func (r *RaftRPC) Status(args *StatusArgs, reply *StatusReply) error {
	r.node.mu.Lock()
	defer r.node.mu.Unlock()
	reply.ID = r.node.id
	reply.State = r.node.state
	reply.Term = r.node.currentTerm
	reply.LogLen = len(r.node.log)
	reply.Committed = r.node.commitIndex
	reply.Leader = r.node.votedFor
	return nil
}

func startRPCServer(n *Node, port string) error {
	handler := &RaftRPC{node: n}
	server := rpc.NewServer()
	if err := server.Register(handler); err != nil {
		return fmt.Errorf("register failed: %v", err)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen on %s failed: %v", port, err)
	}
	fmt.Printf("[RPC] Raft node listening on port %s\n", port)
	go server.Accept(ln)
	return nil
}

func callRPC(addr, method string, args, reply interface{}) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial %s: %v", addr, err)
	}
	defer client.Close()
	return client.Call(method, args, reply)
}
