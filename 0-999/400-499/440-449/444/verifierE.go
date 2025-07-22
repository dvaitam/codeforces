package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Edge struct{ u, v, w int }

type DSU struct{ p []int }

func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = -1
	}
	return &DSU{p}
}
func (d *DSU) Find(x int) int {
	if d.p[x] < 0 {
		return x
	}
	d.p[x] = d.Find(d.p[x])
	return d.p[x]
}
func (d *DSU) Union(x, y int) {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return
	}
	if d.p[x] > d.p[y] {
		x, y = y, x
	}
	d.p[x] += d.p[y]
	d.p[y] = x
}

func check(T int, n int, edges []Edge, x []int, total int) bool {
	d := NewDSU(n)
	for _, e := range edges {
		if e.w >= T {
			break
		}
		d.Union(e.u, e.v)
	}
	compSize := make(map[int]int)
	compCap := make(map[int]int)
	for i := 0; i < n; i++ {
		r := d.Find(i)
		compSize[r]++
		compCap[r] += x[i]
	}
	for r, sz := range compSize {
		if sz+compCap[r] > total {
			return false
		}
	}
	return true
}

func expected(n int, edges []Edge, x []int) int {
	weights := []int{0}
	for _, e := range edges {
		weights = append(weights, e.w)
	}
	sort.Ints(weights)
	uniq := weights[:0]
	for _, v := range weights {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	total := 0
	for _, v := range x {
		total += v
	}
	lo, hi := 0, len(uniq)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if check(uniq[mid], n, edges, x, total) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return uniq[lo]
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 1
	edges := make([]Edge, n-1)
	for i := 1; i < n; i++ {
		u := rng.Intn(i)
		v := i
		w := rng.Intn(10) + 1
		edges[i-1] = Edge{u, v, w}
	}
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u+1, e.v+1, e.w))
	}
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", x[i]))
	}
	ans := expected(n, edges, x)
	return sb.String(), ans
}

func runCase(bin string, input string, expect int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
