package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Mission struct {
	y, l int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binaryPath := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	passed := 0
	numTests := 100

	for t := 0; t < numTests; t++ {
		// Generate test case
		// Random N.
		// Mix of small and large N.
		var n int
		if t < 10 {
			n = 4 // Example case
		} else if t < 90 {
			n = rand.Intn(2000) + 2 // Small to medium
		} else {
			// Larger cases, up to 250000
			n = rand.Intn(250000-2000) + 2000
		}

		inputStr := fmt.Sprintf("%d\n", n)

		// Run binary
		cmd := exec.Command(binaryPath)
		cmd.Stdin = strings.NewReader(inputStr)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: Binary failed to run: %v\n", t, err)
			fmt.Println("Input:", inputStr)
			os.Exit(1)
		}

		// Parse Output
		outputStr := out.String()
		scanner := bufio.NewScanner(strings.NewReader(outputStr))
		
		if !scanner.Scan() {
			fmt.Printf("Test %d: Empty output\n", t)
			os.Exit(1)
		}
		kLine := strings.TrimSpace(scanner.Text())
		k, err := strconv.Atoi(kLine)
		if err != nil {
			fmt.Printf("Test %d: Invalid k: %s\n", t, kLine)
			os.Exit(1)
		}

		missions := make([]Mission, 0, k)
		for i := 0; i < k; i++ {
			if !scanner.Scan() {
				fmt.Printf("Test %d: Not enough mission lines. Expected %d, got %d\n", t, k, i)
				os.Exit(1)
			}
			line := strings.Fields(scanner.Text())
			if len(line) != 2 {
				fmt.Printf("Test %d: Invalid mission format at line %d: %s\n", t, i+1, scanner.Text())
				os.Exit(1)
			}
			y, _ := strconv.Atoi(line[0])
			l, _ := strconv.Atoi(line[1])
			missions = append(missions, Mission{y, l})
		}

		// Verification
		err = verify(n, missions)
		if err != nil {
			fmt.Printf("Test %d FAILED (N=%d)\n", t, n)
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		passed++
	}
	fmt.Printf("All %d tests passed!\n", passed)
}

func verify(n int, missions []Mission) error {
	// 1. Check Budget
	// Budget = 10^6
	cost := 0
	lastDiff := -1 // undefined

	for i, m := range missions {
		if m.y < 1 || m.l < 1 {
			return fmt.Errorf("mission %d has invalid values: y=%d, l=%d", i, m.y, m.l)
		}
		if m.y > 1000000 || m.l > 1000000 {
			// Problem says y, l <= 10^6.
			// Strict check?
		}

		stepCost := 0
		if i == 0 {
			stepCost = m.l
		} else {
			if m.y <= lastDiff {
				stepCost = m.l
			} else {
				stepCost = m.l + 1000
			}
		}
		
		cost += stepCost
		lastDiff = m.y
	}

	if cost > 1000000 {
		return fmt.Errorf("total cost %d exceeds budget 1000000", cost)
	}

	// 2. Check Logic (Coverage)
	// We need to ensure that for all start s in [1, n], final skill >= n.
	// Equivalently, all s in [1, n-1] must eventually move to >= n.
	
	// active[i] is true if there is a starting scenario currently at skill level i.
	// We track indices 1 to n-1.
	// If a scenario moves to >= n, it drops out of the map/array.
	
	// Since n <= 250000, we can use a boolean slice.
	active := make([]bool, n) // 0..n-1. We use indices 1..n-1.
	activeCount := 0
	for i := 1; i < n; i++ {
		active[i] = true
		activeCount++
	}

	for _, m := range missions {
		if activeCount == 0 {
			break
		}
		
		// If mission difficulty y is within range of tracked skills
		if m.y < n && active[m.y] {
			// All scenarios currently at m.y move to m.y + m.l
			active[m.y] = false
			activeCount--
			
			next := m.y + m.l
			if next < n {
				if !active[next] {
					active[next] = true
					activeCount++
				}
			}
		}
	}

	if activeCount > 0 {
		// Collect some failed start points for debug
		failed := []int{}
		for i := 1; i < n; i++ {
			if active[i] {
				failed = append(failed, i)
				if len(failed) >= 5 {
					break
				}
			}
		}
		return fmt.Errorf("%d starting skills did not reach %d. Examples: %v", activeCount, n, failed)
	}

	return nil
}