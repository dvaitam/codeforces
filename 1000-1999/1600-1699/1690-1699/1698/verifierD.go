package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// generateArray creates a permutation of [1..n] with exactly one fixed point
// by performing (n-1)/2 disjoint swaps on all elements except one.
func generateArray(rng *rand.Rand, n int) ([]int, int) {
	a := make([]int, n)
	for i := range a {
		a[i] = i + 1
	}
	// Pick the fixed point
	fixedIdx := rng.Intn(n)
	// Collect all other indices and shuffle them
	others := make([]int, 0, n-1)
	for i := 0; i < n; i++ {
		if i != fixedIdx {
			others = append(others, i)
		}
	}
	rng.Shuffle(len(others), func(i, j int) {
		others[i], others[j] = others[j], others[i]
	})
	// Pair them up and swap
	for i := 0; i+1 < len(others); i += 2 {
		a[others[i]], a[others[i+1]] = a[others[i+1]], a[others[i]]
	}
	return a, fixedIdx + 1 // return 1-indexed fixed point value (which equals fixedIdx+1)
}

func runInteractive(bin string, n int, a []int, fixedPoint int) error {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)

	// Send: 1 test case, then n
	fmt.Fprintf(writer, "1\n%d\n", n)
	writer.Flush()

	queries := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("unexpected EOF from candidate")
			}
			return fmt.Errorf("read: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		if parts[0] == "!" {
			// Answer
			if len(parts) != 2 {
				return fmt.Errorf("invalid answer format: %q", line)
			}
			ans, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid answer: %q", parts[1])
			}
			if ans != fixedPoint {
				return fmt.Errorf("wrong answer: got %d, expected %d", ans, fixedPoint)
			}
			stdin.Close()
			cmd.Wait()
			return nil
		}

		if parts[0] == "?" {
			queries++
			if queries > 15 {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("too many queries: %d > 15", queries)
			}
			if len(parts) != 3 {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("invalid query format: %q", line)
			}
			l, _ := strconv.Atoi(parts[1])
			r, _ := strconv.Atoi(parts[2])
			if l < 1 || r > n || l > r {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("invalid query range: l=%d r=%d n=%d", l, r, n)
			}
			// Extract subarray [l..r], sort, and respond
			sub := make([]int, r-l+1)
			copy(sub, a[l-1:r])
			sort.Ints(sub)
			var sb strings.Builder
			for i, v := range sub {
				if i > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", v)
			}
			fmt.Fprintln(writer, sb.String())
			writer.Flush()
		} else {
			// Unknown command, skip
			continue
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []int{3, 5, 7, 9, 11}
	for i := 0; i < 20; i++ {
		n := rng.Intn(7)*2 + 3 // random odd between 3 and 17
		tests = append(tests, n)
	}

	for i, n := range tests {
		a, fixedPoint := generateArray(rng, n)
		if err := runInteractive(bin, n, a, fixedPoint); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: %v\narray: %v\nfixed point: %d\n",
				i+1, n, err, a, fixedPoint)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
