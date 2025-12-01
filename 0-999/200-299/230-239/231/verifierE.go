package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD = 1000000007

const embeddedTestcasesE = `100
5 4
2 1
3 2
4 2
5 1
2
5 2
2 5
4 3
2 1
3 2
4 3
1
4 2
5 4
2 1
3 2
4 3
5 1
2
5 3
4 1
4 3
2 1
3 1
4 3
1
3 1
3 3
2 1
3 1
3 2
1
3 1
5 4
2 1
3 2
4 3
5 3
2
2 4
5 4
5 4
2 1
3 2
4 1
5 3
2
4 3
3 5
3 2
2 1
3 2
1
3 1
2 1
2 1
2
2 1
2 1
4 3
2 1
3 1
4 3
2
3 1
3 1
5 4
2 1
3 2
4 3
5 3
2
1 4
2 3
3 3
2 1
3 1
2 3
1
2 1
2 1
2 1
1
1 2
2 1
2 1
2
1 2
2 1
4 3
2 1
3 2
4 2
2
1 4
3 1
3 3
2 1
3 1
3 2
1
3 1
5 5
2 1
3 2
4 1
5 4
1 3
2
3 2
4 1
2 1
2 1
2
1 2
2 1
4 3
2 1
3 2
4 3
1
1 3
4 3
2 1
3 1
4 3
2
3 1
4 2
4 3
2 1
3 1
4 1
1
2 1
4 4
2 1
3 1
4 2
3 4
1
4 3
2 1
2 1
1
2 1
4 4
2 1
3 2
4 3
2 4
1
2 4
5 5
2 1
3 1
4 2
5 2
4 5
1
1 4
4 3
2 1
3 2
4 1
1
4 1
4 3
2 1
3 1
4 2
2
2 3
1 4
4 3
2 1
3 2
4 3
2
1 4
2 3
3 2
2 1
3 1
2
3 1
1 3
5 4
2 1
3 2
4 3
5 2
2
5 4
3 4
3 2
2 1
3 2
2
2 1
2 3
3 2
2 1
3 1
2
2 1
3 2
3 2
2 1
3 2
1
1 2
3 3
2 1
3 2
3 1
1
1 2
4 3
2 1
3 1
4 1
2
3 1
1 3
5 5
2 1
3 2
4 3
5 2
4 1
2
5 3
2 1
2 1
2 1
2
1 2
2 1
2 1
2 1
2
1 2
1 2
2 1
2 1
1
1 2
3 2
2 1
3 1
2
2 1
1 3
4 4
2 1
3 1
4 3
3 2
2
2 1
1 3
5 4
2 1
3 1
4 1
5 3
1
2 3
2 1
2 1
1
2 1
3 2
2 1
3 1
2
2 3
3 1
5 4
2 1
3 1
4 3
5 1
2
4 2
4 3
5 4
2 1
3 2
4 3
5 2
2
3 4
3 4
3 3
2 1
3 2
1 3
2
3 2
1 2
5 4
2 1
3 2
4 1
5 3
1
3 1
5 5
2 1
3 1
4 3
5 1
5 3
2
3 1
5 4
2 1
2 1
1
2 1
2 1
2 1
1
1 2
4 4
2 1
3 1
4 1
2 4
2
4 2
1 4
2 1
2 1
2
1 2
2 1
2 1
2 1
2
1 2
2 1
2 1
2 1
1
2 1
5 4
2 1
3 2
4 2
5 3
1
1 3
4 4
2 1
3 1
4 1
2 3
1
4 2
4 3
2 1
3 2
4 1
2
4 1
2 1
4 3
2 1
3 1
4 2
1
4 3
5 4
2 1
3 1
4 2
5 4
1
1 2
3 2
2 1
3 2
2
1 3
1 2
3 2
2 1
3 1
2
2 3
3 1
5 4
2 1
3 1
4 1
5 2
1
5 4
4 3
2 1
3 2
4 1
2
2 3
4 1
5 4
2 1
3 2
4 1
5 2
1
5 2
5 5
2 1
3 2
4 1
5 1
3 1
1
1 4
2 1
2 1
2
1 2
1 2
2 1
2 1
2
2 1
1 2
2 1
2 1
1
2 1
4 3
2 1
3 2
4 2
1
1 4
2 1
2 1
1
1 2
3 2
2 1
3 1
2
2 1
1 3
5 4
2 1
3 1
4 3
5 1
2
5 4
2 5
2 1
2 1
1
2 1
4 3
2 1
3 2
4 1
2
3 1
1 4
5 5
2 1
3 1
4 3
5 1
2 3
2
5 3
5 2
2 1
2 1
1
2 1
3 2
2 1
3 2
1
2 1
3 2
2 1
3 1
2
2 3
1 3
4 4
2 1
3 1
4 3
4 1
2
1 4
3 2
2 1
2 1
2
1 2
2 1
2 1
2 1
2
1 2
1 2
5 4
2 1
3 1
4 1
5 1
1
4 3
3 3
2 1
3 2
1 3
2
3 2
1 2
4 4
2 1
3 2
4 3
4 1
2
3 4
3 4
2 1
2 1
2
2 1
2 1
4 3
2 1
3 1
4 3
1
2 1
5 4
2 1
3 2
4 1
5 1
2
3 4
1 2
5 4
2 1
3 1
4 2
5 4
1
3 1
3 2
2 1
3 1
2
1 3
1 3
2 1
2 1
2
2 1
1 2
4 3
2 1
3 2
4 2
1
2 1
4 3
2 1
3 1
4 3
1
3 2
3 2
2 1
3 2
2
3 1
1 2
5 4
2 1
3 1
4 1
5 4
1
5 3
4 3
2 1
3 1
4 1
1
2 3
2 1
2 1
2
1 2
1 2
4 3
2 1
3 1
4 3
1
4 3
4 3
2 1
3 1
4 1
2
3 1
2 1
4 4
2 1
3 1
4 3
4 2
2
1 4
2 3
`

