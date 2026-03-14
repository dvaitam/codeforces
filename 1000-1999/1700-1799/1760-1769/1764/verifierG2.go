package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// queryDoremy computes Q(l, r, k) on permutation p.
// l and r are 1-indexed, the subarray is p[l..r].
func queryDoremy(p []int, l, r, k int) int {
	seen := make(map[int]bool)
	for i := l - 1; i < r; i++ {
		seen[p[i]/k] = true
	}
	return len(seen)
}

func runTest(bin string, rng *rand.Rand, n int, testNum int) error {
	// Generate a random permutation of [1..n].
	perm := rng.Perm(n)
	p := make([]int, n)
	posOfOne := -1
	for i := 0; i < n; i++ {
		p[i] = perm[i] + 1 // 1-indexed values
		if p[i] == 1 {
			posOfOne = i + 1 // 1-indexed position
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %v", err)
	}
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}

	reader := bufio.NewReader(stdoutPipe)
	writer := bufio.NewWriter(stdinPipe)

	// Send n.
	fmt.Fprintf(writer, "%d\n", n)
	writer.Flush()

	const maxQueries = 25
	queries := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			stdinPipe.Close()
			cmd.Wait()
			return fmt.Errorf("read from candidate: %v", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '!' {
			// Answer line: "! y"
			var y int
			_, err := fmt.Sscanf(line, "! %d", &y)
			if err != nil {
				stdinPipe.Close()
				cmd.Wait()
				return fmt.Errorf("invalid answer format: %q", line)
			}
			stdinPipe.Close()
			cmd.Wait()
			if y != posOfOne {
				return fmt.Errorf("wrong answer: expected %d got %d (perm=%v)", posOfOne, y, p)
			}
			return nil
		} else if line[0] == '?' {
			// Query line: "? l r k"
			var l, r, k int
			_, err := fmt.Sscanf(line, "? %d %d %d", &l, &r, &k)
			if err != nil {
				stdinPipe.Close()
				cmd.Wait()
				return fmt.Errorf("invalid query format: %q", line)
			}
			queries++
			if queries > maxQueries {
				stdinPipe.Close()
				cmd.Wait()
				return fmt.Errorf("too many queries: %d > %d", queries, maxQueries)
			}
			if l < 1 || r > n || l > r || k < 1 || k > n {
				stdinPipe.Close()
				cmd.Wait()
				return fmt.Errorf("invalid query parameters: l=%d r=%d k=%d n=%d", l, r, k, n)
			}
			ans := queryDoremy(p, l, r, k)
			fmt.Fprintf(writer, "%d\n", ans)
			writer.Flush()
		} else {
			stdinPipe.Close()
			cmd.Wait()
			return fmt.Errorf("unexpected output: %q", line)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierG2 /path/to/candidate")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const numTests = 100
	for i := 0; i < numTests; i++ {
		n := rng.Intn(48) + 3 // n in [3, 50]; CF guarantees n >= 3
		if err := runTest(cand, rng, n, i+1); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
