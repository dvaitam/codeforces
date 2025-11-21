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
	refSource   = "2000-2999/2100-2199/2100-2109/2110/2110D.go"
	maxTotalN   = 20000
	maxTotalM   = 30000
	targetTests = 120
	maxValue    = 1_000_000_000
)

type edge struct {
	u int
	v int
	w int64
}

type testCase struct {
	n int
	b []int64
	e []edge
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d, got %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d test cases).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2110D-ref-*")
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

func parseAnswers(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	ans := make([]int64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d: %v", tok, i+1, err)
		}
		ans[i] = val
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, len(tc.e))
		for i, v := range tc.b {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for _, e := range tc.e {
			fmt.Fprintf(&b, "%d %d %d\n", e.u, e.v, e.w)
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
		if totalN+tc.n > maxTotalN || totalM+len(tc.e) > maxTotalM {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
		totalM += len(tc.e)
	}

	// Manual cases (from statement and small edges).
	add(testCase{
		n: 3,
		b: []int64{2, 0, 0},
		e: []edge{{1, 2, 1}, {2, 3, 1}, {1, 3, 2}},
	})
	add(testCase{
		n: 5,
		b: []int64{2, 2, 5, 0, 1},
		e: []edge{{1, 2, 2}, {1, 3, 1}, {1, 4, 3}, {3, 5, 5}, {2, 4, 4}, {4, 5, 3}},
	})
	add(testCase{
		n: 3,
		b: []int64{0, 1, 0},
		e: []edge{{1, 2, 1}},
	})
	add(testCase{
		n: 4,
		b: []int64{1, 9, 0, 0},
		e: []edge{{1, 2, 1}, {1, 3, 3}, {2, 4, 10}},
	})

	for len(tests) < targetTests && totalN < maxTotalN && totalM < maxTotalM {
		remainN := maxTotalN - totalN
		if remainN < 2 {
			break
		}
		n := rng.Intn(min(200, remainN-1)) + 2
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			if rng.Intn(100) < 10 {
				b[i] = 0
			} else {
				b[i] = rng.Int63n(maxValue) + 1
			}
		}

		guaranteePath := rng.Intn(100) < 70
		var edges []edge
		if guaranteePath {
			for i := 1; i < n; i++ {
				edges = append(edges, edge{u: i, v: i + 1, w: rng.Int63n(maxValue) + 1})
			}
		}

		possiblePairs := n * (n - 1) / 2
		maxExtra := min(3*n, possiblePairs-len(edges))
		extra := rng.Intn(maxExtra + 1)
		seen := make(map[[2]int]struct{})
		for _, e := range edges {
			seen[[2]int{e.u, e.v}] = struct{}{}
		}
		for i := 0; i < extra; i++ {
			u := rng.Intn(n-1) + 1
			v := rng.Intn(n-u) + u + 1 // ensure u < v
			key := [2]int{u, v}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			w := rng.Int63n(maxValue) + 1
			edges = append(edges, edge{u: u, v: v, w: w})
		}

		add(testCase{n: n, b: b, e: edges})
	}

	if len(tests) == 0 {
		add(testCase{n: 2, b: []int64{1, 0}, e: []edge{{1, 2, 1}}})
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
