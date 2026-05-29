// ============================================================
// Lab 01 - Topic 1: Distributed Systems Architectures
// FILE: client.go
// ROLE: Runs on Nodes 2, 3, 4, 5
//
// WHAT THIS CLIENT DOES:
//   - Connects to the server node over TCP
//   - Sends concurrent requests using goroutines
//   - Measures latency per request
//   - Handles server failures, timeouts, and rejections gracefully
//
// HOW TO RUN:
// HOW TO RUN (inside your container — Nodes 2-5):
//   ./client_bin -server 10.8.100.22:7000 -workers 3 -requests 30
//
// Replace 7000 with YOUR group's port (formula: 7000 + group_number - 1)
// Example for Group 3:  ./client_bin -server 10.8.100.22:7002 -workers 3 -requests 30
//
// EXPERIMENT FLAGS:
//   -workers   Number of concurrent goroutines (try 3, 10, 20)
//   -requests  Total requests to send        (try 30, 50, 100)
//   -delay     Delay between requests in ms  (default: 100)
//
// EXPERIMENT FLAGS:
//   -workers   Number of concurrent goroutines (simulates concurrent clients)
//   -requests  Total number of requests to send
//   -delay     Delay between requests in milliseconds
// ============================================================

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// ClientConfig holds runtime configuration.
type ClientConfig struct {
	ServerAddr  string
	NumWorkers  int
	NumRequests int
	DelayMs     int
	Timeout     time.Duration
}

// Result captures the outcome of a single request.
type Result struct {
	WorkerID  int
	RequestID int
	Success   bool
	Latency   time.Duration
	Error     string
}

// Metrics tracks aggregate results across all workers.
type Metrics struct {
	Successes    atomic.Int64
	Failures     atomic.Int64
	TotalLatency atomic.Int64 // in microseconds
}

func main() {
	// --- Parse command-line flags ---
	serverAddr := flag.String("server", "10.8.100.22:7000", "Server address (ip:port)")
	numWorkers := flag.Int("workers", 3, "Number of concurrent worker goroutines")
	numRequests := flag.Int("requests", 30, "Total requests to send")
	delayMs := flag.Int("delay", 100, "Delay between requests in milliseconds")
	flag.Parse()

	config := ClientConfig{
		ServerAddr:  *serverAddr,
		NumWorkers:  *numWorkers,
		NumRequests: *numRequests,
		DelayMs:     *delayMs,
		Timeout:     5 * time.Second,
	}

	log.Printf("[CLIENT] Connecting to server: %s", config.ServerAddr)
	log.Printf("[CLIENT] Workers: %d | Requests: %d | Delay: %dms",
		config.NumWorkers, config.NumRequests, config.DelayMs)

	metrics := &Metrics{}
	results := make(chan Result, config.NumRequests)

	// Distribute requests evenly among workers
	requestsPerWorker := config.NumRequests / config.NumWorkers

	// --- Launch worker goroutines ---
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runWorker(workerID, requestsPerWorker, config, metrics, results)
		}(i)
	}

	// Close results channel when all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// --- Collect and print results ---
	collectResults(results, metrics, startTime)
}

// runWorker sends a series of requests to the server.
// Each worker maintains its own TCP connection.
func runWorker(workerID, numRequests int, config ClientConfig, metrics *Metrics, results chan<- Result) {
	// -------------------------------------------------------
	// TODO (Students): Establish a TCP connection to the server.
	//
	//   1. Use net.DialTimeout to connect to config.ServerAddr
	//      with config.Timeout as the deadline.
	//      Hint: conn, err := net.DialTimeout("tcp", config.ServerAddr, config.Timeout)
	//
	//   2. If connection fails:
	//      - Log the error with the workerID
	//      - Send a failed Result to the results channel for each
	//        request this worker would have sent (so stats are accurate)
	//      - Return early
	//
	//   3. Defer conn.Close() so the connection is always cleaned up
	//
	// OBSERVE: When the server is killed mid-lab, which error do you get?
	//   - "connection refused" (server not running)
	//   - "i/o timeout" (server running but not responding)
	//   - "connection reset by peer" (server crashed mid-connection)
	// -------------------------------------------------------

	conn, err := net.DialTimeout("tcp", config.ServerAddr, config.Timeout)

	if err != nil {
		log.Printf("[WORKER %d] connection failed: %v", workerID, err)
		for i := 0; i < numRequests; i++ {
			results <- Result{WorkerID: workerID, RequestID: i, Success: false, Error: err.Error()}

		}
		return

	}

	defer conn.Close()

	log.Printf("[WORKER %d] Connected to server", workerID)

	reader := bufio.NewReader(conn) // TODO: replace nil with your conn

	for i := 0; i < numRequests; i++ {
		result := sendRequest(workerID, i, reader, conn, config) // TODO: pass conn
		results <- result

		// Update metrics
		if result.Success {
			metrics.Successes.Add(1)
			metrics.TotalLatency.Add(result.Latency.Microseconds())
		} else {
			metrics.Failures.Add(1)
		}

		// Delay between requests (simulates realistic client behaviour)
		time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
	}

	log.Printf("[WORKER %d] Finished sending %d requests", workerID, numRequests)
}

