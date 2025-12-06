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

func generateCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(100) + 2 // n between 2 and 101
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}

	// Apply 0 to 3 reversals
	k := rng.Intn(4)
	for i := 0; i < k; i++ {
		l := rng.Intn(n)
		r := rng.Intn(n)
		if l > r {
			l, r = r, l
		}
		// Reverse perm[l...r] (inclusive, 0-based indices)
		for l < r {
			perm[l], perm[r] = perm[r], perm[l]
			l++
			r--
		}
	}
	return n, perm
}

func runSolution(bin string, n int, perm []int) (string, error) {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verify(n int, initialPerm []int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return fmt.Errorf("empty output")
	}
	kStr := scanner.Text()
	k, err := strconv.Atoi(kStr)
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}

	if k < 0 || k > 3 {
		return fmt.Errorf("k out of bounds: %d", k)
	}

	// Start with sorted identity permutation
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}

	// Apply k reversals
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected l for op %d", i+1)
		}
		lStr := scanner.Text()
		l, err := strconv.Atoi(lStr)
		if err != nil {
			return fmt.Errorf("invalid l at op %d: %v", i+1, err)
		}

		if !scanner.Scan() {
			return fmt.Errorf("expected r for op %d", i+1)
		}
		rStr := scanner.Text()
		r, err := strconv.Atoi(rStr)
		if err != nil {
			return fmt.Errorf("invalid r at op %d: %v", i+1, err)
		}

		if l < 1 || l > n || r < 1 || r > n {
			return fmt.Errorf("indices out of bounds: %d %d", l, r)
		}
		if l > r {
			l, r = r, l
		}

		// Apply reversal (1-based to 0-based)
		start, end := l-1, r-1
		for start < end {
			perm[start], perm[end] = perm[end], perm[start]
			start++
			end--
		}
	}

	// Check if matches initialPerm (the target input)
	for i := 0; i < n; i++ {
		if perm[i] != initialPerm[i] {
			return fmt.Errorf("result does not match input permutation. Element at %d is %d, expected %d", i, perm[i], initialPerm[i])
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		n, perm := generateCase(rng)
		output, err := runSolution(bin, n, perm)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			fmt.Printf("Input n=%d\nPerm: %v\n", n, perm)
			os.Exit(1)
		}
		if err := verify(n, perm, output); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			fmt.Printf("Input n=%d\nPerm: %v\nOutput:\n%s\n", n, perm, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All 100 tests passed\n")
}
