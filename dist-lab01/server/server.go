// ============================================================
// Lab 01 - Topic 1: Distributed Systems Architectures
// FILE: server.go
// ROLE: Runs on Node 1 (10.8.100.22)
//
// WHAT THIS SERVER DOES:
//   - Listens for incoming TCP connections from client nodes
//   - Handles each client in a separate goroutine (concurrently)
//   - Tracks active connections and request statistics
//   - Demonstrates what happens during failure / overload
//
// HOW TO RUN:
//   go run server.go
// HOW TO RUN (inside your container — Node 1 only):
//   First kill auto-started server:  pkill server_bin
//   Recompile after editing:         go build -o server_bin server.go
//   Run with your group port:         ./server_bin -port 7002
//
// Replace 7002 with YOUR group's port (formula: 7000 + group_number - 1)
// Example: Group 1 → port 7000, Group 5 → port 7004, Group 20 → port 7019
// ============================================================

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// ServerConfig holds runtime configuration for the server.
// Students can experiment by changing these values via flags.
type ServerConfig struct {
	Port        string
	MaxConns    int
	ReadTimeout time.Duration
}

// Stats tracks live server metrics.
// atomic.Int64 ensures thread-safe updates without locks.
type Stats struct {
	TotalRequests  atomic.Int64
	ActiveConns    atomic.Int64
	RejectedConns  atomic.Int64
	FailedRequests atomic.Int64
}

// Server is the main server struct holding state.
type Server struct {
	config   ServerConfig
	stats    Stats
	listener net.Listener
	wg       sync.WaitGroup // waits for all goroutines to finish on shutdown
	quit     chan struct{}  // signals goroutines to stop
}

func main() {
	// --- Parse command-line flags ---
	port := flag.String("port", "7000", "Port to listen on")
	maxconn := flag.Int("maxconn", 100, "Maximum concurrent connections")
	flag.Parse()

	config := ServerConfig{
		Port:        *port,
		MaxConns:    *maxconn,
		ReadTimeout: 5 * time.Second,
	}

	srv := &Server{
		config: config,
		quit:   make(chan struct{}),
	}

	// --- Start the server ---
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// --- Print stats every 2 seconds so students can watch live ---
	go srv.printStats()

	// --- Wait for Ctrl+C to gracefully shut down ---
	srv.waitForShutdown()
}

// Start begins listening for TCP connections.
func (s *Server) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%s", s.config.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not bind to %s: %w", addr, err)
	}
	s.listener = ln
	log.Printf("[SERVER] Listening on %s (max connections: %d)", addr, s.config.MaxConns)

	// Accept connections in a goroutine so main() is not blocked
	s.wg.Add(1)
	go s.acceptLoop()

	return nil
}

// acceptLoop continuously accepts new connections until the server is stopped.
func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// Check if we deliberately closed the listener (shutdown)
			select {
			case <-s.quit:
				log.Println("[SERVER] Accept loop shutting down.")
				return
			default:
				log.Printf("[SERVER] Accept error: %v", err)
				continue
			}
		}

		// -------------------------------------------------------
		// TODO (Students): Implement connection limiting.
		//
		// Currently the server accepts ALL connections regardless
		// of how many are active. This is unrealistic.
		//
		// Your task:
		//   1. Check s.stats.ActiveConns against s.config.MaxConns
		//   2. If at capacity:
		//        - Increment s.stats.RejectedConns
		//        - Send a rejection message to the client
		//        - Close the connection
		//        - Log a warning
		//        - return (do NOT call handleConnection)
		//   3. If not at capacity, fall through to handleConnection
		//
		// Hint: use s.stats.ActiveConns.Load() to read the value
		// -------------------------------------------------------
		if s.stats.ActiveConns.Load() >= int64(s.config.MaxConns) {
			s.stats.RejectedConns.Add(1)
			fmt.Fprintf(conn, "REJECTED: Server at maximum capacity\n")
			log.Printf("[WARNING] Rejecting connection from %s: Max Capacity (%d).", conn.RemoteAddr().String(), s.config.MaxConns)
			conn.Close()
			return
		}
		// Handle accepted connection in its own goroutine
		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// handleConnection processes a single client connection.
