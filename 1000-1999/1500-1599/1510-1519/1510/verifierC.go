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
)

const refSource = "./1510C.go"

type testCase struct {
	n     int
	edges [][2]int
}

type edgePair struct {
	u int
	v int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refA, err := referenceAnswer(refBin, input)
		if err != nil {
			fail("reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
		}

		out, err := runCommand(commandFor(candidate), input)
		if err != nil {
			fail("runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, out)
		}

		edges, err := parseOutput(out, tc.n)
		if err != nil {
			fail("invalid output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, out)
		}

		if err := verifyAnswer(tc, refA, edges); err != nil {
			fail("wrong answer on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, out)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1510C-ref-*")
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

func referenceAnswer(bin, input string) (int, error) {
	out, err := runCommand(exec.Command(bin), input)
	if err != nil {
		return 0, err
	}
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return 0, fmt.Errorf("reference produced empty output")
	}
	a, err := strconv.Atoi(sc.Text())
	if err != nil {
		return 0, fmt.Errorf("failed to parse added edge count: %w", err)
	}
	return a, nil
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutput(out string, n int) ([]edgePair, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	readInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected end of output")
		}
		val, err := strconv.Atoi(sc.Text())
		if err != nil {
			return 0, err
		}
		return val, nil
	}

	a, err := readInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read number of edges: %w", err)
	}
	if a < 0 {
		return nil, fmt.Errorf("negative edge count")
	}

	edges := make([]edgePair, a)
	for i := 0; i < a; i++ {
		u, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("failed to read edge %d u: %w", i+1, err)
		}
		v, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("failed to read edge %d v: %w", i+1, err)
		}
		if u < 1 || u > n || v < 1 || v > n {
			return nil, fmt.Errorf("edge %d endpoint out of range", i+1)
		}
		if u == v {
			return nil, fmt.Errorf("self-loop in edge %d", i+1)
		}
		edges[i] = edgePair{u: u, v: v}
	}

	if sc.Scan() {
		return nil, fmt.Errorf("extraneous data at the end of output")
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return edges, nil
}

func verifyAnswer(tc testCase, refA int, add []edgePair) error {
	if len(add) != refA {
		return fmt.Errorf("expected %d edges, got %d", refA, len(add))
	}

	existing := make(map[string]struct{})
	for _, e := range tc.edges {
		existing[pairKey(e[0], e[1])] = struct{}{}
	}

	seen := make(map[string]struct{})
	for i, e := range add {
		key := pairKey(e.u, e.v)
		if _, ok := existing[key]; ok {
			return fmt.Errorf("edge %d duplicates existing graph edge (%d, %d)", i+1, e.u, e.v)
		}
		if _, ok := seen[key]; ok {
			return fmt.Errorf("edge %d repeats an earlier added edge (%d, %d)", i+1, e.u, e.v)
		}
		seen[key] = struct{}{}
	}

	all := make([][2]int, 0, len(tc.edges)+len(add))
	all = append(all, tc.edges...)
	for _, e := range add {
		all = append(all, [2]int{e.u, e.v})
	}

	bridges, err := cactusCheck(tc.n, all)
	if err != nil {
		return fmt.Errorf("result is not a cactus: %v", err)
	}
	if !bridgesFormMatching(tc.n, bridges, all) {
		return fmt.Errorf("resulting cactus is not strong (a vertex is incident to multiple bridges)")
	}
	return nil
}

func cactusCheck(n int, edges [][2]int) ([]bool, error) {
	if n == 0 {
		return nil, fmt.Errorf("empty graph")
	}
	m := len(edges)
	g := make([][]int, n+1)
	edgeKey := make(map[string]struct{})
	for i, e := range edges {
		if e[0] < 1 || e[0] > n || e[1] < 1 || e[1] > n {
			return nil, fmt.Errorf("edge %d has endpoint out of range", i+1)
		}
		if e[0] == e[1] {
			return nil, fmt.Errorf("edge %d is a loop", i+1)
		}
		key := pairKey(e[0], e[1])
		if _, ok := edgeKey[key]; ok {
			return nil, fmt.Errorf("multi-edge detected on (%d, %d)", e[0], e[1])
		}
		edgeKey[key] = struct{}{}
		g[e[0]] = append(g[e[0]], i)
		g[e[1]] = append(g[e[1]], i)
	}

	tin := make([]int, n+1)
	low := make([]int, n+1)
	var stack []int
	var timer int
	bridges := make([]bool, m)

	validateComponent := func(comp []int) error {
		if len(comp) == 0 {
			return nil
		}
		if len(comp) == 1 {
			return nil
		}
		deg := make(map[int]int)
		for _, id := range comp {
			u, v := edges[id][0], edges[id][1]
			deg[u]++
			deg[v]++
		}
		if len(deg) != len(comp) {
			return fmt.Errorf("biconnected component with %d edges and %d vertices", len(comp), len(deg))
		}
		for _, d := range deg {
			if d != 2 {
				return fmt.Errorf("non-cycle component detected")
			}
		}
		return nil
	}

	var dfs func(v, pe int) error
	dfs = func(v, pe int) error {
		timer++
		tin[v] = timer
		low[v] = timer
		for _, id := range g[v] {
			if id == pe {
				continue
			}
			to := edges[id][0] ^ edges[id][1] ^ v
			if tin[to] == 0 {
				stack = append(stack, id)
				if err := dfs(to, id); err != nil {
					return err
				}
				if low[to] < low[v] {
					low[v] = low[to]
				}
				if low[to] >= tin[v] {
					comp := make([]int, 0)
					for {
						last := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						comp = append(comp, last)
						if last == id {
							break
						}
					}
					if err := validateComponent(comp); err != nil {
						return err
					}
				}
				if low[to] > tin[v] {
					bridges[id] = true
				}
			} else if tin[to] < tin[v] {
				stack = append(stack, id)
				if tin[to] < low[v] {
					low[v] = tin[to]
				}
			}
		}
		return nil
	}

	if err := dfs(1, -1); err != nil {
		return nil, err
	}
	for v := 1; v <= n; v++ {
		if tin[v] == 0 {
			return nil, fmt.Errorf("graph is disconnected")
		}
	}
	if len(stack) != 0 {
		if err := validateComponent(stack); err != nil {
			return nil, err
		}
	}
	return bridges, nil
}

func bridgesFormMatching(n int, bridges []bool, edges [][2]int) bool {
	deg := make([]int, n+1)
	for i, b := range bridges {
		if !b {
			continue
		}
		u, v := edges[i][0], edges[i][1]
		deg[u]++
		deg[v]++
		if deg[u] > 1 || deg[v] > 1 {
			return false
		}
	}
	return true
}

func formatInput(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&b, "2 %d %d\n", e[0], e[1])
	}
	b.WriteString("0 0\n")
	return b.String()
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)

	rng := rand.New(rand.NewSource(151001510))
	for i := 0; i < 20; i++ {
		n := rng.Intn(18) + 2
		cycles := rng.Intn(n/2 + 1)
		tests = append(tests, randomCactus(rng, n, cycles))
	}

	// Some larger stress cases.
	tests = append(tests, randomCactus(rng, 80, 30))
	tests = append(tests, randomCactus(rng, 120, 60))
	return tests
}

