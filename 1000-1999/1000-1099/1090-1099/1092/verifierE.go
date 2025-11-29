package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `10 4 6 9 1 8 4 1 3 2
6 3 2 4 5 1 5 2
1 0
7 2 2 7 7 4
3 0
3 2 3 2 1 3
4 1 2 4
5 0
6 3 3 1 1 5 3 4
6 1 6 2
1 0
1 0
7 0
9 6 6 7 1 8 1 3 4 2 4 8 6 9
6 4 3 4 1 5 6 3 3 1
7 0
4 2 3 2 1 3
6 2 2 1 6 2
5 3 2 1 1 5 5 4
1 0
10 5 5 8 7 3 8 6 4 3 10 3
7 0
3 1 2 1
1 0
5 1 4 5
3 2 2 3 1 3
8 5 3 2 8 5 6 2 6 1 5 6
9 4 8 5 5 6 3 1 8 9
5 2 3 4 3 5
6 2 3 6 4 3
3 2 2 3 1 3
3 0
6 3 3 6 1 6 4 2
10 9 9 7 5 10 9 5 1 4 3 10 8 10 3 4 3 1 8 4
3 0
3 0
6 1 4 2
9 0
7 3 3 4 6 5 1 5
4 1 3 1
6 3 3 4 1 6 5 3
1 0
2 1 2 1
7 3 7 5 6 5 3 4
3 0
10 6 10 8 4 3 10 2 6 1 7 2 6 10
10 8 3 6 10 7 4 8 5 8 3 10 10 5 5 7 1 6
5 3 3 2 4 1 1 5
8 3 5 1 3 7 1 8
9 8 5 4 8 1 4 8 5 3 9 2 1 3 9 6 5 9
1 0
6 2 5 2 3 5
8 1 4 1
7 6 7 4 6 7 6 2 3 4 1 3 7 1
7 4 3 6 7 4 3 2 2 4
9 7 6 9 3 7 9 1 2 4 5 2 2 1 6 7
2 1 2 1
2 0
4 0
3 1 3 2
3 0
10 2 7 3 2 10
4 0
9 1 7 2
5 0
10 9 2 7 10 3 1 7 2 6 10 8 8 6 6 1 3 5 9 5
9 8 4 5 2 9 5 9 3 4 6 8 7 3 3 1 6 2
10 0
2 0
9 7 4 7 6 2 9 1 9 7 1 3 8 4 3 7
6 1 5 6
6 3 4 2 6 3 6 4
6 4 1 5 5 2 6 2 4 1
2 0
1 0
4 1 1 4
8 5 1 7 8 5 7 6 2 4 2 6
7 6 4 2 5 3 1 3 7 2 7 3 6 2
1 0
1 0
6 2 1 6 2 1
6 2 4 6 2 6
3 0
4 0
10 3 2 10 7 6 5 3
8 5 5 2 8 3 4 3 4 1 7 5
1 0
1 0
7 1 5 7
9 5 4 8 2 8 1 4 1 8 7 1
1 0
7 3 3 6 6 5 6 4
4 1 1 3
2 0
9 2 2 4 1 8
9 5 8 7 7 9 3 1 5 8 4 9
2 0
5 3 1 3 4 3 5 1
8 3 2 7 4 5 1 8
4 0`

type edge struct{ u, v int }

type testCase struct {
	n     int
	edges []edge
}

type dsu struct{ p []int }

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(x, y int) {
	rx, ry := d.find(x), d.find(y)
	if rx != ry {
		d.p[rx] = ry
	}
}

type bfsRes struct {
	node   int
	dist   int
	parent []int
}

func bfs(start, n int, adj [][]int) bfsRes {
	dist := make([]int, n+1)
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	q = append(q, start)
	dist[start] = 0
	parent[start] = -1
	resNode := start
	resDist := 0
	for idx := 0; idx < len(q); idx++ {
		u := q[idx]
		if dist[u] > resDist {
			resDist = dist[u]
			resNode = u
		}
		for _, v := range adj[u] {
			if dist[v] < 0 {
				dist[v] = dist[u] + 1
				parent[v] = u
				q = append(q, v)
			}
		}
	}
	return bfsRes{node: resNode, dist: resDist, parent: parent}
}

