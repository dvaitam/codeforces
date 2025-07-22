package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Edge struct {
	to   int
	cost int64
	ramp int
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type Ramp struct {
	x, d, t, p int64
	idx        int
}

func expected(n int, L int64, ramps []Ramp) string {
	coords := make([]int64, 0, 3*n+2)
	coords = append(coords, 0, L)
	usable := make([]Ramp, 0, n)
	for _, r := range ramps {
		if r.x-r.p < 0 || r.x+r.d > L {
			continue
		}
		usable = append(usable, r)
		coords = append(coords, r.x-r.p, r.x, r.x+r.d)
	}
	sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
	coords = unique(coords)
	m := len(coords)
	idxOf := func(v int64) int {
		return sort.Search(len(coords), func(i int) bool { return coords[i] >= v })
	}
	g := make([][]Edge, m)
	for i := 0; i+1 < m; i++ {
		w := coords[i+1] - coords[i]
		g[i] = append(g[i], Edge{i + 1, w, 0})
		g[i+1] = append(g[i+1], Edge{i, w, 0})
	}
	for _, r := range usable {
		u := idxOf(r.x - r.p)
		v := idxOf(r.x + r.d)
		cost := r.p + r.t
		g[u] = append(g[u], Edge{v, cost, r.idx})
	}
	const INF = int64(4e18)
	dist := make([]int64, m)
	prevNode := make([]int, m)
	prevRamp := make([]int, m)
	for i := range dist {
		dist[i] = INF
		prevNode[i] = -1
	}
	src := idxOf(0)
	dst := idxOf(L)
	dist[src] = 0
	pq := &PriorityQueue{{src, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		u := it.node
		if it.dist != dist[u] {
			continue
		}
		if u == dst {
			break
		}
		for _, e := range g[u] {
			nd := it.dist + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				prevNode[e.to] = u
				prevRamp[e.to] = e.ramp
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}
	path := make([]int, 0)
	for u := dst; u != src; u = prevNode[u] {
		if u < 0 {
			break
		}
		if prevRamp[u] != 0 {
			path = append(path, prevRamp[u])
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", dist[dst]))
	sb.WriteString(fmt.Sprintf("%d\n", len(path)))
	if len(path) > 0 {
		for i, v := range path {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
	}
	return strings.TrimSpace(sb.String())
}

func unique(a []int64) []int64 {
	j := 0
	for i := 0; i < len(a); i++ {
		if i == 0 || a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	tests := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		n, _ := strconv.Atoi(parts[0])
		L, _ := strconv.ParseInt(parts[1], 10, 64)
		if len(parts) != 2+4*n {
			fmt.Printf("test %d: invalid line\n", tests+1)
			os.Exit(1)
		}
		ramps := make([]Ramp, n)
		idx := 2
		for i := 0; i < n; i++ {
			x, _ := strconv.ParseInt(parts[idx], 10, 64)
			idx++
			d, _ := strconv.ParseInt(parts[idx], 10, 64)
			idx++
			t, _ := strconv.ParseInt(parts[idx], 10, 64)
			idx++
			p, _ := strconv.ParseInt(parts[idx], 10, 64)
			idx++
			ramps[i] = Ramp{x, d, t, p, i + 1}
		}
		expect := expected(n, L, ramps)
		// build input
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, L))
		for i := 0; i < n; i++ {
			r := ramps[i]
			input.WriteString(fmt.Sprintf("%d %d %d %d\n", r.x, r.d, r.t, r.p))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tests+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", tests+1, expect, got)
			os.Exit(1)
		}
		tests++
	}
	if err := sc.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", tests)
}
