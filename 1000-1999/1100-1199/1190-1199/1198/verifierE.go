package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const INF = int64(4e18)

type Edge struct {
	to, rev int
	cap     int64
}

type Dinic struct {
	N     int
	G     [][]Edge
	level []int
	ptr   []int
}

func NewDinic(n int) *Dinic {
	d := &Dinic{N: n, G: make([][]Edge, n), level: make([]int, n), ptr: make([]int, n)}
	return d
}

func (d *Dinic) AddEdge(u, v int, c int64) {
	d.G[u] = append(d.G[u], Edge{v, len(d.G[v]), c})
	d.G[v] = append(d.G[v], Edge{u, len(d.G[u]) - 1, 0})
}

func (d *Dinic) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	q := make([]int, 0, d.N)
	d.level[s] = 0
	q = append(q, s)
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, e := range d.G[u] {
			if d.level[e.to] < 0 && e.cap > 0 {
				d.level[e.to] = d.level[u] + 1
				q = append(q, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t int, f int64) int64 {
	if u == t || f == 0 {
		return f
	}
	for ; d.ptr[u] < len(d.G[u]); d.ptr[u]++ {
		e := &d.G[u][d.ptr[u]]
		if d.level[e.to] == d.level[u]+1 && e.cap > 0 {
			pushed := d.dfs(e.to, t, min64(f, e.cap))
			if pushed > 0 {
				e.cap -= pushed
				d.G[e.to][e.rev].cap += pushed
				return pushed
			}
		}
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
	flow := int64(0)
	for d.bfs(s, t) {
		for i := range d.ptr {
			d.ptr[i] = 0
		}
		for {
			pushed := d.dfs(s, t, INF)
			if pushed == 0 {
				break
			}
			flow += pushed
		}
	}
	return flow
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	res := []int{a[0]}
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			res = append(res, a[i])
		}
	}
	return res
}

func solveE(n int, rects [][4]int) int64 {
	if len(rects) == 0 {
		return 0
	}
	xs := make([]int, 0, 2*len(rects))
	ys := make([]int, 0, 2*len(rects))
	for _, r := range rects {
		xs = append(xs, r[0])
		xs = append(xs, r[2]+1)
		ys = append(ys, r[1])
		ys = append(ys, r[3]+1)
	}
	sort.Ints(xs)
	sort.Ints(ys)
	xs = uniqueInts(xs)
	ys = uniqueInts(ys)
	nx := len(xs) - 1
	ny := len(ys) - 1
	rowLen := make([]int64, nx)
	colLen := make([]int64, ny)
	for i := 0; i < nx; i++ {
		rowLen[i] = int64(xs[i+1] - xs[i])
	}
	for j := 0; j < ny; j++ {
		colLen[j] = int64(ys[j+1] - ys[j])
	}
	has := make([][]bool, nx)
	for i := range has {
		has[i] = make([]bool, ny)
	}
	for _, r := range rects {
		x1 := sort.SearchInts(xs, r[0])
		x2 := sort.SearchInts(xs, r[2]+1)
		y1 := sort.SearchInts(ys, r[1])
		y2 := sort.SearchInts(ys, r[3]+1)
		for i := x1; i < x2; i++ {
			for j := y1; j < y2; j++ {
				has[i][j] = true
			}
		}
	}
	N := nx + ny + 2
	src := 0
	sink := N - 1
	d := NewDinic(N)
	for i := 0; i < nx; i++ {
		if rowLen[i] > 0 {
			d.AddEdge(src, 1+i, rowLen[i])
		}
	}
	for j := 0; j < ny; j++ {
		if colLen[j] > 0 {
			d.AddEdge(1+nx+j, sink, colLen[j])
		}
	}
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			if has[i][j] {
				d.AddEdge(1+i, 1+nx+j, INF)
			}
		}
	}
	return d.MaxFlow(src, sink)
}

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) ([]byte, int64) {
	n := rng.Intn(20) + 1
	m := rng.Intn(4)
	rects := make([][4]int, m)
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		x1 := rng.Intn(n) + 1
		x2 := rng.Intn(n-x1+1) + x1
		y1 := rng.Intn(n) + 1
		y2 := rng.Intn(n-y1+1) + y1
		rects[i] = [4]int{x1, y1, x2, y2}
		fmt.Fprintf(&b, "%d %d %d %d\n", x1, y1, x2, y2)
	}
	expect := solveE(n, rects)
	return b.Bytes(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(out) != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i, expect, strings.TrimSpace(out), string(input))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
