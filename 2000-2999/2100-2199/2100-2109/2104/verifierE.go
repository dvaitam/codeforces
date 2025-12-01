package main

import (
	"bytes"
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
	// refSource points to the local reference solution to avoid GOPATH resolution.
	refSource       = "2104E.go"
	maxTotalN       = 200000
	maxTotalQ       = 200000
	maxTotalTLength = 200000
	targetTests     = 80
)

type testCase struct {
	n int
	k int
	s string
	q int
	t []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	refAns, err := parseAnswers(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d values, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch at answer %d: expected %d, got %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d test cases).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2104E-ref-*")
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

func parseAnswers(out string, tests []testCase) ([]int64, error) {
	total := 0
	for _, tc := range tests {
		total += tc.q
	}
	tokens := strings.Fields(out)
	if len(tokens) != total {
		return nil, fmt.Errorf("expected %d answers, got %d", total, len(tokens))
	}
	res := make([]int64, total)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d: %v", tok, i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
		fmt.Fprintf(&b, "%s\n", tc.s)
		fmt.Fprintf(&b, "%d\n", tc.q)
		for _, t := range tc.t {
			fmt.Fprintf(&b, "%s\n", t)
		}
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN, totalQ, totalLenT := 0, 0, 0

	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN || totalQ+tc.q > maxTotalQ {
			return
		}
		sumLen := 0
		for _, s := range tc.t {
			sumLen += len(s)
		}
		if totalLenT+sumLen > maxTotalTLength {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalQ += tc.q
		totalLenT += sumLen
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

	for len(tests) < targetTests && totalN < maxTotalN && totalQ < maxTotalQ {
		remainN := maxTotalN - totalN
		if remainN <= 0 {
			break
		}
		n := rng.Intn(min(5000, remainN)) + 1
		k := rng.Intn(26) + 1
		s := randString(rng, n, k)

		remainQ := maxTotalQ - totalQ
		if remainQ <= 0 {
			break
		}
		q := rng.Intn(min(500, remainQ)) + 1
		qs := make([]string, q)
		sumLen := 0
		for i := 0; i < q; i++ {
			maxLen := min(20, n+5)
			length := rng.Intn(maxLen) + 1
			qs[i] = randString(rng, length, k)
			sumLen += length
		}
		if totalLenT+sumLen > maxTotalTLength {
			break
		}
		add(testCase{n: n, k: k, s: s, q: q, t: qs})
	}

	if len(tests) == 0 {
		add(testCase{n: 1, k: 1, s: "a", q: 1, t: []string{"a"}})
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
