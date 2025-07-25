package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Edge struct {
	u, v int
	w    int64
	idx  int
}

type DSU struct {
	p  []int
	sz []int
}

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n+1), sz: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.p[i] = i
		d.sz[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(x, y int) bool {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return false
	}
	if d.sz[x] < d.sz[y] {
		x, y = y, x
	}
	d.p[y] = x
	d.sz[x] += d.sz[y]
	return true
}

const LOG = 20

type Node struct {
	to int
	w  int64
}

func expectedMST(n int, edges []Edge) []int64 {
	sorted := append([]Edge(nil), edges...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].w < sorted[j].w })
	dsu := NewDSU(n)
	mark := make([]bool, len(edges))
	adj := make([][]Node, n+1)
	mstWeight := int64(0)
	for _, e := range sorted {
		if dsu.Union(e.u, e.v) {
			mark[e.idx] = true
			mstWeight += e.w
			adj[e.u] = append(adj[e.u], Node{to: e.v, w: e.w})
			adj[e.v] = append(adj[e.v], Node{to: e.u, w: e.w})
		}
	}
	depth := make([]int, n+1)
	parent := make([][]int, LOG)
	maxUp := make([][]int64, LOG)
	for k := 0; k < LOG; k++ {
		parent[k] = make([]int, n+1)
		maxUp[k] = make([]int64, n+1)
	}
	queue := []int{1}
	parent[0][1] = 0
	depth[1] = 0
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		for _, nb := range adj[v] {
			if nb.to == parent[0][v] {
				continue
			}
			parent[0][nb.to] = v
			maxUp[0][nb.to] = nb.w
			depth[nb.to] = depth[v] + 1
			queue = append(queue, nb.to)
		}
	}
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			p := parent[k-1][v]
			parent[k][v] = parent[k-1][p]
			if maxUp[k-1][v] > maxUp[k-1][p] {
				maxUp[k][v] = maxUp[k-1][v]
			} else {
				maxUp[k][v] = maxUp[k-1][p]
			}
		}
	}
	getMax := func(u, v int) int64 {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		res := int64(0)
		diff := depth[u] - depth[v]
		for k := LOG - 1; k >= 0; k-- {
			if diff>>k&1 == 1 {
				if maxUp[k][u] > res {
					res = maxUp[k][u]
				}
				u = parent[k][u]
			}
		}
		if u == v {
			return res
		}
		for k := LOG - 1; k >= 0; k-- {
			if parent[k][u] != parent[k][v] {
				if maxUp[k][u] > res {
					res = maxUp[k][u]
				}
				if maxUp[k][v] > res {
					res = maxUp[k][v]
				}
				u = parent[k][u]
				v = parent[k][v]
			}
		}
		if maxUp[0][u] > res {
			res = maxUp[0][u]
		}
		if maxUp[0][v] > res {
			res = maxUp[0][v]
		}
		return res
	}
	ans := make([]int64, len(edges))
	for _, e := range edges {
		if mark[e.idx] {
			ans[e.idx] = mstWeight
		} else {
			maxW := getMax(e.u, e.v)
			ans[e.idx] = mstWeight + e.w - maxW
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		pos := 0
		n, _ := strconv.Atoi(parts[pos])
		pos++
		m, _ := strconv.Atoi(parts[pos])
		pos++
		if len(parts) < pos+3*m {
			fmt.Printf("test %d invalid length\n", idx)
			os.Exit(1)
		}
		edges := make([]Edge, m)
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(parts[pos])
			pos++
			v, _ := strconv.Atoi(parts[pos])
			pos++
			w, _ := strconv.ParseInt(parts[pos], 10, 64)
			pos++
			edges[i] = Edge{u: u, v: v, w: w, idx: i}
		}
		exp := expectedMST(n, edges)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d", n, m)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&buf, " %d %d %d", edges[i].u, edges[i].v, edges[i].w)
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outParts := strings.Fields(strings.TrimSpace(string(out)))
		if len(outParts) < m {
			fmt.Printf("Test %d invalid output length\n", idx)
			os.Exit(1)
		}
		for i := 0; i < m; i++ {
			got, _ := strconv.ParseInt(outParts[i], 10, 64)
			if got != exp[i] {
				fmt.Printf("Test %d failed at edge %d: expected %d got %d\n", idx, i+1, exp[i], got)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
