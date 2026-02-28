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

const targetTests = 80

type testCase struct {
	n int
	k int
	s string
	q int
	t []string
}

// solveE implements the correct reference solution using next-occurrence
// arrays and DP.  For each position i in s (and position 0 = "before s"),
// nxt[i][c] = next occurrence of character c after position i.
// dp[i] = minimum letters to append so that starting from position i the
// string becomes unpleasant (not a subsequence). dp[n+1] = 0 (past end).
// dp[i] = 1 + min over c of dp[nxt[i][c]].
func solveE(tc testCase) []int64 {
	n := tc.n
	k := tc.k
	s := tc.s

	// Build nxt table: (n+1) positions x k characters.
	// nxt[i*k+c] = next position of char c at or after position i+1 (1-indexed in s).
	// If not found, n+1 (sentinel).
	nxt := make([]int32, (n+1)*k)
	nextPos := make([]int32, k)
	for c := 0; c < k; c++ {
		nextPos[c] = int32(n + 1)
	}
	for i := n; i >= 0; i-- {
		base := i * k
		for c := 0; c < k; c++ {
			nxt[base+c] = nextPos[c]
		}
		if i > 0 {
			ch := int(s[i-1] - 'a')
			nextPos[ch] = int32(i)
		}
	}

	// Build DP.
	dp := make([]int32, n+2)
	for i := n; i >= 0; i-- {
		base := i * k
		mn := int32(1 << 30)
		for c := 0; c < k; c++ {
			j := int(nxt[base+c])
			if dp[j] < mn {
				mn = dp[j]
			}
		}
		dp[i] = mn + 1
	}

	// Answer queries.
	results := make([]int64, tc.q)
	for qi, t := range tc.t {
		pos := 0
		for _, bb := range t {
			c := int(bb - 'a')
			pos = int(nxt[pos*k+c])
			if pos == n+1 {
				break
			}
		}
		if pos == n+1 {
			results[qi] = 0
		} else {
			results[qi] = int64(dp[pos])
		}
	}
	return results
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()

	for idx, tc := range tests {
		input := formatSingleTest(tc)
		expected := solveE(tc)

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate runtime error: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}
		candAns, err := parseInts(candOut, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse candidate output: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := range expected {
			if expected[i] != candAns[i] {
				fmt.Fprintf(os.Stderr, "test %d, query %d: expected %d, got %d\n", idx+1, i+1, expected[i], candAns[i])
				fmt.Fprintf(os.Stderr, "input:\n%s", input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Accepted (%d test cases).\n", len(tests))
}

func formatSingleTest(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
	fmt.Fprintf(&b, "%s\n", tc.s)
	fmt.Fprintf(&b, "%d\n", tc.q)
	for _, t := range tc.t {
		fmt.Fprintf(&b, "%s\n", t)
	}
	return b.String()
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	return out.String(), nil
}

func parseInts(out string, expected int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d: %v", tok, i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	add := func(tc testCase) {
		tests = append(tests, tc)
	}

	// Manual cases from statement.
	add(testCase{
		n: 7, k: 3,
		s: "abacaba",
		q: 3,
		t: []string{"cc", "bcb", "b"},
	})
	add(testCase{
		n: 5, k: 1,
		s: "aaaaa",
		q: 6,
		t: []string{"a", "aa", "aaa", "aaaa", "aaaaa", "aaaaaa"},
	})

	// Additional small crafted cases.
	add(testCase{
		n: 3, k: 2,
		s: "aba",
		q: 3,
		t: []string{"a", "b", "ab"},
	})
	add(testCase{
		n: 1, k: 1,
		s: "a",
		q: 2,
		t: []string{"a", "aa"},
	})

	for len(tests) < targetTests {
		n := rng.Intn(5000) + 1
		k := rng.Intn(26) + 1
		s := randString(rng, n, k)
		q := rng.Intn(500) + 1
		qs := make([]string, q)
		for i := 0; i < q; i++ {
			maxLen := min(20, n+5)
			length := rng.Intn(maxLen) + 1
			qs[i] = randString(rng, length, k)
		}
		add(testCase{n: n, k: k, s: s, q: q, t: qs})
	}

	return tests
}

func randString(rng *rand.Rand, length int, k int) string {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = byte('a' + rng.Intn(k))
	}
	return string(b)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
