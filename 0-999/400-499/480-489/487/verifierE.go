package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Embedded reference solver for 487E (block-cut tree + HLD + segment tree)

const refINF = int64(1) << 60

type refEdge struct{ u, v int }

func refMinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func refMin64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solveE(n, m, q int, weights []int, edgeList [][2]int, queries []string) []string {
	wgt := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		wgt[i] = int64(weights[i-1])
	}
	adj := make([][]int, n+1)
	for _, e := range edgeList {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}

	dfn := make([]int, n+1)
	low := make([]int, n+1)
	timer := 0
	stk := []refEdge{}
	mark := make([]int, n+1)
	curMark := 0
	bcCnt := 0

	totMax := n + m + 5
	btAdj := make([][]int, totMax)

	var tarjan func(u, pe int)
	tarjan = func(u, pe int) {
		timer++
		dfn[u] = timer
		low[u] = timer
		for _, v := range adj[u] {
			if v == pe {
				pe = -1
				continue
			}
			if dfn[v] == 0 {
				stk = append(stk, refEdge{u, v})
				tarjan(v, u)
				low[u] = refMinInt(low[u], low[v])
				if low[v] >= dfn[u] {
					bcCnt++
					bid := n + bcCnt
					curMark++
					for {
						e := stk[len(stk)-1]
						stk = stk[:len(stk)-1]
						for _, x := range []int{e.u, e.v} {
							if mark[x] != curMark {
								mark[x] = curMark
								btAdj[bid] = append(btAdj[bid], x)
								btAdj[x] = append(btAdj[x], bid)
							}
						}
						if e.u == u && e.v == v {
							break
						}
					}
				}
			} else if dfn[v] < dfn[u] {
				stk = append(stk, refEdge{u, v})
				low[u] = refMinInt(low[u], dfn[v])
			}
		}
	}

	for i := 1; i <= n; i++ {
		if dfn[i] == 0 {
			tarjan(i, -1)
		}
	}

	totN := n + bcCnt
	parent := make([]int, totN+1)
	depth := make([]int, totN+1)
	heavy := make([]int, totN+1)
	sz := make([]int, totN+1)

	var dfs1 func(u, p int)
	dfs1 = func(u, p int) {
		sz[u] = 1
		parent[u] = p
		depth[u] = depth[p] + 1
		maxSz := 0
		for _, v := range btAdj[u] {
			if v == p {
				continue
			}
			dfs1(v, u)
			if sz[v] > maxSz {
				maxSz = sz[v]
				heavy[u] = v
			}
			sz[u] += sz[v]
		}
	}
	dfs1(1, 0)

	headArr := make([]int, totN+1)
	pos := make([]int, totN+1)
	flatW := make([]int64, totN+1)
	curPos := 1

	var dfs2 func(u, h int)
	dfs2 = func(u, h int) {
		headArr[u] = h
		pos[u] = curPos
		if u <= n {
			flatW[curPos] = wgt[u]
		} else {
			flatW[curPos] = refINF
		}
		curPos++
		if heavy[u] != 0 {
			dfs2(heavy[u], h)
		}
		for _, v := range btAdj[u] {
			if v == parent[u] || v == heavy[u] {
				continue
			}
			dfs2(v, v)
		}
	}
	dfs2(1, 1)

	// Segment tree
	segN := 1
	for segN < totN {
		segN <<= 1
	}
	segT := make([]int64, 2*segN)
	for i := range segT {
		segT[i] = refINF
	}
	for i := 1; i <= totN; i++ {
		segT[segN+i-1] = flatW[i]
	}
	for i := segN - 1; i >= 1; i-- {
		segT[i] = refMin64(segT[2*i], segT[2*i+1])
	}

	segUpdate := func(p int, v int64) {
		i := segN + p - 1
		segT[i] = v
		for i >>= 1; i > 0; i >>= 1 {
			segT[i] = refMin64(segT[2*i], segT[2*i+1])
		}
	}

	segQuery := func(l, r int) int64 {
		res := refINF
		l += segN - 1
		r += segN - 1
		for l <= r {
			if l&1 == 1 {
				res = refMin64(res, segT[l])
				l++
			}
			if r&1 == 0 {
				res = refMin64(res, segT[r])
				r--
			}
			l >>= 1
			r >>= 1
		}
		return res
	}

	pathQuery := func(u, v int) int64 {
		res := refINF
		for headArr[u] != headArr[v] {
			if depth[headArr[u]] > depth[headArr[v]] {
				res = refMin64(res, segQuery(pos[headArr[u]], pos[u]))
				u = parent[headArr[u]]
			} else {
				res = refMin64(res, segQuery(pos[headArr[v]], pos[v]))
				v = parent[headArr[v]]
			}
		}
		l, r := pos[u], pos[v]
		if l > r {
			l, r = r, l
		}
		res = refMin64(res, segQuery(l, r))
		return res
	}

	var results []string
	for _, line := range queries {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		typ := fields[0]
		a := parseInt(fields[1])
		b := parseInt(fields[2])
		if typ == "C" {
			wgt[a] = int64(b)
			segUpdate(pos[a], int64(b))
		} else if typ == "A" {
			ans := pathQuery(a, b)
			results = append(results, fmt.Sprint(ans))
		}
	}
	return results
}

