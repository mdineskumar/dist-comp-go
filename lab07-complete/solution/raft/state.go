package main

import (
	"fmt"
	"strings"
)

func (n *Node) applyToStateMachine(command string) {
	parts := strings.Fields(command)
	if len(parts) == 0 { return }
	switch strings.ToUpper(parts[0]) {
	case "SET":
		if len(parts) >= 3 {
			n.stateMachine[parts[1]] = parts[2]
			fmt.Printf("[STATE %s] SET %s=%s\n", n.id, parts[1], parts[2])
		}
	case "DEL":
		if len(parts) >= 2 {
			delete(n.stateMachine, parts[1])
			fmt.Printf("[STATE %s] DEL %s\n", n.id, parts[1])
		}
	default:
		fmt.Printf("[STATE %s] Unknown command: %s\n", n.id, command)
	}
}

func (n *Node) GetValue(key string) (string, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	val, ok := n.stateMachine[key]
	return val, ok
}

func (n *Node) PrintStatus() {
	n.mu.Lock()
	defer n.mu.Unlock()
	fmt.Printf("\n── Node Status ────────────────────────────────\n")
	fmt.Printf("  ID:           %s\n", n.id)
	fmt.Printf("  State:        %s\n", n.state)
	fmt.Printf("  Term:         %d\n", n.currentTerm)
	fmt.Printf("  Log length:   %d\n", len(n.log))
	fmt.Printf("  Committed:    %d\n", n.commitIndex)
	fmt.Printf("  Applied:      %d\n", n.lastApplied)
	fmt.Printf("  State machine: %d keys\n", len(n.stateMachine))
	for k, v := range n.stateMachine {
		fmt.Printf("    %s = %s\n", k, v)
	}
	fmt.Println()
}
