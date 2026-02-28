package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	maxW           = 100000
	maxQueries     = 2
	timeoutSeconds = 30
)

type testCase struct {
	w int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[0]

	tests := buildTests()

	if err := runInteractive(candidate, tests); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

// simulateEditor computes the number of lines the editor uses to display
// the article with the given word lengths and width W.
// Returns 0 if any word exceeds W.
func simulateEditor(words []int, W int) int {
	for _, a := range words {
		if a > W {
			return 0
		}
	}
	lines := 1
	s := 0
	for _, a := range words {
		if s+a <= W {
			s += a
		} else {
			lines++
			s = a
		}
	}
	return lines
}

func runInteractive(path string, tests []testCase) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	cmd := commandFor(ctx, path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start candidate: %w", err)
	}
	defer func() {
		stdin.Close()
		cmd.Wait()
	}()

	writer := bufio.NewWriter(stdin)
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)

	// Send the number of test cases
	fmt.Fprintf(writer, "%d\n", len(tests))
	writer.Flush()

	for idx, tc := range tests {
		W := tc.w
		queries := 0

		for {
			if !scanner.Scan() {
				if ctx.Err() == context.DeadlineExceeded {
					return fmt.Errorf("test %d (W=%d): timed out waiting for output", idx+1, W)
				}
				return fmt.Errorf("test %d (W=%d): unexpected EOF from candidate\nstderr: %s", idx+1, W, stderrBuf.String())
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "!") {
				// Answer line: "! W"
				parts := strings.Fields(line)
				if len(parts) != 2 {
					return fmt.Errorf("test %d (W=%d): malformed answer line: %q", idx+1, W, line)
				}
				ans, err := strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("test %d (W=%d): invalid answer %q", idx+1, W, parts[1])
				}
				if ans != W {
					return fmt.Errorf("test %d: wrong answer %d (expected %d)", idx+1, ans, W)
				}
				break // move to next test case
			} else if strings.HasPrefix(line, "?") {
				// Query line: "? n a1 a2 ... an"
				queries++
				if queries > maxQueries {
					// Send -1 to signal error, then fail
					fmt.Fprintln(writer, -1)
					writer.Flush()
					return fmt.Errorf("test %d (W=%d): too many queries (%d > %d)", idx+1, W, queries, maxQueries)
				}

				parts := strings.Fields(line)
				if len(parts) < 2 {
					fmt.Fprintln(writer, -1)
					writer.Flush()
					return fmt.Errorf("test %d (W=%d): malformed query: %q", idx+1, W, line)
				}
				n, err := strconv.Atoi(parts[1])
				if err != nil || n < 1 || n > maxW {
					fmt.Fprintln(writer, -1)
					writer.Flush()
					return fmt.Errorf("test %d (W=%d): invalid n=%q in query", idx+1, W, parts[1])
				}
				if len(parts) != n+2 {
					fmt.Fprintln(writer, -1)
					writer.Flush()
					return fmt.Errorf("test %d (W=%d): expected %d words but got %d tokens", idx+1, W, n, len(parts)-2)
				}
				words := make([]int, n)
				for i := 0; i < n; i++ {
					val, err := strconv.Atoi(parts[i+2])
					if err != nil || val < 1 || val > maxW {
						fmt.Fprintln(writer, -1)
						writer.Flush()
						return fmt.Errorf("test %d (W=%d): invalid word length %q", idx+1, W, parts[i+2])
					}
					words[i] = val
				}

				result := simulateEditor(words, W)
				fmt.Fprintln(writer, result)
				writer.Flush()
			} else {
				return fmt.Errorf("test %d (W=%d): unexpected line: %q", idx+1, W, line)
			}
		}
	}

	stdin.Close()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("candidate runtime error: %v\nstderr: %s", err, stderrBuf.String())
	}
	return nil
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

func buildTests() []testCase {
	tests := make([]testCase, 0)
	add := func(w int) {
		if w < 1 {
			w = 1
		}
		if w > maxW {
			w = maxW
		}
		tests = append(tests, testCase{w: w})
	}

	// Deterministic edge cases.
	add(1)
	add(maxW)
	add(50000)
	add(99999)

	// Randomized cases.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 10 {
		add(1 + rng.Intn(maxW))
	}
	return tests
}
