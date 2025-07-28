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

type Edge struct {
	u, v int
	w    int
}

type DSU struct {
	p, r []int
}

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n), r: make([]int, n)}
	for i := 0; i < n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(x, y int) bool {
	fx := d.find(x)
	fy := d.find(y)
	if fx == fy {
		return false
	}
	if d.r[fx] < d.r[fy] {
		fx, fy = fy, fx
	}
	d.p[fy] = fx
	if d.r[fx] == d.r[fy] {
		d.r[fx]++
	}
	return true
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func mstEdges(edges []Edge, n int, x int) []int {
	es := make([]Edge, len(edges))
	copy(es, edges)
	sort.Slice(es, func(i, j int) bool {
		d1 := abs(es[i].w - x)
		d2 := abs(es[j].w - x)
		if d1 == d2 {
			return es[i].w < es[j].w
		}
		return d1 < d2
	})
	dsu := NewDSU(n + 1)
	ws := make([]int, 0, n-1)
	for _, e := range es {
		if dsu.union(e.u, e.v) {
			ws = append(ws, e.w)
			if len(ws) == n-1 {
				break
			}
		}
	}
	sort.Ints(ws)
	return ws
}

type MSTInfo struct {
	weights []int
	prefix  []int64
}

func buildInfo(ws []int) MSTInfo {
	prefix := make([]int64, len(ws)+1)
	for i, w := range ws {
		prefix[i+1] = prefix[i] + int64(w)
	}
	return MSTInfo{weights: ws, prefix: prefix}
}

func (info MSTInfo) cost(x int) int64 {
	w := info.weights
	idx := sort.Search(len(w), func(i int) bool { return w[i] > x })
	sumLess := info.prefix[idx]
	total := info.prefix[len(w)]
	return int64(x)*int64(idx) - sumLess + (total - sumLess) - int64(x)*int64(len(w)-idx)
}

func runCandidate(bin, input string) (string, error) {
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(n, m int, edges []Edge, p, k, a, b, c int, q []int) int64 {
	boundariesMap := make(map[int]struct{})
	boundariesMap[0] = struct{}{}
	boundariesMap[c] = struct{}{}
	weights := make([]int, m)
	for i, e := range edges {
		weights[i] = e.w
	}
	for i := 0; i < m; i++ {
		for j := i; j < m; j++ {
			v := (weights[i] + weights[j] + 1) / 2
			if v < 0 {
				v = 0
			}
			if v > c {
				v = c
			}
			boundariesMap[v] = struct{}{}
		}
	}
	boundaries := make([]int, 0, len(boundariesMap))
	for v := range boundariesMap {
		boundaries = append(boundaries, v)
	}
	sort.Ints(boundaries)
	starts := make([]int, 0)
	for _, v := range boundaries {
		if v < c {
			starts = append(starts, v)
		}
	}
	starts = append(starts, c)
	infos := make([]MSTInfo, len(starts))
	for i, st := range starts {
		if st == c {
			infos[i] = MSTInfo{}
		} else {
			ws := mstEdges(edges, n, st)
			infos[i] = buildInfo(ws)
		}
	}
	var xorResult int64
	for _, x := range q {
		idx := sort.Search(len(starts), func(i int) bool { return starts[i] > x }) - 1
		if idx < 0 {
			idx = 0
		}
		ans := infos[idx].cost(x)
		xorResult ^= ans
	}
	return xorResult
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(5) + 2
	m := r.Intn(5) + n - 1
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		for v == u {
			v = r.Intn(n) + 1
		}
		w := r.Intn(20)
		edges[i] = Edge{u: u, v: v, w: w}
	}
	p := r.Intn(3) + 1
	k := p + r.Intn(3)
	a := r.Intn(5) + 1
	b := r.Intn(5) + 1
	c := r.Intn(20) + 1
	q := make([]int, k)
	for i := 0; i < p; i++ {
		q[i] = r.Intn(c)
	}
	for i := p; i < k; i++ {
		q[i] = (q[i-1]*a + b) % c
	}
	expect := fmt.Sprintf("%d", solveE(n, m, edges, p, k, a, b, c, q))
	input := fmt.Sprintf("%d %d\n", n, m)
	for i := 0; i < m; i++ {
		input += fmt.Sprintf("%d %d %d\n", edges[i].u, edges[i].v, edges[i].w)
	}
	input += fmt.Sprintf("%d %d %d %d %d\n", p, k, a, b, c)
	for i := 0; i < p; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", q[i])
	}
	if p > 0 {
		input += "\n"
	}
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