func manualTests() []testCase {
	return []testCase{
		{n: 2, edges: [][2]int{{1, 2}}},
		{n: 3, edges: [][2]int{{1, 2}, {2, 3}}},
		{n: 5, edges: [][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}, {4, 5}}},
	}
}

func randomCactus(rng *rand.Rand, n, cycles int) testCase {
	if n < 1 {
		n = 1
	}
	edges := make([][2]int, 0, n-1+cycles)
	edgeSet := make(map[string]struct{})

	treeAdj := make([][]struct {
		to int
		id int
	}, n+1)

	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		id := len(edges)
		edges = append(edges, [2]int{v, p})
		treeAdj[v] = append(treeAdj[v], struct {
			to int
			id int
		}{to: p, id: id})
		treeAdj[p] = append(treeAdj[p], struct {
			to int
			id int
		}{to: v, id: id})
		edgeSet[pairKey(v, p)] = struct{}{}
	}

	parent := make([]int, n+1)
	parentEdge := make([]int, n+1)
	depth := make([]int, n+1)
	stack := []int{1}
	parent[1] = -1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range treeAdj[v] {
			if e.to == parent[v] {
				continue
			}
			parent[e.to] = v
			parentEdge[e.to] = e.id
			depth[e.to] = depth[v] + 1
			stack = append(stack, e.to)
		}
	}

	usedInCycle := make([]bool, len(edges))

	for c := 0; c < cycles; c++ {
		for attempt := 0; attempt < 30; attempt++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if _, ok := edgeSet[pairKey(u, v)]; ok {
				continue
			}
			path := treePath(u, v, parent, parentEdge, depth)
			if len(path) < 2 {
				continue
			}
			ok := true
			for _, id := range path {
				if usedInCycle[id] {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			for _, id := range path {
				usedInCycle[id] = true
			}
			edges = append(edges, [2]int{u, v})
			edgeSet[pairKey(u, v)] = struct{}{}
			break
		}
	}

	return testCase{n: n, edges: edges}
}

func treePath(u, v int, parent, parentEdge, depth []int) []int {
	var pathU, pathV []int
	uu, vv := u, v
	for depth[uu] > depth[vv] {
		pathU = append(pathU, parentEdge[uu])
		uu = parent[uu]
	}
	for depth[vv] > depth[uu] {
		pathV = append(pathV, parentEdge[vv])
		vv = parent[vv]
	}
	for uu != vv {
		pathU = append(pathU, parentEdge[uu])
		pathV = append(pathV, parentEdge[vv])
		uu = parent[uu]
		vv = parent[vv]
	}
	for i, j := 0, len(pathV)-1; i < j; i, j = i+1, j-1 {
		pathV[i], pathV[j] = pathV[j], pathV[i]
	}
	return append(pathU, pathV...)
}

func pairKey(u, v int) string {
	if u > v {
		u, v = v, u
	}
	return fmt.Sprintf("%d_%d", u, v)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
