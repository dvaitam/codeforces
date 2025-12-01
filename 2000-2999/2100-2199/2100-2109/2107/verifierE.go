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
)

const refSource = "./2107E.go"

type testCase struct {
	n int
	k int64
}

type caseAns struct {
	ok    bool
	edges [][2]int
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[len(args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := validateCase(tc, refAns[i].ok, candAns[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (n=%d k=%d): %v\n", i+1, tc.n, tc.k, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2107E-ref-*")
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

func buildTests() []testCase {
	var tcs []testCase
	add := func(n int, k int64) { tcs = append(tcs, testCase{n: n, k: k}) }

	add(2, 0)
	add(2, 1)
	add(3, 0)
	add(3, 1)
	add(4, 3)
	add(5, 7)
	add(6, 8)
	add(8, 20)
	add(10, 45)
	add(12, 120)
	add(20, 0)

	rng := rand.New(rand.NewSource(2107))
	for len(tcs) < 30 {
		n := rng.Intn(80) + 2
		if len(tcs)%7 == 0 {
			n = rng.Intn(3000) + 500
		}
		maxW := choose3(int64(n))
		k := rng.Int63n(maxW + 5)
		if rng.Intn(5) == 0 {
			k = 0
		}
		add(n, k)
	}

	// Large stress tests but respecting total n <= 2e5.
	add(100000, choose3(10)) // small target with large n
	add(50000, 1e12)

	return tcs
}

func buildInput(tcs []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	}
	return sb.String()
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string, tests []testCase) ([]caseAns, error) {
	tokens := strings.Fields(out)
	pos := 0
	res := make([]caseAns, len(tests))
	for i, tc := range tests {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("output ended early at case %d", i+1)
		}
		tok := strings.ToLower(tokens[pos])
		pos++
		if tok[0] == 'n' {
			res[i] = caseAns{ok: false}
			continue
		}
		if tok[0] != 'y' {
			return nil, fmt.Errorf("case %d: expected Yes/No, got %q", i+1, tokens[pos-1])
		}
		need := tc.n - 1
		if pos+2*need > len(tokens) {
			return nil, fmt.Errorf("case %d: expected %d edges, got fewer tokens", i+1, need)
		}
		edges := make([][2]int, need)
		for e := 0; e < need; e++ {
			u, err1 := strconv.Atoi(tokens[pos+2*e])
			v, err2 := strconv.Atoi(tokens[pos+2*e+1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("case %d: invalid edge token around position %d", i+1, pos+2*e)
			}
			edges[e] = [2]int{u, v}
		}
		pos += 2 * need
		res[i] = caseAns{ok: true, edges: edges}
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra output tokens detected (%d unused)", len(tokens)-pos)
	}
	return res, nil
}

func validateCase(tc testCase, expectedPossible bool, ans caseAns) error {
	if !expectedPossible {
		if ans.ok {
			return fmt.Errorf("expected No (reference unsatisfiable), but got a tree")
		}
		return nil
	}
	if !ans.ok {
		return fmt.Errorf("expected Yes (solution exists), but got No")
	}
	if len(ans.edges) != tc.n-1 {
		return fmt.Errorf("expected %d edges, got %d", tc.n-1, len(ans.edges))
	}
	if err := checkTree(tc.n, ans.edges, tc.k); err != nil {
		return err
	}
	return nil
}

func checkTree(n int, edges [][2]int, k int64) error {
	uf := newDSU(n)
	for i, e := range edges {
		u, v := e[0], e[1]
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge %d has invalid node indices (%d,%d)", i+1, u, v)
		}
		if !uf.union(u-1, v-1) {
			return fmt.Errorf("edge %d introduces a cycle between %d and %d", i+1, u, v)
		}
	}
	if uf.count != 1 {
		return fmt.Errorf("graph is not connected")
	}

	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	depth := make([]int, n)
	order := make([]int, 0, n)
	stack := []int{0}
	parent[0] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if parent[to] != -1 {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	if len(order) != n {
		return fmt.Errorf("tree not connected to root 1")
	}

	size := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		size[v] = 1
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			size[v] += size[to]
		}
	}

	var weight int64
	for v := 0; v < n; v++ {
		s := size[v]
		pairs := s * (s - 1) / 2
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			sChild := size[to]
			pairs -= sChild * (sChild - 1) / 2
		}
		weight += int64(depth[v]) * pairs
	}

	diff := weight - k
	if diff < 0 {
		diff = -diff
	}
	if diff > 1 {
		return fmt.Errorf("tree weight %d differs from k=%d by %d", weight, k, diff)
	}
	return nil
}

func choose3(x int64) int64 {
	return x * (x - 1) * (x - 2) / 6
}

type dsu struct {
	parent []int
	rank   []byte
	count  int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &dsu{parent: p, rank: make([]byte, n), count: n}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(x, y int) bool {
	fx, fy := d.find(x), d.find(y)
	if fx == fy {
		return false
	}
	if d.rank[fx] < d.rank[fy] {
		fx, fy = fy, fx
	}
	d.parent[fy] = fx
	if d.rank[fx] == d.rank[fy] {
		d.rank[fx]++
	}
	d.count--
	return true
}
