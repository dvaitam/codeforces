package main

import (
	"bufio"
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

const refSource = "1723GP2.go"

type testCase struct {
	input string
	info  testInfo
}

type testInfo struct {
	n   int
	adj [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierGP2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests, err := buildTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build tests:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refOut, err := runExecutable(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(refOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(candOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\nInput:\n%sCandidate output:\n%s\n", i+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1723GP2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutable(path, input string) (string, error) {
	cmd := exec.Command(path)
	return execute(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return execute(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func execute(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildTests() ([]testCase, error) {
	baseTests := []string{
		"2 1\n0 1 5\n1 1.0\n0\n",
		"4 4\n0 1 2\n1 2 3\n2 3 4\n3 0 5\n2 0.5\n0\n2\n",
	}
	var tests []testCase
	for _, input := range baseTests {
		info, err := prepareTest(input)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: input, info: info})
	}

	randomConfigs := []struct {
		n, m int
		seed int64
	}{
		{5, 6, 1},
		{8, 12, 2},
		{12, 20, 3},
		{20, 35, 4},
		{40, 70, 5},
		{80, 150, 6},
		{150, 300, 7},
		{300, 600, time.Now().UnixNano()},
	}

	for _, cfg := range randomConfigs {
		input := randomTest(cfg.n, cfg.m, cfg.seed)
		info, err := prepareTest(input)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: input, info: info})
	}
	return tests, nil
}

func randomTest(n, m int, seed int64) string {
	if n < 2 {
		n = 2
	}
	maxEdges := n * (n - 1) / 2
	if m < n-1 {
		m = n - 1
	}
	if m > maxEdges {
		m = maxEdges
	}

	r := rand.New(rand.NewSource(seed))
	type edge struct {
		u, v, w int
	}
	edges := make([]edge, 0, m)
	used := make(map[int64]struct{})

	addEdge := func(u, v, w int) {
		if u > v {
			u, v = v, u
		}
		key := int64(u)*int64(n) + int64(v)
		if _, ok := used[key]; ok {
			return
		}
		used[key] = struct{}{}
		edges = append(edges, edge{u, v, w})
	}

	for v := 1; v < n; v++ {
		u := r.Intn(v)
		w := r.Intn(100) + 1
		addEdge(u, v, w)
	}
	for len(edges) < m {
		u := r.Intn(n)
		v := r.Intn(n - 1)
		if v >= u {
			v++
		}
		w := r.Intn(100000) + 1
		addEdge(u, v, w)
	}

	setSize := r.Intn(min(50, n)) + 1
	kVal := r.Float64() * 10
	nodes := r.Perm(n)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	sb.WriteString(fmt.Sprintf("%d %.6f\n", setSize, kVal))
	for i := 0; i < setSize; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", nodes[i]))
	}
	return sb.String()
}

func prepareTest(input string) (testInfo, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return testInfo{}, fmt.Errorf("failed to read n and m: %v", err)
	}
	if n < 1 || m < n-1 {
		return testInfo{}, fmt.Errorf("invalid n or m")
	}
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v, w int
		if _, err := fmt.Fscan(reader, &u, &v, &w); err != nil {
			return testInfo{}, fmt.Errorf("failed to read edge %d: %v", i+1, err)
		}
		if u < 0 || u >= n || v < 0 || v >= n || u == v {
			return testInfo{}, fmt.Errorf("edge %d has invalid endpoints", i+1)
		}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var setSize int
	var k float64
	if _, err := fmt.Fscan(reader, &setSize, &k); err != nil {
		return testInfo{}, fmt.Errorf("failed to read set size and k: %v", err)
	}
	if setSize < 1 || setSize > min(50, n) {
		return testInfo{}, fmt.Errorf("invalid set size")
	}
	for i := 0; i < setSize; i++ {
		var node int
		if _, err := fmt.Fscan(reader, &node); err != nil {
			return testInfo{}, fmt.Errorf("failed to read S node %d: %v", i+1, err)
		}
		if node < 0 || node >= n {
			return testInfo{}, fmt.Errorf("node %d in S is out of range", i+1)
		}
	}
	return testInfo{n: n, adj: adj}, nil
}

func checkOutput(output string, info testInfo) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Buffer(make([]byte, 1024), 16*1024*1024)
	scanner.Split(bufio.ScanWords)

	nextInt := func() (int, error) {
		if !scanner.Scan() {
			return 0, fmt.Errorf("unexpected end of output")
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("expected integer, got %q", scanner.Text())
		}
		return val, nil
	}

	k, err := nextInt()
	if err != nil {
		return fmt.Errorf("failed to read subset count: %v", err)
	}
	if k < 1 || k > info.n {
		return fmt.Errorf("invalid subset count %d", k)
	}

	seen := make([]bool, info.n)
	subsetMark := make([]int, info.n)
	visitMark := make([]int, info.n)
	curMark := 0
	curVisit := 0

	for idx := 0; idx < k; idx++ {
		cnt, err := nextInt()
		if err != nil {
			return fmt.Errorf("failed to read size of subset %d: %v", idx+1, err)
		}
		if cnt < 1 {
			return fmt.Errorf("subset %d has non-positive size", idx+1)
		}
		nodes := make([]int, cnt)
		for i := 0; i < cnt; i++ {
			node, err := nextInt()
			if err != nil {
				return fmt.Errorf("failed to read node %d in subset %d: %v", i+1, idx+1, err)
			}
			if node < 0 || node >= info.n {
				return fmt.Errorf("node %d in subset %d out of range", node, idx+1)
			}
			if seen[node] {
				return fmt.Errorf("node %d appears multiple times", node)
			}
			seen[node] = true
			nodes[i] = node
		}
		curMark++
		for _, node := range nodes {
			subsetMark[node] = curMark
		}
		curVisit++
		queue := []int{nodes[0]}
		visitMark[nodes[0]] = curVisit
		visitedCount := 0
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			visitedCount++
			for _, to := range info.adj[v] {
				if subsetMark[to] == curMark && visitMark[to] != curVisit {
					visitMark[to] = curVisit
					queue = append(queue, to)
				}
			}
		}
		if visitedCount != len(nodes) {
			return fmt.Errorf("subset %d is not connected", idx+1)
		}
	}

	for v, ok := range seen {
		if !ok {
			return fmt.Errorf("node %d missing from partition", v)
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
