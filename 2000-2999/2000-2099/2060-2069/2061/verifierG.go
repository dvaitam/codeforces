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
	refSource       = "2000-2999/2000-2099/2060-2069/2061/2061G.go"
	maxTotalN       = 4000
	randomTests     = 50
	maxNPerRandom   = 160
	denseEdgeChance = 0.5
)

type testCase struct {
	n     int
	edges []byte // length n*(n-1)/2, order (1,2),(1,3),...,(n-1,n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)
	expectedK := make([]int, len(tests))
	for i, tc := range tests {
		expectedK[i] = (tc.n + 1) / 3
	}

	// Sanity-check reference output on generated tests.
	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference runtime error: %v\n%s", err, refOut)
	}
	if err := parseAndValidate(refOut, tests, expectedK); err != nil {
		fail("reference output invalid: %v\n%s", err, refOut)
	}

	candCmd := commandFor(candidate)
	candOut, err := runProgram(candCmd, input)
	if err != nil {
		fail("candidate runtime error: %v\n%s", err, candOut)
	}
	if err := parseAndValidate(candOut, tests, expectedK); err != nil {
		fail("wrong answer: %v\n%s", err, candOut)
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2061G-ref-*")
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

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
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

func parseAndValidate(output string, tests []testCase, expectedK []int) error {
	tokens := strings.Fields(output)
	pos := 0
	for caseIdx, tc := range tests {
		if pos >= len(tokens) {
			return fmt.Errorf("test %d: missing k", caseIdx+1)
		}
		k, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return fmt.Errorf("test %d: invalid k token %q", caseIdx+1, tokens[pos])
		}
		pos++
		if k < 1 || k > tc.n/2 {
			return fmt.Errorf("test %d: k=%d out of range", caseIdx+1, k)
		}
		if expectedK != nil && k != expectedK[caseIdx] {
			return fmt.Errorf("test %d: expected k=%d, got %d", caseIdx+1, expectedK[caseIdx], k)
		}

		need := 2 * k
		if pos+need > len(tokens) {
			return fmt.Errorf("test %d: expected %d nodes, only %d tokens left", caseIdx+1, need, len(tokens)-pos)
		}

		used := make([]bool, tc.n+1)
		var relation byte
		haveRelation := false

		for i := 0; i < k; i++ {
			u, err1 := strconv.Atoi(tokens[pos])
			v, err2 := strconv.Atoi(tokens[pos+1])
			pos += 2
			if err1 != nil || err2 != nil {
				return fmt.Errorf("test %d pair %d: invalid node ids", caseIdx+1, i+1)
			}
			if u < 1 || u > tc.n || v < 1 || v > tc.n || u == v {
				return fmt.Errorf("test %d pair %d: nodes out of range or equal (%d, %d)", caseIdx+1, i+1, u, v)
			}
			if used[u] || used[v] {
				return fmt.Errorf("test %d: node reused in multiple pairs", caseIdx+1)
			}
			used[u], used[v] = true, true

			cur := edgeValue(tc, u, v)
			if !haveRelation {
				relation = cur
				haveRelation = true
			} else if cur != relation {
				return fmt.Errorf("test %d: pairs are not uniform (mixed friends and non-friends)", caseIdx+1)
			}
		}
	}

	if pos != len(tokens) {
		return fmt.Errorf("unexpected extra tokens in output (got %d extra)", len(tokens)-pos)
	}

	return nil
}

func edgeValue(tc testCase, u, v int) byte {
	if u > v {
		u, v = v, u
	}
	// Offset contributed by vertices before u plus position within u's row.
	offset := (u - 1) * (2*tc.n - u) / 2
	index := offset + (v - u - 1)
	return tc.edges[index]
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "manual\n%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		b.Write(tc.edges)
		b.WriteByte('\n')
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Deterministic small cases to sanity-check parsing and both relationship types.
	add(fullGraph(2, '1'))
	add(fullGraph(3, '0'))
	add(customStarCase())

	// Alternating pattern to make sure indexing is correct.
	add(patternCase(6))
	add(patternCase(9))

	for len(tests) < randomTests && totalN < maxTotalN {
		maxN := maxNPerRandom
		if rem := maxTotalN - totalN; maxN > rem {
			maxN = rem
		}
		if maxN < 2 {
			break
		}
		n := rng.Intn(maxN-1) + 2
		tc := randomCase(n, rng)
		add(tc)
	}

	return tests
}

func fullGraph(n int, val byte) testCase {
	size := n * (n - 1) / 2
	return testCase{n: n, edges: bytes.Repeat([]byte{val}, size)}
}

func customStarCase() testCase {
	// Friends are only edges from node 1 to others; enemies elsewhere.
	n := 7
	size := n * (n - 1) / 2
	edges := make([]byte, size)
	pos := 0
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if i == 1 {
				edges[pos] = '1'
			} else {
				edges[pos] = '0'
			}
			pos++
		}
	}
	return testCase{n: n, edges: edges}
}

func patternCase(n int) testCase {
	size := n * (n - 1) / 2
	edges := make([]byte, size)
	pos := 0
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if (i+j)%2 == 0 {
				edges[pos] = '1'
			} else {
				edges[pos] = '0'
			}
			pos++
		}
	}
	return testCase{n: n, edges: edges}
}

func randomCase(n int, rng *rand.Rand) testCase {
	size := n * (n - 1) / 2
	edges := make([]byte, size)
	for i := 0; i < size; i++ {
		if rng.Float64() < denseEdgeChance {
			edges[i] = '1'
		} else {
			edges[i] = '0'
		}
	}
	return testCase{n: n, edges: edges}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
