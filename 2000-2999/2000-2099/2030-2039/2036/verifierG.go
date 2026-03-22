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
)

type testCase struct {
	n, a, b, c int64
	desc       string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()

	for i, tc := range tests {
		if err := runInteractive(candidate, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\n", i+1, tc.desc, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

// xorOfMissingInRange returns XOR of the missing values (a, b, c) that fall within [l, r].
func xorOfMissingInRange(l, r int64, missing [3]int64) int64 {
	var result int64
	for _, v := range missing {
		if v >= l && v <= r {
			result ^= v
		}
	}
	return result
}

// runInteractive runs the candidate binary with a single test case, simulating
// the interactive protocol. The candidate sends "xor L R" queries and receives
// the XOR of the missing values in [L, R]. Then it sends "ans a b c".
func runInteractive(binPath string, tc testCase) error {
	cmd := exec.Command(binPath)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %v", err)
	}
	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}

	reader := bufio.NewReader(stdoutPipe)
	writer := bufio.NewWriter(stdinPipe)

	missing := [3]int64{tc.a, tc.b, tc.c}

	// Write: t=1, then n
	fmt.Fprintf(writer, "1\n%d\n", tc.n)
	writer.Flush()

	queryCount := 0
	maxQueries := 200
	done := false

	for !done {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			_ = cmd.Process.Kill()
			return fmt.Errorf("read from candidate: %v (stderr: %s)", err, stderrBuf.String())
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		switch strings.ToLower(tokens[0]) {
		case "xor":
			if len(tokens) != 3 {
				_ = cmd.Process.Kill()
				return fmt.Errorf("invalid xor query: %q", line)
			}
			l, err1 := strconv.ParseInt(tokens[1], 10, 64)
			r, err2 := strconv.ParseInt(tokens[2], 10, 64)
			if err1 != nil || err2 != nil {
				_ = cmd.Process.Kill()
				return fmt.Errorf("invalid xor args: %q", line)
			}
			queryCount++
			if queryCount > maxQueries {
				_ = cmd.Process.Kill()
				return fmt.Errorf("too many queries (%d)", queryCount)
			}

			// The query returns XOR of the missing values in [l, r].
			result := xorOfMissingInRange(l, r, missing)

			fmt.Fprintf(writer, "%d\n", result)
			writer.Flush()

		case "ans":
			if len(tokens) != 4 {
				_ = cmd.Process.Kill()
				return fmt.Errorf("invalid ans line: %q", line)
			}
			va, err1 := strconv.ParseInt(tokens[1], 10, 64)
			vb, err2 := strconv.ParseInt(tokens[2], 10, 64)
			vc, err3 := strconv.ParseInt(tokens[3], 10, 64)
			if err1 != nil || err2 != nil || err3 != nil {
				_ = cmd.Process.Kill()
				return fmt.Errorf("invalid ans args: %q", line)
			}

			got := []int64{va, vb, vc}
			sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
			exp := []int64{tc.a, tc.b, tc.c}
			sort.Slice(exp, func(i, j int) bool { return exp[i] < exp[j] })

			if got[0] != exp[0] || got[1] != exp[1] || got[2] != exp[2] {
				_ = cmd.Process.Kill()
				return fmt.Errorf("wrong answer: got %v, expected %v", got, exp)
			}
			done = true

		default:
			_ = cmd.Process.Kill()
			return fmt.Errorf("unexpected command: %q", line)
		}
	}

	stdinPipe.Close()
	cmd.Wait()

	if !done {
		return fmt.Errorf("candidate exited without providing answer")
	}

	return nil
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 6, a: 2, b: 3, c: 5, desc: "sample-like"},
		{n: 3, a: 1, b: 2, c: 3, desc: "small-all-missing"},
		{n: 10, a: 1, b: 5, c: 9, desc: "small-mixed"},
		{n: 100, a: 10, b: 20, c: 30, desc: "mid-range"},
		{n: 10_000, a: 1, b: 9999, c: 5000, desc: "spread"},
		{n: 1_000_000_000_000, a: 123456789, b: 234567890, c: 345678901, desc: "large"},
		{n: 999_999_999_999_999_999, a: 111111111111111111, b: 222222222222222222, c: 333333333333333333, desc: "near-limit"},
	}

	rng := rand.New(rand.NewSource(2036))
	for i := 0; i < 20; i++ {
		n := randRange(rng, 3, 1_000_000_000_000_000_000)
		a, b, c := distinctTriple(rng, n)
		tests = append(tests, testCase{
			n:    n,
			a:    a,
			b:    b,
			c:    c,
			desc: fmt.Sprintf("random-%d", i+1),
		})
	}
	return tests
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}

func distinctTriple(rng *rand.Rand, n int64) (int64, int64, int64) {
	if n < 3 {
		return 1, 1, 1
	}
	choose := func() int64 { return randRange(rng, 1, n) }
	a := choose()
	var b int64
	for {
		b = choose()
		if b != a {
			break
		}
	}
	var c int64
	for {
		c = choose()
		if c != a && c != b {
			break
		}
	}
	return a, b, c
}