type edge struct{ to, id int }
type frame struct{ u, peid, dep, idx int }
type testCase struct {
	n, m    int
	edges   [][2]int // zero-based
	queries [][2]int // zero-based
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(embeddedTestcasesE))
	scan.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scan.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		m, err := nextInt()
		if err != nil {
			return nil, err
		}
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			u, err := nextInt()
			if err != nil {
				return nil, err
			}
			v, err := nextInt()
			if err != nil {
				return nil, err
			}
			edges[j] = [2]int{u - 1, v - 1}
		}
		k, err := nextInt()
		if err != nil {
			return nil, err
		}
		queries := make([][2]int, k)
		for j := 0; j < k; j++ {
			x, err := nextInt()
			if err != nil {
				return nil, err
			}
			y, err := nextInt()
			if err != nil {
				return nil, err
			}
			queries[j] = [2]int{x - 1, y - 1}
		}
		cases = append(cases, testCase{n: n, m: m, edges: edges, queries: queries})
	}
	return cases, nil
}

func lca(a, b int, depth []int, up [][]int, LOG int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := 0; diff > 0; k++ {
		if diff&1 != 0 {
			a = up[k][a]
		}
		diff >>= 1
	}
	if a == b {
		return a
	}
	for k := LOG - 1; k >= 0; k-- {
		ua := up[k][a]
		ub := up[k][b]
		if ua != ub {
			a = ua
			b = ub
		}
	}
	return up[0][a]
}

