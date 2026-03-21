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
	"sync"
)

// Generate a random permutation of 0..n-1.
func randPerm(rng *rand.Rand, n int) []int {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	rng.Shuffle(n, func(i, j int) { p[i], p[j] = p[j], p[i] })
	return p
}

// runInteractive spawns the candidate binary, feeds it n, answers "? i j"
// queries with (perm[i-1] | perm[j-1]), and verifies the "! p0 p1 ... pn-1"
// answer line.
func runInteractive(exe string, n int, perm []int) error {
	cmd := exec.Command(exe)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	// Send n
	if _, err := fmt.Fprintf(stdin, "%d\n", n); err != nil {
		return fmt.Errorf("failed to write n: %v", err)
	}

	queries := 0
	maxQueries := 4269
	var answerLine string
	var parseErr error

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
		for scanner.Scan() {
			ln := strings.TrimSpace(scanner.Text())
			if ln == "" {
				continue
			}
			if strings.HasPrefix(ln, "?") {
				parts := strings.Fields(ln)
				if len(parts) != 3 {
					fmt.Fprintf(stdin, "-1\n")
					continue
				}
				i, e1 := strconv.Atoi(parts[1])
				j, e2 := strconv.Atoi(parts[2])
				if e1 != nil || e2 != nil || i < 1 || i > n || j < 1 || j > n || i == j {
					fmt.Fprintf(stdin, "-1\n")
					continue
				}
				queries++
				if queries > maxQueries {
					fmt.Fprintf(stdin, "-1\n")
					continue
				}
				ans := perm[i-1] | perm[j-1]
				fmt.Fprintf(stdin, "%d\n", ans)
			} else if strings.HasPrefix(ln, "!") {
				answerLine = ln
				stdin.Close()
				return
			}
		}
	}()

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	if parseErr != nil {
		return parseErr
	}

	if answerLine == "" {
		return fmt.Errorf("no answer line received")
	}
	if queries > maxQueries {
		return fmt.Errorf("too many queries: %d > %d", queries, maxQueries)
	}

	// Parse answer: "! p0 p1 ... pn-1"
	parts := strings.Fields(answerLine)
	if len(parts) != n+1 || parts[0] != "!" {
		return fmt.Errorf("bad answer format: %q", answerLine)
	}
	got := make([]int, n)
	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(parts[i+1])
		if err != nil || v < 0 || v >= n {
			return fmt.Errorf("bad value in answer: %q", parts[i+1])
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d in answer", v)
		}
		seen[v] = true
		got[i] = v
	}

	// Verify: for all pairs i,j: (got[i] | got[j]) must equal (perm[i] | perm[j]).
	// Equivalently, got must equal perm (the permutation is unique given
	// the OR structure). Just check equality.
	for i := 0; i < n; i++ {
		if got[i] != perm[i] {
			return fmt.Errorf("mismatch at position %d: expected %d got %d", i+1, perm[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	// Test sizes: small cases and some medium ones
	sizes := []int{3, 3, 4, 4, 5, 5, 6, 7, 8, 10, 16, 32, 64, 128, 256, 512}
	for _, n := range sizes {
		perm := randPerm(rng, n)
		if err := runInteractive(exe, n, perm); err != nil {
			fmt.Fprintf(os.Stderr, "n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(sizes))
}
