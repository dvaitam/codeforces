package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u int
	v int
}

type testCase struct {
	n     int
	vols  []int
	pits  []int
	edges []edge
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2110E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "2110E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}
	sb.Grow(totalN*16 + len(tests)*24)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d %d", tc.vols[i], tc.pits[i]))
			if i+1 < tc.n {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOracleResult(out string, tests []testCase) ([]bool, error) {
	tokens := strings.Fields(out)
	res := make([]bool, len(tests))
	idx := 0
	for i := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("oracle output ended early at test %d", i+1)
		}
		token := strings.ToLower(tokens[idx])
		idx++
		if token == "yes" {
			res[i] = true
		} else if token == "no" {
			res[i] = false
		} else {
			return nil, fmt.Errorf("oracle test %d: invalid token %q", i+1, token)
		}
		if res[i] {
			// consume the trail
			for j := 0; j < tests[i].n; j++ {
				if idx >= len(tokens) {
					return nil, fmt.Errorf("oracle test %d: missing edge index", i+1)
				}
				idx++
			}
		}
	}
	return res, nil
}

func parseCandidate(out string, tests []testCase) ([]bool, [][]int, error) {
	tokens := strings.Fields(out)
	ans := make([]bool, len(tests))
	orders := make([][]int, len(tests))
	idx := 0
	for i, tc := range tests {
		if idx >= len(tokens) {
			return nil, nil, fmt.Errorf("output ended early at test %d", i+1)
		}
		token := strings.ToLower(tokens[idx])
		idx++
		if token == "yes" {
			ans[i] = true
		} else if token == "no" {
			ans[i] = false
		} else {
			return nil, nil, fmt.Errorf("test %d: invalid yes/no token %q", i+1, token)
		}
		if ans[i] {
			order := make([]int, tc.n)
			for j := 0; j < tc.n; j++ {
				if idx >= len(tokens) {
					return nil, nil, fmt.Errorf("test %d: missing edge index %d", i+1, j+1)
				}
				v, err := strconv.Atoi(tokens[idx])
				idx++
				if err != nil || v < 1 || v > tc.n {
					return nil, nil, fmt.Errorf("test %d: invalid edge id %q", i+1, tokens[idx-1])
				}
				order[j] = v - 1
			}
			orders[i] = order
		}
	}
	if idx != len(tokens) {
		return nil, nil, fmt.Errorf("extra tokens after processing all tests")
	}
	return ans, orders, nil
}

func compress(t testCase) testCase {
	volMap := make(map[int]int)
	pitMap := make(map[int]int)
	for i := 0; i < t.n; i++ {
		if _, ok := volMap[t.vols[i]]; !ok {
			volMap[t.vols[i]] = len(volMap)
		}
		if _, ok := pitMap[t.pits[i]]; !ok {
			pitMap[t.pits[i]] = len(pitMap)
		}
	}
	volCnt := len(volMap)
	t.edges = make([]edge, t.n)
	for i := 0; i < t.n; i++ {
		u := volMap[t.vols[i]]
		v := pitMap[t.pits[i]] + volCnt
		t.edges[i] = edge{u: u, v: v}
	}
	return t
}

func validateTrail(tc testCase, order []int) error {
	if len(order) != tc.n {
		return fmt.Errorf("trail length %d does not match edge count %d", len(order), tc.n)
	}
	seen := make([]bool, tc.n)
	for _, id := range order {
		if seen[id] {
			return fmt.Errorf("edge %d used multiple times", id+1)
		}
		seen[id] = true
	}

	curPossible := make(map[int]struct{}, 2)
	first := tc.edges[order[0]]
	curPossible[first.u] = struct{}{}
	curPossible[first.v] = struct{}{}

	for i := 1; i < len(order); i++ {
		e := tc.edges[order[i]]
		nextPossible := make(map[int]struct{}, 2)
		for node := range curPossible {
			if node == e.u {
				nextPossible[e.v] = struct{}{}
			}
			if node == e.v {
				nextPossible[e.u] = struct{}{}
			}
		}
		if len(nextPossible) == 0 {
			return fmt.Errorf("edges %d and %d are not connected", order[i-1]+1, order[i]+1)
		}
		curPossible = nextPossible
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, vols: []int{1}, pits: []int{2}},
		{n: 2, vols: []int{1, 1}, pits: []int{2, 3}},
		{n: 3, vols: []int{1, 2, 3}, pits: []int{4, 5, 6}},       // disconnected nodes but edges isolated => possible? no euler trail
		{n: 3, vols: []int{1, 2, 1}, pits: []int{2, 3, 3}},       // euler path
		{n: 4, vols: []int{1, 1, 2, 2}, pits: []int{3, 4, 3, 4}}, // euler circuit
	}
}

func randomTests(totalN int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 128)
	used := 0
	for used < totalN {
		remain := totalN - used
		n := rng.Intn(min(5000, remain)) + 1
		vols := make([]int, n)
		pits := make([]int, n)
		for i := 0; i < n; i++ {
			vols[i] = rng.Intn(1000000)
			pits[i] = rng.Intn(1000000)
		}
		tests = append(tests, testCase{n: n, vols: vols, pits: pits})
		used += n
	}
	return tests
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	const nLimit = 100_000
	used := totalN(tests)
	if used < nLimit {
		tests = append(tests, randomTests(nLimit-used)...)
	}

	compTests := make([]testCase, len(tests))
	for i, tc := range tests {
		compTests[i] = compress(tc)
	}

	input := buildInput(tests)

	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	expectedYes, err := parseOracleResult(oracleOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, oracleOut)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	candYes, orders, err := parseCandidate(actOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\noutput:\n%s", err, actOut)
		os.Exit(1)
	}

	for i, tc := range compTests {
		if candYes[i] != expectedYes[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %v, got %v\ninput:\n%s", i+1, expectedYes[i], candYes[i], input)
			os.Exit(1)
		}
		if candYes[i] {
			if err := validateTrail(tc, orders[i]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d invalid trail: %v\ninput:\n%s", i+1, err, input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