func solveCase(n, m int, edges [][2]int, queries [][2]int) []int {
	g := make([][]edge, n)
	uE := make([]int, m)
	vE := make([]int, m)
	for i, e := range edges {
		u, v := e[0], e[1]
		uE[i], vE[i] = u, v
		g[u] = append(g[u], edge{to: v, id: i})
		g[v] = append(g[v], edge{to: u, id: i})
	}

	visited := make([]int8, n)
	parent := make([]int, n)
	parentE := make([]int, n)
	depth := make([]int, n)
	edgeInCycle := make([]bool, m)
	var cycles [][]int
	var stk []frame
	stk = append(stk, frame{u: 0, peid: -1, dep: 0, idx: 0})
	parent[0], parentE[0], depth[0] = -1, -1, 0

	for len(stk) > 0 {
		f := &stk[len(stk)-1]
		u := f.u
		if visited[u] == 0 {
			visited[u] = 1
			parentE[u] = f.peid
			depth[u] = f.dep
		}
		if f.idx < len(g[u]) {
			e := g[u][f.idx]
			f.idx++
			v := e.to
			eid := e.id
			if eid == parentE[u] {
				continue
			}
			if visited[v] == 0 {
				parent[v] = u
				stk = append(stk, frame{u: v, peid: eid, dep: f.dep + 1, idx: 0})
				continue
			}
			if visited[v] == 1 && depth[v] < f.dep {
				edgeInCycle[eid] = true
				w := u
				var cyc []int
				for w != v {
					cyc = append(cyc, w)
					edgeInCycle[parentE[w]] = true
					w = parent[w]
				}
				cyc = append(cyc, v)
				cycles = append(cycles, cyc)
			}
		} else {
			visited[u] = 2
			stk = stk[:len(stk)-1]
		}
	}

	cCnt := len(cycles)
	tot := n + cCnt
	bAdj := make([][]int, tot)
	for i := 0; i < m; i++ {
		if edgeInCycle[i] {
			continue
		}
		u, v := uE[i], vE[i]
		bAdj[u] = append(bAdj[u], v)
		bAdj[v] = append(bAdj[v], u)
	}
	isCycleNode := make([]bool, tot)
	for cid, cyc := range cycles {
		bid := n + cid
		isCycleNode[bid] = true
		for _, u := range cyc {
			bAdj[bid] = append(bAdj[bid], u)
			bAdj[u] = append(bAdj[u], bid)
		}
	}

	LOG := 1
	for (1 << LOG) <= tot {
		LOG++
	}
	up := make([][]int, LOG)
	for i := range up {
		up[i] = make([]int, tot)
	}
	depthB := make([]int, tot)
	distC := make([]int, tot)
	parentB := make([]int, tot)
	queue := make([]int, 0, tot)
	queue = append(queue, 0)
	parentB[0] = 0
	up[0][0] = 0
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		for _, v := range bAdj[u] {
			if v == parentB[u] {
				continue
			}
			parentB[v] = u
			depthB[v] = depthB[u] + 1
			distC[v] = distC[u]
			if isCycleNode[v] {
				distC[v]++
			}
			up[0][v] = u
			queue = append(queue, v)
		}
	}
	for k := 1; k < LOG; k++ {
		for v := 0; v < tot; v++ {
			up[k][v] = up[k-1][up[k-1][v]]
		}
	}

	maxC := cCnt
	pow2 := make([]int, maxC+2)
	pow2[0] = 1
	for i := 1; i <= maxC; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}

	res := make([]int, len(queries))
	for i, q := range queries {
		x, y := q[0], q[1]
		l := lca(x, y, depthB, up, LOG)
		cnt := distC[x] + distC[y] - 2*distC[l]
		if isCycleNode[l] {
			cnt++
		}
		if cnt < 0 {
			cnt = 0
		}
		if cnt <= maxC {
			res[i] = pow2[cnt]
		} else {
			v := 1
			for j := 0; j < cnt; j++ {
				v = v * 2 % MOD
			}
			res[i] = v
		}
	}
	return res
}

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		expectedVals := solveCase(tc.n, tc.m, tc.edges, tc.queries)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0]+1, e[1]+1)
		}
		fmt.Fprintf(&input, "%d\n", len(tc.queries))
		for _, q := range tc.queries {
			fmt.Fprintf(&input, "%d %d\n", q[0]+1, q[1]+1)
		}
		var want strings.Builder
		for _, v := range expectedVals {
			fmt.Fprintf(&want, "%d\n", v)
		}
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want.String()) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed: expected\n%s\ngot\n%s\n", i+1, strings.TrimSpace(want.String()), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