// referenceSolve embeds logic from 1092E.go.
func referenceSolve(tc testCase) (string, int, []edge) {
	n := tc.n
	adj := make([][]int, n+1)
	d := newDSU(n)
	for _, e := range tc.edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
		d.union(e.u, e.v)
	}
	compMap := make(map[int][]int)
	for i := 1; i <= n; i++ {
		r := d.find(i)
		compMap[r] = append(compMap[r], i)
	}
	type comp struct{ center, radius int }
	comps := make([]comp, 0, len(compMap))
	for _, nodes := range compMap {
		u0 := nodes[0]
		r1 := bfs(u0, n, adj)
		r2 := bfs(r1.node, n, adj)
		path := []int{r2.node}
		for cur := r2.node; cur != r1.node; cur = r2.parent[cur] {
			path = append(path, r2.parent[cur])
		}
		L := len(path)
		center := path[L/2]
		d1 := L / 2
		d2 := L - 1 - L/2
		radius := d1
		if d2 > radius {
			radius = d2
		}
		comps = append(comps, comp{center: center, radius: radius})
	}
	sort.Slice(comps, func(i, j int) bool {
		return comps[i].radius > comps[j].radius
	})
	var added []edge
	if len(comps) > 0 {
		mainCenter := comps[0].center
		for i := 1; i < len(comps); i++ {
			c := comps[i].center
			added = append(added, edge{mainCenter, c})
			adj[mainCenter] = append(adj[mainCenter], c)
			adj[c] = append(adj[c], mainCenter)
		}
	}
	// compute diameter after adding edges
	start := 1
	r1 := bfs(start, n, adj)
	r2 := bfs(r1.node, n, adj)
	diameter := r2.dist
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(diameter))
	sb.WriteByte('\n')
	for _, e := range added {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return strings.TrimSpace(sb.String()), diameter, added
}

func runSolution(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func countComponents(tc testCase) int {
	d := newDSU(tc.n)
	for _, e := range tc.edges {
		d.union(e.u, e.v)
	}
	seen := make(map[int]struct{})
	for i := 1; i <= tc.n; i++ {
		seen[d.find(i)] = struct{}{}
	}
	return len(seen)
}

func parseOutput(tc testCase, output string) (int, []edge, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	diam, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid diameter")
	}
	expectedEdges := countComponents(tc) - 1
	if len(fields) != 1+expectedEdges*2 {
		return 0, nil, fmt.Errorf("expected %d edges, got %d tokens", expectedEdges, (len(fields)-1)/2)
	}
	edges := make([]edge, expectedEdges)
	for i := 0; i < expectedEdges; i++ {
		u, err := strconv.Atoi(fields[1+2*i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid edge u")
		}
		v, err := strconv.Atoi(fields[2+2*i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid edge v")
		}
		if u < 1 || u > tc.n || v < 1 || v > tc.n {
			return 0, nil, fmt.Errorf("edge vertex out of range")
		}
		edges[i] = edge{u, v}
	}
	return diam, edges, nil
}

func diameterWithEdges(tc testCase, extra []edge) (int, bool) {
	n := tc.n
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	for _, e := range extra {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	// check connectivity
	visited := make([]bool, n+1)
	stack := []int{1}
	visited[1] = true
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				stack = append(stack, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return 0, false
		}
	}
	r1 := bfs(1, n, adj)
	r2 := bfs(r1.node, n, adj)
	return r2.dist, true
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		vals := strings.Fields(line)
		if len(vals) < 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		n, err := strconv.Atoi(vals[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(vals[1])
		if err != nil {
			return nil, err
		}
		if len(vals) != 2+2*m {
			return nil, fmt.Errorf("invalid edge count for n=%d m=%d line=%q", n, m, line)
		}
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			u, err := strconv.Atoi(vals[2+2*i])
			if err != nil {
				return nil, err
			}
			v, err := strconv.Atoi(vals[3+2*i])
			if err != nil {
				return nil, err
			}
			edges[i] = edge{u: u, v: v}
		}
		tests = append(tests, testCase{n: n, edges: edges})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		_, expectedDiam, _ := referenceSolve(tc)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		got, err := runSolution(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		diam, edges, err := parseOutput(tc, got)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if diam != expectedDiam {
			fmt.Printf("case %d failed: expected diameter %d got %d\n", i+1, expectedDiam, diam)
			os.Exit(1)
		}
		if d, ok := diameterWithEdges(tc, edges); !ok || d != diam {
			fmt.Printf("case %d failed: reported diameter mismatch after validation (got %d)\n", i+1, d)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
