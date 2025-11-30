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

const testcases = `6 9 3 4 3 1 4 4 1 3 8 1 5 4 1 6 3 3 6 9 2 4 9 1 6 4 2 4 2
3 3 1 2 6 1 2 7 1 3 4
4 3 3 4 3 1 2 5 1 3 5
2 1 1 2 9
8 13 3 8 1 4 6 9 2 5 5 1 8 1 5 6 8 3 6 6 2 3 5 3 4 2 5 6 3 2 6 5 4 8 6 6 7 10 1 7 2
2 1 1 2 6
8 10 6 7 3 1 7 5 3 8 10 3 8 6 2 3 8 5 8 7 2 3 7 3 8 8 6 8 5 2 7 3
4 5 3 4 3 3 4 5 1 3 10 1 3 6 1 2 8
8 14 1 7 8 3 4 3 1 5 4 3 8 10 4 6 8 2 6 3 1 8 4 1 6 6 2 5 7 2 4 4 3 7 5 1 3 3 6 7 10 4 8 9
5 6 3 5 9 1 4 9 3 5 2 3 5 5 2 3 7 1 3 9
5 10 3 4 3 2 5 2 4 5 4 2 4 8 2 5 7 2 4 10 1 3 7 1 3 10 2 3 10 3 4 7
6 9 3 4 3 5 6 8 1 4 2 2 4 1 2 3 8 4 5 9 3 4 1 3 6 8 2 3 1
8 7 5 6 9 4 8 5 1 2 3 1 5 9 3 5 5 1 8 6 3 6 1
5 9 3 4 3 1 5 4 1 4 7 3 4 9 1 3 1 2 4 10 2 3 8 2 4 9 4 5 5
6 5 1 3 2 4 6 2 1 6 6 4 6 8 1 6 4
2 1 1 2 3
2 1 1 2 5
2 1 1 2 1
5 4 3 5 8 2 5 9 3 4 6 1 2 5
4 4 1 3 2 3 4 7 1 3 9 2 3 5
6 15 1 2 7 3 6 7 3 5 6 1 5 7 2 6 6 2 6 3 2 4 7 1 5 10 2 3 5 5 6 7 2 6 2 1 6 8 1 4 2 3 5 1 1 3 9
2 1 1 2 7
5 6 2 5 5 1 2 6 1 2 3 1 3 3 4 5 6 2 3 1
5 4 1 2 6 1 3 1 3 4 4 1 2 4
7 6 2 5 2 3 6 3 2 7 4 3 4 5 1 6 8 4 5 6
3 2 2 3 5 1 3 10
2 1 1 2 7
3 2 2 3 2 1 3 9
7 7 2 6 1 3 7 7 1 7 1 1 2 3 2 7 8 3 4 1 1 4 1
7 14 2 5 2 2 3 3 1 5 4 3 7 7 4 5 7 3 4 2 1 6 5 1 5 3 3 5 8 4 6 4 2 6 2 3 7 2 5 6 3 1 5 8
8 13 1 2 1 5 7 2 4 5 1 4 8 1 1 2 3 3 7 6 5 6 1 1 5 3 1 8 3 3 5 8 6 7 7 4 5 6 2 5 3
8 12 4 6 6 2 6 6 4 8 8 1 2 3 3 7 6 2 3 5 7 8 7 1 4 9 1 3 10 5 7 10 3 7 8 1 3 9
8 13 4 6 6 2 7 5 1 2 3 4 6 8 6 7 8 4 8 7 4 6 1 2 8 4 1 7 2 7 8 6 4 5 5 1 4 1 1 4 8
6 14 4 5 8 1 3 5 1 6 9 3 6 6 3 5 9 2 4 1 2 3 5 5 6 7 2 6 2 4 5 9 1 5 9 1 3 10 2 4 6 3 6 1
7 15 6 7 9 3 7 10 2 3 9 2 6 10 3 4 6 2 4 8 2 3 2 2 3 5 4 6 7 1 2 2 2 6 2 2 4 9 1 7 9 4 6 3 1 6 7
6 14 2 6 10 2 3 6 2 6 3 3 5 9 2 4 7 2 4 4 1 4 3 1 2 2 3 5 2 1 5 6 3 6 2 1 3 10 2 3 7 1 4 1
6 6 1 3 2 5 6 2 2 6 9 4 5 6 2 4 9 2 4 2
5 7 3 4 3 1 3 2 1 3 8 3 5 3 1 4 7 3 4 2 3 4 4
6 6 1 2 7 3 4 6 4 5 7 1 5 7 2 3 2 2 5 7
8 10 2 8 6 4 8 1 1 8 1 2 8 5 3 7 9 1 5 9 2 3 1 7 8 9 3 6 1 1 2 4
5 10 2 3 3 1 5 4 1 4 7 1 2 6 1 4 6 2 4 7 3 4 5 1 4 1 2 3 4 2 4 2
5 5 2 4 4 3 5 8 1 3 10 1 2 8 2 4 2
4 4 3 4 9 3 4 5 2 4 10 1 3 5
3 3 1 2 1 1 3 1 2 3 3
2 1 1 2 4
6 14 4 6 6 3 6 10 2 4 8 4 5 7 3 5 6 3 5 9 1 6 2 2 5 1 1 4 9 1 6 8 2 3 1 1 6 4 2 3 10 1 5 5
4 3 1 4 10 2 4 1 1 3 3
4 5 3 4 3 1 2 6 1 4 10 1 3 7 2 4 2
5 9 2 5 2 1 3 6 3 4 9 4 5 10 3 4 5 4 5 2 1 5 2 1 4 5 1 5 5
6 15 1 2 7 2 6 7 1 2 6 2 5 8 3 5 6 2 4 4 3 5 2 4 5 9 3 5 8 2 5 4 4 5 6 2 4 9 2 3 7 4 6 3 1 6 7
4 5 1 3 2 1 2 9 1 3 10 1 4 2 2 4 6
7 10 6 7 9 1 4 4 1 3 2 1 4 10 3 6 6 4 7 7 1 4 6 3 5 2 1 3 3 1 4 8
7 14 1 7 8 2 3 9 5 7 9 1 4 7 1 3 5 4 5 7 3 4 9 1 2 9 4 7 10 1 5 9 1 2 5 2 4 3 2 6 5 1 2 4
2 1 1 2 2
5 5 3 5 8 3 5 5 1 2 8 1 3 9 2 4 2
7 14 1 2 1 3 4 3 3 6 4 3 6 7 3 7 7 3 5 6 5 6 8 4 6 8 3 6 9 1 4 3 1 3 3 6 7 1 4 5 5 3 5 10
6 10 4 5 8 4 5 1 1 3 4 2 5 7 1 5 3 4 6 4 1 6 4 1 5 8 3 6 1 1 6 7
3 2 2 3 4 1 3 9
7 7 1 4 4 2 7 2 2 4 1 5 7 8 6 7 1 3 5 1 3 5 10
3 3 1 3 1 2 3 9 2 3 10
3 2 1 3 6 2 3 4
2 1 1 2 1
2 1 1 2 4
7 12 6 7 9 3 4 9 2 7 5 1 7 4 2 4 1 4 7 7 1 2 9 6 7 8 4 6 4 3 5 5 2 6 4 2 4 2
8 9 6 7 3 2 8 3 2 6 10 6 8 5 6 8 10 3 7 2 1 8 6 3 5 1 1 6 1
6 11 1 6 10 2 4 5 1 3 2 4 5 4 3 6 7 2 4 8 4 5 1 2 3 6 2 4 7 1 4 9 1 5 2
4 4 1 4 9 3 4 10 1 3 7 2 3 2
4 6 1 2 6 2 3 2 1 4 10 1 4 6 3 4 2 3 4 1
3 3 2 3 7 2 3 9 1 3 10
2 1 1 2 6
2 1 1 2 3
7 10 5 6 9 2 3 6 4 7 1 4 7 10 5 7 1 5 7 7 3 6 5 1 5 2 3 4 4 1 4 8
3 3 1 3 7 1 3 2 1 2 3
7 14 1 4 4 5 6 9 3 5 3 1 4 7 1 7 1 5 7 5 2 4 1 2 6 8 3 4 5 4 5 9 2 7 7 1 3 6 2 5 3 1 4 5
4 6 3 4 3 1 3 4 1 4 3 2 3 8 1 2 2 1 2 4
3 3 1 2 8 1 2 10 1 3 5
2 1 1 2 1
8 7 5 8 3 1 4 4 6 7 6 4 7 10 4 5 3 5 6 4 1 5 2
3 2 1 2 8 1 3 10
7 10 1 7 5 2 5 5 1 6 6 3 4 8 4 5 9 2 5 10 3 6 2 1 5 3 5 7 10 4 7 9
7 11 2 4 5 1 5 4 2 5 5 3 6 3 1 6 6 2 5 1 2 6 2 4 5 9 5 7 1 3 4 5 3 6 2
8 13 5 6 9 3 7 4 3 6 7 1 5 4 7 8 8 5 8 2 1 3 1 1 4 9 1 7 10 3 7 8 3 8 2 3 8 8 4 6 3
8 7 6 7 9 1 6 10 1 4 4 2 7 8 1 6 8 2 3 1 4 6 3
6 9 4 5 8 2 3 6 4 6 2 1 4 6 1 3 9 2 3 4 2 3 7 2 6 4 1 2 4
5 6 3 4 3 1 5 9 1 3 3 1 4 2 2 3 1 2 5 6
4 6 1 3 5 3 4 8 1 3 7 2 4 3 1 3 9 1 4 8
8 8 4 5 4 4 7 8 4 6 8 1 8 6 2 7 7 2 4 9 4 6 3 4 7 2
8 14 3 6 4 3 7 4 3 4 6 1 4 7 5 8 2 5 6 5 1 6 6 3 8 3 6 8 6 5 7 1 2 5 10 1 8 3 6 8 2 2 4 9
5 8 1 2 7 3 4 10 1 4 7 1 3 1 1 4 10 2 3 4 2 5 3 3 4 7
2 1 1 2 6
3 3 1 2 6 1 2 4 2 3 6
4 3 3 4 3 2 4 7 2 3 10
5 10 3 5 7 4 5 8 1 5 4 4 5 5 2 5 6 1 5 10 1 4 3 2 4 3 2 3 10 3 5 10
4 5 2 4 5 1 4 10 1 4 6 2 4 7 1 3 7
4 5 1 2 7 1 2 6 1 3 9 1 2 8 2 3 10
2 1 1 2 4
7 13 3 6 10 2 6 7 3 6 3 4 5 10 4 7 1 1 6 2 3 4 2 5 6 7 1 7 7 1 6 5 5 6 4 2 7 10 1 6 7
3 3 1 2 6 2 3 7 2 3 8
2 1 1 2 6
4 6 1 3 5 2 4 7 1 2 2 2 3 4 1 4 2 3 4 4`

