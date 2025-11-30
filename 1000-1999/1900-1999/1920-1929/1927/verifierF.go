package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

//go:embed testcasesF.txt
var testcases string

type dsu struct {
	parent, size []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	sz := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &dsu{parent: p, size: sz}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

type edge struct {
	u, v, w int
}

func solveCase(n int, edges []edge) (int, []int) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w > edges[j].w
	})
	d := newDSU(n)
	g := make([][]int, n+1)
	st, en, cost := 0, 0, 0
	for _, e := range edges {
		u, v, w := e.u, e.v, e.w
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
		if d.find(u) == d.find(v) {
			st, en, cost = u, v, w
		}
		d.union(u, v)
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
	edges []edge
}

func parseCases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
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
		es := make([]edge, m)
		pos := 2
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(vals[pos])
			v, _ := strconv.Atoi(vals[pos+1])
			w, _ := strconv.Atoi(vals[pos+2])
			pos += 3
			es[i] = edge{u: u, v: v, w: w}
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