func parseInt(s string) int {
	v := 0
	for _, c := range s {
		v = v*10 + int(c-'0')
	}
	return v
}

func solveEFromInput(input string) string {
	fields := strings.Fields(input)
	pos := 0
	next := func() string {
		s := fields[pos]
		pos++
		return s
	}
	nextInt := func() int { return parseInt(next()) }

	n := nextInt()
	m := nextInt()
	q := nextInt()

	weights := make([]int, n)
	for i := 0; i < n; i++ {
		weights[i] = nextInt()
	}

	edgeList := make([][2]int, m)
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		edgeList[i] = [2]int{u, v}
	}

	queries := make([]string, q)
	for i := 0; i < q; i++ {
		typ := next()
		a := next()
		b := next()
		queries[i] = typ + " " + a + " " + b
	}

	results := solveE(n, m, q, weights, edgeList, queries)
	return strings.Join(results, "\n")
}

func genGraph(rng *rand.Rand, n, m int) [][2]int {
	edges := make([][2]int, 0, m)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		exist := false
		for _, e := range edges {
			if (e[0] == a && e[1] == b) || (e[0] == b && e[1] == a) {
				exist = true
				break
			}
		}
		if !exist {
			edges = append(edges, [2]int{a, b})
		}
	}
	return edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if !filepath.IsAbs(bin) {
		abs, err := filepath.Abs(bin)
		if err == nil {
			bin = abs
		}
	}
	rng := rand.New(rand.NewSource(42))
	for t := 0; t < 30; t++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(3) + n - 1
		q := rng.Intn(5) + 1
		weights := make([]int, n)
		for i := 0; i < n; i++ {
			weights[i] = rng.Intn(100) + 1
		}
		edges := genGraph(rng, n, m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", weights[i]))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		queryStrs := make([]string, q)
		for i := 0; i < q; i++ {
			if rng.Intn(2) == 0 {
				a := rng.Intn(n) + 1
				b := rng.Intn(n) + 1
				line := fmt.Sprintf("A %d %d", a, b)
				sb.WriteString(line + "\n")
				queryStrs[i] = line
			} else {
				a := rng.Intn(n) + 1
				w := rng.Intn(100) + 1
				line := fmt.Sprintf("C %d %d", a, w)
				sb.WriteString(line + "\n")
				queryStrs[i] = line
			}
		}
		input := sb.String()

		expected := solveEFromInput(input)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: runtime error: %v\n%s\ninput:\n%s", t+1, err, stderr.String(), input)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
