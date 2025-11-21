package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxQueriesPerTest = 10
	maxK              = 20000
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	factDigits, err := precomputeFactorials(tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to precompute factorial digits: %v\n", err)
		os.Exit(1)
	}

	if err := runInteractive(candidate, tests, factDigits); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildTests() []int {
	manual := []int{
		1, 2, 3, 4, 5, 6, 9, 10, 12, 15,
		19, 25, 32, 50, 73, 100, 150, 227, 350, 512,
		777, 999, 1234, 1500, 1773, 2000, 2345, 3210, 4321, 5000, 5501, 5800, 5900, 5982,
	}
	seen := make(map[int]bool, len(manual))
	for _, x := range manual {
		if x >= 1 && x <= 5982 {
			seen[x] = true
		}
	}
	// Add a few deterministic "random" numbers to cover more of the range.
	extras := []int{42, 87, 271, 941, 1414, 2789, 3870, 4671}
	for _, x := range extras {
		if x >= 1 && x <= 5982 {
			seen[x] = true
		}
	}
	out := make([]int, 0, len(seen))
	for v := range seen {
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}

func precomputeFactorials(ns []int) (map[int][]byte, error) {
	if len(ns) == 0 {
		return nil, fmt.Errorf("no tests to precompute")
	}
	targets := make(map[int]bool, len(ns))
	var maxN int
	for _, v := range ns {
		targets[v] = true
		if v > maxN {
			maxN = v
		}
	}

	fact := big.NewInt(1)
	var mul big.Int
	res := make(map[int][]byte, len(ns))
	for i := 1; i <= maxN; i++ {
		mul.SetInt64(int64(i))
		fact.Mul(fact, &mul)
		if targets[i] {
			// Store digits in little-endian order for fast access by position.
			s := fact.String()
			digits := make([]byte, len(s))
			for idx := range s {
				digits[idx] = s[len(s)-1-idx]
			}
			res[i] = digits
		}
	}
	return res, nil
}

func digitAt(digits []byte, k int) int {
	if k < 0 || k >= len(digits) {
		return 0
	}
	return int(digits[k] - '0')
}

func runInteractive(path string, tests []int, factDigits map[int][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := commandFor(ctx, path)
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
		return fmt.Errorf("failed to start candidate: %w", err)
	}

	writer := bufio.NewWriter(stdin)
	reader := bufio.NewReader(stdout)

	if _, err := fmt.Fprintf(writer, "%d\n", len(tests)); err != nil {
		return fmt.Errorf("failed to send test count: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush test count: %w", err)
	}

	for idx, n := range tests {
		digits, ok := factDigits[n]
		if !ok {
			return fmt.Errorf("missing factorial digits for n=%d", n)
		}
		queries := 0
		for {
			tok, err := nextToken(reader)
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					return fmt.Errorf("test %d: timed out waiting for output", idx+1)
				}
				return fmt.Errorf("test %d: failed to read token: %v\nstderr:\n%s", idx+1, err, stderr.String())
			}
			if len(tok) == 0 {
				continue
			}
			switch tok[0] {
			case '?':
				kStr := strings.TrimSpace(tok[1:])
				if kStr == "" {
					kStr, err = nextToken(reader)
					if err != nil {
						return fmt.Errorf("test %d: expected index after '?': %v", idx+1, err)
					}
				}
				k, err := strconv.Atoi(kStr)
				if err != nil {
					return fmt.Errorf("test %d: invalid query index %q", idx+1, kStr)
				}
				if k < 0 || k >= maxK {
					return fmt.Errorf("test %d: query index out of range (%d)", idx+1, k)
				}
				queries++
				if queries > maxQueriesPerTest {
					return fmt.Errorf("test %d: query limit exceeded (%d)", idx+1, queries)
				}
				fmt.Fprintf(writer, "%d\n", digitAt(digits, k))
				if err := writer.Flush(); err != nil {
					return fmt.Errorf("test %d: failed to flush answer: %w", idx+1, err)
				}
			case '!':
				answerStr := strings.TrimSpace(tok[1:])
				if answerStr == "" {
					answerStr, err = nextToken(reader)
					if err != nil {
						return fmt.Errorf("test %d: expected answer after '!': %v", idx+1, err)
					}
				}
				guess, err := strconv.Atoi(answerStr)
				if err != nil {
					return fmt.Errorf("test %d: invalid guess %q", idx+1, answerStr)
				}
				if guess != n {
					fmt.Fprintln(writer, "NO")
					writer.Flush()
					return fmt.Errorf("test %d: wrong guess %d (expected %d)", idx+1, guess, n)
				}
				fmt.Fprintln(writer, "YES")
				if err := writer.Flush(); err != nil {
					return fmt.Errorf("test %d: failed to flush YES: %w", idx+1, err)
				}
				break
			default:
				return fmt.Errorf("test %d: unexpected token %q", idx+1, tok)
			}
			if tok[0] == '!' {
				break
			}
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush final data: %w", err)
	}
	if err := stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("candidate timed out")
		}
		return fmt.Errorf("candidate runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return nil
}

func nextToken(r *bufio.Reader) (string, error) {
	var tok string
	_, err := fmt.Fscan(r, &tok)
	return tok, err
}

func commandFor(ctx context.Context, path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.CommandContext(ctx, "go", "run", path)
	case ".py":
		return exec.CommandContext(ctx, "python3", path)
	default:
		return exec.CommandContext(ctx, path)
	}
}