// sendRequest sends a single request to the server and returns the result.
func sendRequest(workerID, requestID int, reader *bufio.Reader, conn net.Conn, config ClientConfig) Result {
	result := Result{WorkerID: workerID, RequestID: requestID}
	start := time.Now()

	// -------------------------------------------------------
	// TODO (Students): Implement the full request/response cycle.
	//
	//   STEP 1 — Set a write deadline before sending:
	//     conn.SetWriteDeadline(time.Now().Add(config.Timeout))
	//
	//   STEP 2 — Send the request message to the server:
	//     Format: "REQUEST <workerID>-<requestID>\n"
	//     Hint: fmt.Fprintf(conn, "REQUEST %d-%d\n", workerID, requestID)
	//     If this fails, set result.Success = false, result.Error = err.Error()
	//     and return early.
	//
	//   STEP 3 — Set a read deadline before waiting for response:
	//     conn.SetReadDeadline(time.Now().Add(config.Timeout))
	//
	//   STEP 4 — Read the server's response:
	//     Hint: response, err := reader.ReadString('\n')
	//     If this fails (timeout, server crash), set result.Success = false
	//     and return early.
	//
	//   STEP 5 — Record latency and mark success:
	//     result.Latency = time.Since(start)
	//     result.Success = true
	//
	// EXPERIMENT:
	//   - What happens when you remove the deadlines?
	//   - What happens when the server is killed mid-request?
	//   - What happens with 100 workers vs 5 workers?
	// -------------------------------------------------------

	conn.SetWriteDeadline(time.Now().Add(config.Timeout))

	_, err := fmt.Fprintf(conn, "REQUEST %d-%d\n", workerID, requestID)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}

	conn.SetReadDeadline(time.Now().Add(config.Timeout))

	_, err1 := reader.ReadString('\n')

	if err1 != nil {
		result.Success = false
		result.Error = err1.Error()
		return result
	}

	result.Latency = time.Since(start)
	result.Success = true
	return result
}

// collectResults reads from the results channel and prints a summary.
func collectResults(results <-chan Result, metrics *Metrics, startTime time.Time) {
	// -------------------------------------------------------
	// TODO (Students): Implement result collection and reporting.
	//
	//   1. Range over the results channel (it closes when workers finish)
	//   2. For each failed result, print the error
	//   3. After the channel closes, print a summary:
	//      - Total time elapsed
	//      - Total successes and failures
	//      - Average latency (TotalLatency / Successes)
	//      - Requests per second (Successes / elapsed seconds)
	//
	// This summary will be your key observation during the failure
	// injection part of the lab.
	// -------------------------------------------------------

	for res := range results {
		if !res.Success {
			//print error
			log.Printf("[FAILURE] Worker %d, Req: %d: %s", res.WorkerID, res.RequestID, res.Error)
		}

	}
	elapsed := time.Since(startTime)

	successes := metrics.Successes.Load()
	failures := metrics.Failures.Load()

	fmt.Println("\n========== LAB 01 RESULTS ==========")
	fmt.Printf("Total time:     %.2fs\n", elapsed.Seconds())
	fmt.Printf("Successes:      %d\n", successes)
	fmt.Printf("Failures:       %d\n", failures)

	if successes > 0 {
		avgLatency := time.Duration(metrics.TotalLatency.Load()/successes) * time.Microsecond
		rps := float64(successes) / elapsed.Seconds()
		fmt.Printf("Avg latency:    %v\n", avgLatency)
		fmt.Printf("Throughput:     %.1f req/s\n", rps)
	}

	fmt.Println("=====================================")
}