// Embedded solution logic from 1927F.go.
type DSU struct {
	parent, size []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	sz := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{parent: p, size: sz}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	ra, rb := d.Find(a), d.Find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

type Edge struct {
	u, v, w int
}

func solveCase(n int, edges []Edge) (int, []int) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w > edges[j].w
	})

	dsu := NewDSU(n)
	g := make([][]int, n+1)
	var st, en, cost int
	for _, e := range edges {
		u, v, w := e.u, e.v, e.w
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
		if dsu.Find(u) == dsu.Find(v) {
			st, en, cost = u, v, w
		}
		dsu.Union(u, v)
	}

	vis := make([]bool, n+1)
	parent := make([]int, n+1)
	q := []int{st}
	vis[st] = true
	for head := 0; head < len(q); head++ {
		u := q[head]
		if u == en {
			break
		}
		for _, v := range g[u] {
			if u == st && v == en {
				continue
			}
			if !vis[v] {
				vis[v] = true
				parent[v] = u
				q = append(q, v)
			}
		}
	}

	path := make([]int, 0)
	cur := en
	for {
		path = append(path, cur)
		if cur == st {
			break
		}
		cur = parent[cur]
	}
	return cost, path
}

type testCase struct {
	n, m  int
	edges []Edge
}

func parseCases(raw string) ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		vals := strings.Fields(line)
		if len(vals) < 2 {
			return nil, fmt.Errorf("line %d too short", idx+1)
		}
		n, err := strconv.Atoi(vals[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(vals[1])
		if err != nil {
			return nil, err
		}
		need := 2 + 3*m
		if len(vals) < need {
			return nil, fmt.Errorf("line %d length mismatch", idx+1)
		}
		es := make([]Edge, m)
		pos := 2
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(vals[pos])
			v, _ := strconv.Atoi(vals[pos+1])
			w, _ := strconv.Atoi(vals[pos+2])
			pos += 3
			es[i] = Edge{u: u, v: v, w: w}
		}
		res = append(res, testCase{n: n, m: m, edges: es})
	}
	return res, nil
}

func runProg(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases(testcases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		cost, path := solveCase(tc.n, tc.edges)
		var exp strings.Builder
		exp.WriteString(fmt.Sprintf("%d %d\n", cost, len(path)))
		for _, v := range path {
			exp.WriteString(strconv.Itoa(v))
			exp.WriteByte(' ')
		}
		want := strings.TrimSpace(exp.String())

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
		}

		got, err := runProg(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
