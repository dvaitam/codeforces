package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binaryPath := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	passed := 0
	numTests := 100

	for t := 0; t < numTests; t++ {
		// Generate test case
		// Use small n for easier debugging and speed
		n := rand.Intn(20) + 1
		k := rand.Intn(n) + 1
		
		// Constraint: k <= n.
		
		a := make([]int, n)
		for i := 0; i < n; i++ {
			// Values 1 to 2n
			a[i] = rand.Intn(2*n) + 1
		}

		inputStr := fmt.Sprintf("1\n%d %d\n", n, k)
		for i, v := range a {
			inputStr += fmt.Sprintf("%d", v)
			if i < n-1 {
				inputStr += " "
			}
		}
		inputStr += "\n"

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

		outputStr := strings.TrimSpace(out.String())
		outputInt, err := strconv.Atoi(outputStr)
		if err != nil {
			fmt.Printf("Test %d: Invalid output: %s\n", t, outputStr)
			os.Exit(1)
		}

		// Solve internally
		expected := solve(n, k, a)

		if outputInt != expected {
			fmt.Printf("Test %d FAILED\n", t)
			fmt.Printf("Input:\n%s", inputStr)
			fmt.Printf("Expected: %d\n", expected)
			fmt.Printf("Actual:   %d\n", outputInt)
			
			os.Exit(1)
		}
		passed++
	}
	fmt.Printf("All %d tests passed!\n", passed)
}

// solve implements the "Parallel Rounds" simulation
func solve(n, k int, a []int) int {
	// Frequency map
	cnt := make(map[int]int)
	maxVal := 0
	for _, v := range a {
		cnt[v]++
		if v > maxVal {
			maxVal = v
		}
	}

	ops := 0
	for {
		// Identify all x with cnt[x] > k
		// We need to iterate over current keys.
		// Since map iteration is random, we collect keys or iterate range.
		// Finding keys:
		var badKeys []int
		// Optimization: track min/max of active range?
		// Or just iterate map. n is small (20). Map size is small.
		for x, c := range cnt {
			if c > k {
				badKeys = append(badKeys, x)
			}
		}

		if len(badKeys) == 0 {
			break
		}

		ops++
		
		// Compute moves
		// We use a delta map to apply changes simultaneously
		delta := make(map[int]int)
		
		for _, x := range badKeys {
			c := cnt[x]
			move := c - 1 // Keep 1 logic
			// We remove 'move' from x, add 'move' to x+1
			// But wait: cnt[x] will be set to 1 directly?
			// In parallel simulation:
			// "Every drone ... is marked".
			// "Marked drones energy increased by 1".
			// So 'move' drones change from x to x+1.
			// The remaining '1' drone stays at x.
			// So effectively: cnt[x] -= move; cnt[x+1] += move.
			
			delta[x] -= move
			delta[x+1] += move
		}
		
		// Apply deltas
		for x, d := range delta {
			cnt[x] += d
			if cnt[x] == 0 {
				delete(cnt, x)
			}
		}
	}
	return ops
}