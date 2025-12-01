package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./521E.go"

type testCase struct {
	input string
}

type graph struct {
	n   int
	adj []map[int]struct{}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		g, err := parseGraph(tc.input)
		if err != nil {
			fail("failed to parse generated test %d: %v", i+1, err)
		}
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		refYes, err := parseDecision(refOut)
		if err != nil {
			fail("failed to parse reference output on test %d: %v\noutput:\n%s", i+1, err, refOut)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		candRes, err := parseCandidateOutput(candOut)
		if err != nil {
			fail("could not parse candidate output on test %d: %v\noutput:\n%s", i+1, err, candOut)
		}

		if refYes {
			if !candRes.yes {
				fail("test %d: expected YES but candidate printed NO\ninput:\n%s", i+1, tc.input)
			}
			if err := validateRoutes(g, candRes.routes); err != nil {
				fail("test %d: invalid routes: %v\ninput:\n%s\ncandidate output:\n%s", i+1, err, tc.input, candOut)
			}
		} else {
			if candRes.yes {
				fail("test %d: expected NO but candidate printed YES\ninput:\n%s\ncandidate output:\n%s", i+1, tc.input, candOut)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "521E-ref-*")
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
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseDecision(out string) (bool, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return false, errors.New("missing verdict")
	}
	switch strings.ToUpper(sc.Text()) {
	case "YES":
		return true, nil
	case "NO":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected verdict %q", sc.Text())
	}
}

type candidateResult struct {
	yes    bool
	routes [][]int
}

func parseCandidateOutput(out string) (candidateResult, error) {
	res := candidateResult{}
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	sc.Buffer(make([]byte, 1024), 1<<20)
	if !sc.Scan() {
		return res, errors.New("missing verdict")
	}
	verdict := strings.ToUpper(sc.Text())
	if verdict == "NO" {
		return res, nil
	}
	if verdict != "YES" {
		return res, fmt.Errorf("unexpected verdict %q", verdict)
	}
	res.yes = true
	res.routes = make([][]int, 3)
	for i := 0; i < 3; i++ {
		if !sc.Scan() {
			return res, fmt.Errorf("missing length for route %d", i+1)
		}
		l, err := strconv.Atoi(sc.Text())
		if err != nil || l <= 0 {
			return res, fmt.Errorf("invalid length for route %d", i+1)
		}
		route := make([]int, l)
		for j := 0; j < l; j++ {
			if !sc.Scan() {
				return res, fmt.Errorf("route %d missing node %d", i+1, j+1)
			}
			val, err := strconv.Atoi(sc.Text())
			if err != nil {
				return res, fmt.Errorf("route %d node %d is not an integer", i+1, j+1)
			}
			route[j] = val
		}
		res.routes[i] = route
	}
	return res, nil
}

func parseGraph(input string) (*graph, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return nil, err
	}
	g := &graph{
		n:   n,
		adj: make([]map[int]struct{}, n+1),
	}
	for i := 0; i < m; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return nil, err
		}
		if u < 1 || u > n || v < 1 || v > n || u == v {
			return nil, fmt.Errorf("invalid edge %d-%d", u, v)
		}
		if g.adj[u] == nil {
			g.adj[u] = make(map[int]struct{})
		}
		if g.adj[v] == nil {
			g.adj[v] = make(map[int]struct{})
		}
		g.adj[u][v] = struct{}{}
		g.adj[v][u] = struct{}{}
	}
	return g, nil
}

func (g *graph) hasEdge(u, v int) bool {
	if u < 1 || u > g.n || v < 1 || v > g.n {
		return false
	}
	if g.adj[u] == nil {
		return false
	}
	_, ok := g.adj[u][v]
	return ok
}

func edgeKey(u, v int) uint64 {
	if u > v {
		u, v = v, u
	}
	return uint64(u)<<32 | uint64(v)
}

