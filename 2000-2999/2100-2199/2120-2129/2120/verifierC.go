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

type testCase struct {
	n int
	m int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2120C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "2120C.go")
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
	sb.Grow(len(tests) * 32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	}
	return sb.String()
}

type parsedTree struct {
	possible bool
	root     int
	edges    [][2]int
}

func parseOutput(out string, tests []testCase, fromOracle bool) ([]parsedTree, error) {
	tokens := strings.Fields(out)
	res := make([]parsedTree, len(tests))
	idx := 0
	for i, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("output ended early at test %d", i+1)
		}
		if tokens[idx] == "-1" {
			res[i] = parsedTree{possible: false}
			idx++
			continue
		}
		root, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid root %q", i+1, tokens[idx])
		}
		idx++
		edges := make([][2]int, 0, tc.n-1)
		for j := 0; j < tc.n-1; j++ {
			if idx+1 >= len(tokens) {
				return nil, fmt.Errorf("test %d: insufficient edge data", i+1)
			}
			u, err1 := strconv.Atoi(tokens[idx])
			v, err2 := strconv.Atoi(tokens[idx+1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("test %d: invalid edge endpoints %q %q", i+1, tokens[idx], tokens[idx+1])
			}
			idx += 2
			edges = append(edges, [2]int{u, v})
		}
		res[i] = parsedTree{possible: true, root: root, edges: edges}
	}
	if !fromOracle && idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens detected after parsing output")
	}
	return res, nil
}

func validateTree(tc testCase, pt parsedTree) error {
	if !pt.possible {
		return nil
	}
	if pt.root < 1 || pt.root > tc.n {
		return fmt.Errorf("root %d out of range", pt.root)
	}
	if len(pt.edges) != tc.n-1 {
		return fmt.Errorf("edge count %d does not match n-1", len(pt.edges))
	}
	adj := make([][]int, tc.n+1)
	for i, e := range pt.edges {
		u, v := e[0], e[1]
		if u < 1 || u > tc.n || v < 1 || v > tc.n {
			return fmt.Errorf("edge %d endpoint out of range (%d,%d)", i+1, u, v)
		}
		if u == v {
			return fmt.Errorf("edge %d is a self-loop at %d", i+1, u)
		}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// connectivity
	stack := []int{pt.root}
	visited := make([]bool, tc.n+1)
	visited[pt.root] = true
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range adj[v] {
			if !visited[to] {
				visited[to] = true
				stack = append(stack, to)
			}
		}
	}
	for i := 1; i <= tc.n; i++ {
		if !visited[i] {
			return fmt.Errorf("graph is disconnected, node %d unreachable", i)
		}
	}
	// check for cycles via edges count and connectivity already ensures tree.
	// compute divineness
	parent := make([]int, tc.n+1)
	minLabel := make([]int, tc.n+1)
	stack = []int{pt.root}
	parent[pt.root] = -1
	minLabel[pt.root] = pt.root
	order := []int{}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			if parent[to] != 0 {
				return fmt.Errorf("cycle detected between %d and %d", v, to)
			}
			parent[to] = v
			if vLabel := minLabel[v]; vLabel < to {
				minLabel[to] = vLabel
			} else {
				minLabel[to] = to
			}
			stack = append(stack, to)
		}
	}
	if len(order) != tc.n {
		return fmt.Errorf("tree traversal missing nodes")
	}

	var total int64
	for i := 1; i <= tc.n; i++ {
		total += int64(minLabel[i])
	}
	if total != tc.m {
		return fmt.Errorf("divineness sum mismatch: got %d, expected %d", total, tc.m)
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, m: 1},
		{n: 1, m: 2},  // impossible
		{n: 4, m: 6},  // sample possible
		{n: 3, m: 5},  // small possible
		{n: 3, m: 10}, // impossible large
		{n: 5, m: 15}, // chain minimal sum
		{n: 5, m: 20}, // mid value
		{n: 5, m: 25}, // max sum
	}
}

func randomTests(totalN int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 256)
	used := 0
	for used < totalN {
		remain := totalN - used
		n := rng.Intn(min(5000, remain)) + 1
		// possible sums range between min=sum of i?(div) ? minimal sum when root n? Wait minimal sum = n (root=1 all children 1)?? Actually minimal sum = n (all div=1) yes max sum? chain inc sum of 1..n. choose random within.
		minSum := int64(n) // if root 1 all parents 1
		maxSum := int64(n*(n+1)) / 2
		var m int64
		if rng.Intn(4) == 0 {
			// sometimes impossible outside range
			if rng.Intn(2) == 0 {
				m = minSum - int64(rng.Intn(n)+1)
			} else {
				m = maxSum + int64(rng.Intn(n)+1)
			}
		} else {
			delta := maxSum - minSum
			if delta < 0 {
				delta = 0
			}
			if delta == 0 {
				m = minSum
			} else {
				m = minSum + int64(rng.Int63n(delta+1))
			}
		}
		if m < 1 {
			m = 1
		}
		tests = append(tests, testCase{n: n, m: m})
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	const nLimit = 1_000_000
	used := totalN(tests)
	if used < nLimit {
		tests = append(tests, randomTests(nLimit-used)...)
	}

	input := buildInput(tests)

	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	expected, err := parseOutput(oracleOut, tests, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, oracleOut)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actual, err := parseOutput(actOut, tests, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\noutput:\n%s", err, actOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if expected[i].possible != actual[i].possible {
			fmt.Fprintf(os.Stderr, "test %d mismatch possible flag (expected %v, got %v)\ninput:\n%s", i+1, expected[i].possible, actual[i].possible, input)
			os.Exit(1)
		}
		if actual[i].possible {
			if err := validateTree(tc, actual[i]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d invalid tree: %v\ninput:\n%s", i+1, err, input)
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