// It runs in its own goroutine so the server can handle many clients at once.
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	// Track this connection as active
	s.stats.ActiveConns.Add(1)
	defer s.stats.ActiveConns.Add(-1) // decrement when this goroutine exits

	clientAddr := conn.RemoteAddr().String()
	log.Printf("[SERVER] New connection from %s", clientAddr)

	// -------------------------------------------------------
	// TODO (Students): Implement the request handling loop.
	//
	// A client may send MULTIPLE requests on the same connection.
	// You need a loop that:
	//   1. Sets a read deadline on the connection using s.config.ReadTimeout
	//      Hint: conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
	//
	//   2. Reads a request from the client
	//      Hint: use a bufio.Reader and ReadString('\n')
	//
	//   3. If read fails (timeout or disconnect), break the loop
	//
	//   4. Increments s.stats.TotalRequests
	//
	//   5. Processes the request:
	//      - Parse the message (it will be plain text, e.g. "REQUEST 42\n")
	//      - Simulate work with a small time.Sleep (e.g. 50ms)
	//
	//   6. Sends a response back to the client:
	//      - Format: "RESPONSE OK <timestamp>\n"
	//      - Hint: fmt.Fprintf(conn, "RESPONSE OK %s\n", time.Now().Format(time.RFC3339))
	//
	//   7. On any write error, increment s.stats.FailedRequests and break
	//
	// EXPERIMENT: What happens if you increase the Sleep duration?
	// -------------------------------------------------------

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))

		line, err := reader.ReadString('\n')

		if err != nil {
			//conn.Close()
			break
		}

		log.Printf("Read: %s", line)
		s.stats.TotalRequests.Add(1)

		time.Sleep(time.Millisecond * 50)

		_, write_err := fmt.Fprintf(conn, "RESPONSE OK %s\n", time.Now().Format(time.RFC3339))

		if write_err != nil {
			s.stats.FailedRequests.Add(1)
			break
		}
	}
	log.Printf("[SERVER] Connection closed: %s", clientAddr)
}

// printStats logs server statistics every 2 seconds.
// Students should watch this during the lab to observe the effect of failures.
func (s *Server) printStats() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\n[STATS] Active: %d | Total Requests: %d | Rejected: %d | Failed: %d\n",
				s.stats.ActiveConns.Load(),
				s.stats.TotalRequests.Load(),
				s.stats.RejectedConns.Load(),
				s.stats.FailedRequests.Load(),
			)
		case <-s.quit:
			return
		}
	}
}

// waitForShutdown blocks until SIGINT or SIGTERM is received,
// then gracefully shuts down the server.
func (s *Server) waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("\n[SERVER] Shutdown signal received. Closing listener...")
	close(s.quit)      // signal all goroutines to stop
	s.listener.Close() // unblock acceptLoop

	// -------------------------------------------------------
	// TODO (Students): Implement graceful shutdown.
	//
	// Currently, the server just exits immediately.
	// Your task:
	//   1. Wait for all active connections to finish using s.wg.Wait()
	//   2. Print final stats before exiting
	//   3. Add a timeout: if connections don't finish in 10 seconds,
	//      force exit anyway.
	//
	// THINK ABOUT: In a real distributed system, what happens to
	// clients that were mid-request when the server shuts down?
	// -------------------------------------------------------
	done := make(chan struct{})

	go func() {
		s.wg.Wait()
		close(done)
	}()
	//s.wg.Wait()

	select {
	case <-done:
		log.Println("\n[SERVER] all active connections finished successfully")
	case <-time.After(10 * time.Second):
		log.Println("ln[SERVER] shutdown timeout reached. some connections were forced to close")
	}

	fmt.Printf("total requests: %d\n", s.stats.TotalRequests.Load())
	fmt.Printf("rejected: %d\n", s.stats.RejectedConns.Load())
	fmt.Printf("failed: %d\n", s.stats.FailedRequests.Load())

	log.Println("[SERVER] Shutdown complete.")
}
