package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Edge struct {
	to int
	w  int64
}

type Item struct {
	v    int
	dist int64
}

type PQ []Item

func (h PQ) Len() int            { return len(h) }
func (h PQ) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h PQ) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PQ) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *PQ) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[0 : n-1]
	return it
}

type test struct {
	input    string
	expected string
}

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solve(input string) string {
	r := strings.NewReader(input)
	var n int64
	var m int
	var sx, sy, fx, fy int64
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return ""
	}
	fmt.Fscan(r, &sx, &sy, &fx, &fy)
	type Point struct {
		x, y int64
		idx  int
	}
	pts := make([]Point, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &pts[i].x, &pts[i].y)
		pts[i].idx = i
	}
	ans := abs64(sx-fx) + abs64(sy-fy)
	if m == 0 {
		return fmt.Sprint(ans)
	}
	g := make([][]Edge, m)
	sxord := make([]Point, m)
	syord := make([]Point, m)
	copy(sxord, pts)
	copy(syord, pts)
	sort.Slice(sxord, func(i, j int) bool { return sxord[i].x < sxord[j].x })
	sort.Slice(syord, func(i, j int) bool { return syord[i].y < syord[j].y })
	for i := 0; i < m-1; i++ {
		u := sxord[i].idx
		v := sxord[i+1].idx
		w := sxord[i+1].x - sxord[i].x
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
		u = syord[i].idx
		v = syord[i+1].idx
		w = syord[i+1].y - syord[i].y
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
	}
	const INF = int64(9e18)
	dist := make([]int64, m)
	for i := range dist {
		dist[i] = INF
	}
	pq := &PQ{}
	heap.Init(pq)
	for i, p := range pts {
		d := min64(abs64(sx-p.x), abs64(sy-p.y))
		dist[i] = d
		heap.Push(pq, Item{i, d})
	}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		u, d := it.v, it.dist
		if d != dist[u] {
			continue
		}
		cand := d + abs64(pts[u].x-fx) + abs64(pts[u].y-fy)
		if cand < ans {
			ans = cand
		}
		for _, e := range g[u] {
			nd := d + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}
	return fmt.Sprint(ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	tests := []test{}
	fixed := []string{
		"1 0\n1 1 3 3\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := int64(rng.Intn(10) + 5)
		m := rng.Intn(5)
		sx := int64(rng.Intn(int(n)) + 1)
		sy := int64(rng.Intn(int(n)) + 1)
		fx := int64(rng.Intn(int(n)) + 1)
		fy := int64(rng.Intn(int(n)) + 1)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n%d %d %d %d\n", n, m, sx, sy, fx, fy)
		for i := 0; i < m; i++ {
			x := int64(rng.Intn(int(n)) + 1)
			y := int64(rng.Intn(int(n)) + 1)
			fmt.Fprintf(&sb, "%d %d\n", x, y)
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
