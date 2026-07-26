package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (Follower="follower"; Candidate="candidate"; Leader="leader")

type LogEntry struct { Term, Index int; Command string }

type Node struct {
	mu           sync.Mutex
	id, addr     string
	peers        []string
	currentTerm  int
	votedFor     string
	log          []LogEntry
	state        string
	commitIndex  int
	lastApplied  int
	nextIndex    map[string]int
	matchIndex   map[string]int
	electionTimer *time.Timer
	stateMachine map[string]string
}

func NewNode(id, addr string, peers []string) *Node {
	n := &Node{
		id: id, addr: addr, peers: peers,
		state: Follower, currentTerm: 0, votedFor: "",
		log: []LogEntry{}, commitIndex: 0, lastApplied: 0,
		nextIndex: make(map[string]int),
		matchIndex: make(map[string]int),
		stateMachine: make(map[string]string),
	}
	fmt.Printf("[NODE %s] Starting as Follower (term 0)\n", id)
	n.resetElectionTimer()
	return n
}

func (n *Node) becomeFollower(term int) {
	n.state = Follower
	n.currentTerm = term
	n.votedFor = ""
	fmt.Printf("[NODE %s] Became Follower (term %d)\n", n.id, term)
	n.resetElectionTimer()
}

func (n *Node) becomeCandidate() {
	n.mu.Lock()
	n.currentTerm++
	n.state = Candidate
	n.votedFor = n.id
	fmt.Printf("[NODE %s] Became Candidate (term %d) — starting election\n", n.id, n.currentTerm)
	n.resetElectionTimer()
	n.mu.Unlock()
	go n.startElection()
}

func (n *Node) resetElectionTimer() {
	if n.electionTimer != nil { n.electionTimer.Stop() }
	timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	n.electionTimer = time.AfterFunc(timeout, func() {
		n.mu.Lock()
		isFollowerOrCandidate := n.state != Leader
		n.mu.Unlock()
		if isFollowerOrCandidate { n.becomeCandidate() }
	})
}

func (n *Node) lastLogIndex() int {
	if len(n.log) == 0 { return 0 }
	return n.log[len(n.log)-1].Index
}

func (n *Node) lastLogTerm() int {
	if len(n.log) == 0 { return 0 }
	return n.log[len(n.log)-1].Term
}

func (n *Node) String() string {
	return fmt.Sprintf("Node{id:%s state:%s term:%d log:%d committed:%d}",
		n.id, n.state, n.currentTerm, len(n.log), n.commitIndex)
}
