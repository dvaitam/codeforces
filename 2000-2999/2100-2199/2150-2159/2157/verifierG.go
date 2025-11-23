package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	binaryPath := os.Args[1]

	cmd := exec.Command(binaryPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// Redirect stderr to see any debug output/panics from the solution
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	defer func() {
		// Ensure process is killed on exit
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	scanner := bufio.NewScanner(stdout)
	// Increase buffer size just in case output lines are long (though n=100 lines are short)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	writer := bufio.NewWriter(stdin)

	numTestCases := 30
	budget := 300.0

	// Send total test cases
	fmt.Fprintf(writer, "%d\n", numTestCases)
	writer.Flush()

	rand.Seed(time.Now().UnixNano())

	for t := 1; t <= numTestCases; t++ {
		var n int
		var a []int
		if t == 1 {
			// Example case from problem description
			n = 3
			a = []int{2, 4, 6}
		} else {
			n = 100
			a = make([]int, n)
			for i := 0; i < n; i++ {
				a[i] = rand.Intn(1 << 30)
			}
		}

		// Precompute prefix xors
		// p[i] = a[0] ^ ... ^ a[i-1]
		// XOR sum of a[u...v] (1-based u, v) corresponds to p[v] ^ p[u-1]
		p := make([]int, n+1)
		for i := 0; i < n; i++ {
			p[i+1] = p[i] ^ a[i]
		}

		// Send n
		fmt.Fprintf(writer, "%d\n", n)
		writer.Flush()

		// Interaction loop for current test case
		testFinished := false
		for !testFinished {
			if !scanner.Scan() {
				fmt.Printf("Test %d: Unexpected EOF from binary (scanner error: %v)\n", t, scanner.Err())
				return
			}
			line := strings.TrimSpace(scanner.Text())
			
			if line == "!" {
				// Verify answers
				// We expect n lines
				for i := 1; i <= n; i++ {
					if !scanner.Scan() {
						fmt.Printf("Test %d: Expected answer line for i=%d\n", t, i)
						return
					}
					ansLine := strings.Fields(scanner.Text())
					
					// The i-th line should contain answers for j from i to n
					// Count is n - i + 1
					expectedCount := n - i + 1
					if len(ansLine) != expectedCount {
						fmt.Printf("Test %d: Line %d has wrong number of answers. Expected %d, got %d\n", t, i, expectedCount, len(ansLine))
						return
					}

					for offset, valStr := range ansLine {
						j := i + offset // j ranges from i to n
						
						val, err := strconv.Atoi(valStr)
						if err != nil {
							fmt.Printf("Test %d: Invalid integer answer at (%d, %d): %s\n", t, i, j, valStr)
							return
						}
						
						// Calculate ground truth
						xorVal := p[j] ^ p[i-1]
						expected := -1
						if xorVal != 0 {
							expected = bits.Len(uint(xorVal)) - 1
						}
						
						if val != expected {
							fmt.Printf("Test %d: Wrong answer for query %d %d. Expected %d, got %d\n", t, i, j, expected, val)
							return
						}
					}
				}
				testFinished = true
			} else if strings.HasPrefix(line, "?") {
				parts := strings.Fields(line)
				if len(parts) != 3 {
					fmt.Printf("Test %d: Invalid query format: %s\n", t, line)
					return
				}
				u, err1 := strconv.Atoi(parts[1])
				v, err2 := strconv.Atoi(parts[2])
				
				if err1 != nil || err2 != nil {
					fmt.Printf("Test %d: Invalid query arguments: %s\n", t, line)
					return
				}
				
				if u < 1 || v > n || u > v {
					fmt.Printf("Test %d: Invalid query range: %d %d\n", t, u, v)
					// Protocol: -2
					fmt.Fprintf(writer, "-2\n")
					writer.Flush()
					return
				}

				cost := 1.0 / float64(v - u + 1)
				budget -= cost
				if budget < -1e-9 {
					fmt.Printf("Test %d: Budget exceeded! Cost %.5f, Remaining %.5f\n", t, cost, budget)
					fmt.Fprintf(writer, "-2\n")
					writer.Flush()
					return
				}

				xorVal := p[v] ^ p[u-1]
				resp := -1
				if xorVal != 0 {
					resp = bits.Len(uint(xorVal)) - 1
				}
				
				fmt.Fprintf(writer, "%d\n", resp)
				writer.Flush()
			} else {
				fmt.Printf("Test %d: Unknown command: %s\n", t, line)
				return
			}
		}
		// fmt.Printf("Test %d passed. Remaining budget: %.4f\n", t, budget)
	}

	fmt.Printf("All tests passed! Final remaining budget: %.4f\n", budget)
}