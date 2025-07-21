package main

import (
	"bufio"
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
	g, s int64
	idxG int
}

type DSU struct {
	p, r []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	r := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
	for d.p[x] != x {
		d.p[x] = d.p[d.p[x]]
		x = d.p[x]
	}
	return x
}

func (d *DSU) Union(x, y int) bool {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return false
	}
	if d.r[rx] < d.r[ry] {
		d.p[rx] = ry
	} else if d.r[rx] > d.r[ry] {
		d.p[ry] = rx
	} else {
		d.p[ry] = rx
		d.r[rx]++
	}
	return true
}

func solve(data string) string {
	in := bufio.NewReader(strings.NewReader(data))
	var N, M int
	var G, S int64
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return ""
	}
	fmt.Fscan(in, &G, &S)
	edges := make([]Edge, M)
	for i := 0; i < M; i++ {
		var x, y int
		var gi, si int64
		fmt.Fscan(in, &x, &y, &gi, &si)
		edges[i] = Edge{u: x - 1, v: y - 1, g: gi, s: si}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].g < edges[j].g })
	for i := range edges {
		edges[i].idxG = i
	}
	d0 := NewDSU(N)
	comps := N
	i0 := -1
	for i, e := range edges {
		if d0.Union(e.u, e.v) {
			comps--
			if comps == 1 {
				i0 = i
				break
			}
		}
	}
	if i0 < 0 {
		return "-1\n"
	}
	edgesS := make([]Edge, M)
	copy(edgesS, edges)
	sort.Slice(edgesS, func(i, j int) bool { return edgesS[i].s < edgesS[j].s })
	const INF = int64(4e18)
	ans := INF
	for i := i0; i < M; i++ {
		a := edges[i].g
		if a*G >= ans {
			break
		}
		dsu := NewDSU(N)
		used := 0
		var b int64
		for _, e := range edgesS {
			if e.idxG > i {
				continue
			}
			if dsu.Union(e.u, e.v) {
				used++
				if e.s > b {
					b = e.s
				}
				if used == N-1 {
					break
				}
			}
		}
		if used == N-1 {
			cost := a*G + b*S
			if cost < ans {
				ans = cost
			}
		}
	}
	if ans == INF {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	N := rng.Intn(5) + 2
	maxM := N * N
	M := rng.Intn(maxM) + 1
	if M > 15 {
		M = 15
	}
	G := rng.Int63n(5) + 1
	S := rng.Int63n(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", N, M))
	sb.WriteString(fmt.Sprintf("%d %d\n", G, S))
	for i := 0; i < M; i++ {
		u := rng.Intn(N) + 1
		v := rng.Intn(N) + 1
		g := rng.Int63n(10) + 1
		s := rng.Int63n(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u, v, g, s))
	}
	input := sb.String()
	expected := solve(input)
	return input, expected
}

func runCase(exe string, in, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