func validateRoutes(g *graph, routes [][]int) error {
	if len(routes) != 3 {
		return errors.New("expected three routes")
	}
	start := routes[0][0]
	finish := routes[0][len(routes[0])-1]
	if start == finish {
		return errors.New("start and finish must differ")
	}
	if start < 1 || start > g.n || finish < 1 || finish > g.n {
		return errors.New("start or finish outside graph")
	}
	internalUsed := make(map[int]struct{})
	edgesUsed := make(map[uint64]struct{})

	for idx, route := range routes {
		if len(route) < 2 {
			return fmt.Errorf("route %d is too short", idx+1)
		}
		if route[0] != start || route[len(route)-1] != finish {
			return fmt.Errorf("route %d does not start/end at common vertices", idx+1)
		}
		seen := make(map[int]struct{}, len(route))
		for i, node := range route {
			if node < 1 || node > g.n {
				return fmt.Errorf("route %d uses invalid node %d", idx+1, node)
			}
			if _, ok := seen[node]; ok {
				return fmt.Errorf("route %d visits node %d twice", idx+1, node)
			}
			seen[node] = struct{}{}
			if i > 0 {
				prev := route[i-1]
				if !g.hasEdge(prev, node) {
					return fmt.Errorf("route %d uses missing edge %d-%d", idx+1, prev, node)
				}
				key := edgeKey(prev, node)
				if _, ok := edgesUsed[key]; ok {
					return fmt.Errorf("edge %d-%d reused across routes", prev, node)
				}
				edgesUsed[key] = struct{}{}
			}
			if i != 0 && i != len(route)-1 {
				if _, ok := internalUsed[node]; ok {
					return fmt.Errorf("node %d appears in multiple routes", node)
				}
				internalUsed[node] = struct{}{}
			}
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20240602))
	var tests []testCase
	tests = append(tests, manualNoCase())
	tests = append(tests, manualYesCase())
	tests = append(tests, triangleCase())
	tests = append(tests, triplePathCase())

	for i := 0; i < 40; i++ {
		n := rng.Intn(18) + 2
		maxEdges := int64(n*(n-1)) / 2
		m := rng.Intn(n*2) + 1
		if int64(m) > maxEdges {
			m = int(maxEdges)
		}
		tests = append(tests, randomGraphTest(rng, n, m))
	}

	for i := 0; i < 25; i++ {
		n := rng.Intn(400) + 50
		m := rng.Intn(3*n) + n
		tests = append(tests, randomGraphTest(rng, n, m))
	}

	tests = append(tests, randomGraphTest(rng, 5000, 8000))
	tests = append(tests, randomGraphTest(rng, 20000, 40000))
	tests = append(tests, randomGraphTest(rng, 80000, 150000))
	tests = append(tests, randomGraphTest(rng, 200000, 200000))
	tests = append(tests, largeGuaranteedYesCase())

	return tests
}

func manualNoCase() testCase {
	edges := [][2]int{
		{1, 2},
		{2, 3},
	}
	return makeTestCaseInput(3, edges)
}

func manualYesCase() testCase {
	edges := [][2]int{
		{1, 2}, {2, 5}, {5, 9},
		{1, 3}, {3, 6}, {6, 9},
		{1, 4}, {4, 7}, {7, 9},
		{2, 6}, {3, 7}, {4, 5},
	}
	return makeTestCaseInput(10, edges)
}

func triangleCase() testCase {
	edges := [][2]int{
		{1, 2}, {2, 3}, {3, 1},
		{3, 4}, {4, 5},
	}
	return makeTestCaseInput(5, edges)
}

func triplePathCase() testCase {
	edges := [][2]int{
		{1, 2}, {2, 3}, {3, 10},
		{1, 4}, {4, 5}, {5, 10},
		{1, 6}, {6, 7}, {7, 10},
		{3, 8}, {8, 9}, {9, 10},
	}
	return makeTestCaseInput(10, edges)
}

func largeGuaranteedYesCase() testCase {
	var edges [][2]int
	finish := 400
	for i := 2; i <= 5; i++ {
		edges = append(edges, [2]int{i - 1, i})
	}
	edges = append(edges, [2]int{5, finish})

	path2 := []int{1, 101, 102, 103, 104, finish}
	for i := 1; i < len(path2); i++ {
		edges = append(edges, [2]int{path2[i-1], path2[i]})
	}
	path3 := []int{1, 201, 202, 203, 204, 205, finish}
	for i := 1; i < len(path3); i++ {
		edges = append(edges, [2]int{path3[i-1], path3[i]})
	}
	for i := 206; i <= 390; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return makeTestCaseInput(450, edges)
}

func randomGraphTest(rng *rand.Rand, n, m int) testCase {
	maxEdges := int64(n*(n-1)) / 2
	if int64(m) > maxEdges {
		m = int(maxEdges)
	}
	if m <= 0 {
		m = 1
	}
	edges := make(map[[2]int]struct{})
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		if u > v {
			u, v = v, u
		}
		edges[[2]int{u, v}] = struct{}{}
	}
	list := make([][2]int, 0, len(edges))
	for e := range edges {
		list = append(list, e)
	}
	return makeTestCaseInput(n, list)
}

func makeTestCaseInput(n int, edges [][2]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return testCase{input: sb.String()}
}
