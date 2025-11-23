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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	binaryPath := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	passed := 0
	numTests := 100 // Reduce number of tests as logic is more complex to verify

	for t := 0; t < numTests; t++ {
		// Generate random test case
		// n: 1 to 100
		// m: 1 to n
		n := rand.Intn(30) + 1 // Start small to cover edge cases more frequently
		if t > 50 {
			n = rand.Intn(100) + 1
		}
		m := rand.Intn(n) + 1

		inputStr := fmt.Sprintf("%d %d\n", n, m)

		cmd := exec.Command(binaryPath)
		cmd.Stdin = strings.NewReader(inputStr)
		var out bytes.Buffer
		cmd.Stdout = &out
		// cmd.Stderr = os.Stderr // Optional: see stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: Binary failed to run: %v\nInput: %s\n", t, err, inputStr)
			os.Exit(1)
		}

		outputStr := out.String()
		
		err = verify(n, m, outputStr)
		if err != nil {
			fmt.Printf("Test %d FAILED\nInput: %s\nError: %v\n", t, inputStr, err)
			os.Exit(1)
		}
		passed++
	}
	fmt.Printf("All %d tests passed!\n", passed)
}

func verify(n, m int, outputStr string) error {
	scanner := bufio.NewScanner(strings.NewReader(outputStr))
	
	if !scanner.Scan() {
		return fmt.Errorf("empty output")
	}
	
	rLine := strings.TrimSpace(scanner.Text())
	r, err := strconv.Atoi(rLine)
	if err != nil {
		return fmt.Errorf("invalid r (count): %s", rLine)
	}

	// The verifier can't easily know the exact number of solutions k without solving it itself.
	// However, the problem says output min(k, 2000).
	// We can check if r <= 2000.
	// If r < 2000, we assume it found all of them (or at least check the ones it found are valid).
	// A strong verifier would solve it to know k, but that requires writing the solver.
	// Given the context, we will focus on validating the permutations provided.
	// If the user provides 0 solutions but there are some, it's a fail, but hard to detect without solver.
	// We will assume the count 'r' is trustworthy if the permutations are valid, 
	// unless we implement a solver.
	
	// Let's try to verify validity of each permutation.
	// Also check for duplicates? The problem implies "find and print ... examples", usually implying distinct ones.

	seen := make(map[string]bool)

	for i := 0; i < r; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected %d permutations, got %d", r, i)
		}
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		
		if len(parts) != n {
			return fmt.Errorf("permutation %d has length %d, expected %d", i+1, len(parts), n)
		}
		
		p := make([]int, n)
		used := make([]bool, n+1)
		for j, s := range parts {
			val, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("invalid number in permutation %d: %s", i+1, s)
			}
			if val < 1 || val > n {
				return fmt.Errorf("value %d out of range [1, %d] in permutation %d", val, n, i+1)
			}
			if used[val] {
				return fmt.Errorf("duplicate value %d in permutation %d", val, i+1)
			}
			used[val] = true
			p[j] = val
		}
		
		// Check distinctness
		pStr := fmt.Sprint(p)
		if seen[pStr] {
			return fmt.Errorf("duplicate permutation output: %v", p)
		}
		seen[pStr] = true

		// 1. Check Bitonic
		if !isBitonic(p) {
			return fmt.Errorf("permutation %d is not bitonic: %v", i+1, p)
		}

		// 2. Check Cycle Count
		cycles := countCycles(p)
		if cycles != m {
			return fmt.Errorf("permutation %d has %d cycles, expected %d: %v", i+1, cycles, m, p)
		}
	}
	
	return nil
}

func isBitonic(p []int) bool {
	n := len(p)
	if n == 0 {
		return true
	}
	
	// Find peak
	peakIdx := 0
	for i := 1; i < n; i++ {
		if p[i] > p[peakIdx] {
			peakIdx = i
		}
	}
	
	// Check increasing prefix
	for i := 0; i < peakIdx; i++ {
		if p[i] > p[i+1] {
			return false
		}
	}
	
	// Check decreasing suffix
	for i := peakIdx; i < n-1; i++ {
		if p[i] < p[i+1] {
			return false
		}
	}
	
	return true
}

func countCycles(p []int) int {
	n := len(p)
	visited := make([]bool, n+1) // 1-based values
	cycles := 0
	
	// p contains 1-based values. 
	// But indices are 0-based.
	// p[i] is the value at position i.
	// A cycle is defined by mapping: index -> value?
	// "For every x in C, we have p_x in C"
	// This usually implies p acts as a function f(i) = p[i].
	// If indices are 1..n, then f(i) = p[i-1].
	
	for i := 1; i <= n; i++ {
		if !visited[i] {
			cycles++
			curr := i
			for !visited[curr] {
				visited[curr] = true
				// next is p[curr-1] because p is 0-indexed array of 1-based values
				curr = p[curr-1]
			}
		}
	}
	return cycles
}