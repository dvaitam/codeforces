package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSource       = "2000-2999/2000-2099/2060-2069/2061/2061H1.go"
	maxTotalN       = 1800
	maxTotalM       = 8000
	defaultTCases   = 120
	maxNodesPerTest = 2000
)

type testCase struct {
	n     int
	m     int
	s     string
	t     string
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %v, got %v\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2061H1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func parseAnswers(out string, t int) ([]bool, error) {
	tokens := strings.Fields(out)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("too many tokens in output: expected %d, got %d", t, len(tokens))
	}
	ans := make([]bool, t)
	for i, tok := range tokens {
		low := strings.ToLower(tok)
		switch low[0] {
		case 'y':
			ans[i] = true
		case 'n':
			ans[i] = false
		default:
			return nil, fmt.Errorf("invalid answer token %q at position %d", tok, i+1)
		}
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		fmt.Fprintf(&b, "%s\n", tc.s)
		fmt.Fprintf(&b, "%s\n", tc.t)
		for _, e := range tc.edges {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0
	totalM := 0

	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN || totalM+tc.m > maxTotalM {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalM += tc.m
	}

	// Simple manual cases.
	add(parseManualSample(
		2, 1, "10", "01",
		[][2]int{{1, 2}},
	))
	add(parseManualSample(
		1, 0, "1", "1",
		nil,
	))
	add(parseManualSample(
		4, 2, "1100", "0011",
		[][2]int{{1, 2}, {3, 4}},
	))

	// Trivial cases.
	add(randomCase(1, 0, 0, rng))
	add(randomCase(1, 0, 1, rng))
	add(randomCase(2, 1, 1, rng))

	// Randomized cases.
	targetTests := rng.Intn(defaultTCases/2) + defaultTCases/2
	for len(tests) < targetTests {
		n := rng.Intn(maxNodesPerTest-1) + 1
		m := rng.Intn(min(n*(n-1)/2, 2000)) // keep single test reasonable.
		if n == 1 {
			m = 0
		}
		ones := rng.Intn(n + 1)
		tc := randomCase(n, m, ones, rng)
		add(tc)
	}

	if len(tests) == 0 {
		add(randomCase(3, 2, 1, rng))
	}
	return tests
}

func parseManualSample(n, m int, s, t string, edges [][2]int) testCase {
	return testCase{n: n, m: m, s: s, t: t, edges: edges}
}

func randomCase(n, m, ones int, rng *rand.Rand) testCase {
	degLimit := n * (n - 1) / 2
	if m > degLimit {
		m = degLimit
	}

	edges := make([][2]int, 0, m)
	seen := make(map[[2]int]struct{})
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		edges = append(edges, key)
	}

	if ones > n {
		ones = n
	}
	sBits := make([]byte, n)
	tBits := make([]byte, n)
	for i := 0; i < n; i++ {
		sBits[i] = '0'
		tBits[i] = '0'
	}
	places := rng.Perm(n)
	for i := 0; i < ones; i++ {
		sBits[places[i]] = '1'
	}
	places = rng.Perm(n)
	for i := 0; i < ones; i++ {
		tBits[places[i]] = '1'
	}

	return testCase{
		n:     n,
		m:     len(edges),
		s:     string(sBits),
		t:     string(tBits),
		edges: edges,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
