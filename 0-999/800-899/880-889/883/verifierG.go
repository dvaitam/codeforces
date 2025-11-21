package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edgeInput struct {
	t, u, v int
}

type Edge struct {
	to     int
	undIdx int // -1 for directed edges
	dir    bool
}

type testCase struct {
	input    string
	adj      [][]Edge
	n, s     int
	undCnt   int
	maxCount int
	minCount int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(tc, out); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(tc testCase, output string) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("empty output")
	}
	tokens := strings.Fields(output)
	if len(tokens) != 4 {
		return fmt.Errorf("expected 4 tokens (cnt, orient, cnt, orient) but got %d", len(tokens))
	}
	maxCnt, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("invalid max reachable count: %v", err)
	}
	maxPlan, err := parsePlan(tokens[1], tc.undCnt)
	if err != nil {
		return fmt.Errorf("invalid max plan: %v", err)
	}
	minCnt, err := strconv.Atoi(tokens[2])
	if err != nil {
		return fmt.Errorf("invalid min reachable count: %v", err)
	}
	minPlan, err := parsePlan(tokens[3], tc.undCnt)
	if err != nil {
		return fmt.Errorf("invalid min plan: %v", err)
	}
	if maxCnt < 1 || maxCnt > tc.n {
		return fmt.Errorf("reported max reachable count %d out of range", maxCnt)
	}
	if minCnt < 1 || minCnt > tc.n {
		return fmt.Errorf("reported min reachable count %d out of range", minCnt)
	}
	actualMax := reachableWithPlan(tc.adj, tc.s, maxPlan)
	if actualMax != maxCnt {
		return fmt.Errorf("max plan reaches %d vertices, but reported %d", actualMax, maxCnt)
	}
	if actualMax != tc.maxCount {
		return fmt.Errorf("max plan reachability %d differs from optimal %d", actualMax, tc.maxCount)
	}
	actualMin := reachableWithPlan(tc.adj, tc.s, minPlan)
	if actualMin != minCnt {
		return fmt.Errorf("min plan reaches %d vertices, but reported %d", actualMin, minCnt)
	}
	if actualMin != tc.minCount {
		return fmt.Errorf("min plan reachability %d differs from optimal %d", actualMin, tc.minCount)
	}
	return nil
}

func parsePlan(s string, expected int) ([]bool, error) {
	if len(s) != expected {
		return nil, fmt.Errorf("orientation string length %d differs from %d", len(s), expected)
	}
	plan := make([]bool, expected)
	for i, ch := range s {
		switch ch {
		case '+':
			plan[i] = true
		case '-':
			plan[i] = false
		default:
			return nil, fmt.Errorf("invalid character %q in orientation string", ch)
		}
	}
	return plan, nil
}

func reachableWithPlan(adj [][]Edge, s int, plan []bool) int {
	visited := make([]bool, len(adj))
	queue := []int{s}
	visited[s] = true
	count := 1
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, e := range adj[u] {
			if e.undIdx == -1 {
				if !visited[e.to] {
					visited[e.to] = true
					queue = append(queue, e.to)
					count++
				}
			} else {
				if plan[e.undIdx] == e.dir && !visited[e.to] {
					visited[e.to] = true
					queue = append(queue, e.to)
					count++
				}
			}
		}
	}
	return count
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		buildTest(2, 1, []edgeInput{
			{t: 2, u: 1, v: 2},
		}),
		buildTest(3, 1, []edgeInput{
			{t: 1, u: 1, v: 2},
			{t: 2, u: 2, v: 3},
		}),
	}
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest())
	}
	return tests
}

func randomTest() testCase {
	n := rand.Intn(8) + 2
	m := rand.Intn(10) + 1
	s := rand.Intn(n) + 1
	edges := make([]edgeInput, 0, m)
	edges = append(edges, randomEdge(2, n))
	for len(edges) < m {
		t := rand.Intn(2) + 1
		edges = append(edges, randomEdge(t, n))
	}
	return buildTest(n, s, edges)
}

func randomEdge(t, n int) edgeInput {
	u := rand.Intn(n) + 1
	v := rand.Intn(n-1) + 1
	if v >= u {
		v++
	}
	return edgeInput{t: t, u: u, v: v}
}

func buildTest(n, s int, edges []edgeInput) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, len(edges), s)
	undCnt := 0
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.t, e.u, e.v)
		if e.t == 2 {
			undCnt++
		}
	}
	if undCnt == 0 {
		panic("test without undirected edges")
	}
	adj := make([][]Edge, n+1)
	idx := 0
	for _, e := range edges {
		if e.t == 1 {
			adj[e.u] = append(adj[e.u], Edge{to: e.v, undIdx: -1})
		} else {
			adj[e.u] = append(adj[e.u], Edge{to: e.v, undIdx: idx, dir: true})
			adj[e.v] = append(adj[e.v], Edge{to: e.u, undIdx: idx, dir: false})
			idx++
		}
	}
	maxCount, minCount := computeCounts(adj, n, s, undCnt)
	return testCase{
		input:    sb.String(),
		adj:      adj,
		n:        n,
		s:        s,
		undCnt:   undCnt,
		maxCount: maxCount,
		minCount: minCount,
	}
}

func computeCounts(adj [][]Edge, n, s, undCnt int) (int, int) {
	visDirected := make([]bool, n+1)
	stack := []int{s}
	visDirected[s] = true
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range adj[u] {
			if e.undIdx == -1 && !visDirected[e.to] {
				visDirected[e.to] = true
				stack = append(stack, e.to)
			}
		}
	}
	minCount := 0
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if visDirected[i] {
			minCount++
			queue = append(queue, i)
		}
	}
	visMax := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if visDirected[i] {
			visMax[i] = true
		}
	}
	used := make([]bool, undCnt)
	maxCount := minCount
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, e := range adj[u] {
			if e.undIdx == -1 {
				if !visMax[e.to] {
					visMax[e.to] = true
					queue = append(queue, e.to)
					maxCount++
				}
			} else {
				idx := e.undIdx
				if used[idx] {
					continue
				}
				if !visMax[e.to] {
					used[idx] = true
					visMax[e.to] = true
					queue = append(queue, e.to)
					maxCount++
				}
			}
		}
	}
	return maxCount, minCount
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
