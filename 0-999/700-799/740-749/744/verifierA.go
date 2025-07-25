package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &DSU{p, s}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
	}
	d.parent[ry] = rx
	d.size[rx] += d.size[ry]
}

func solveA(n int, edges [][2]int, gov []int) int64 {
	dsu := NewDSU(n)
	for _, e := range edges {
		dsu.Union(e[0], e[1])
	}
	govRoot := make(map[int]bool)
	var sizes []int64
	for _, g := range gov {
		r := dsu.Find(g)
		if !govRoot[r] {
			govRoot[r] = true
			sizes = append(sizes, int64(dsu.size[r]))
		}
	}
	var extra int64
	for i := 1; i <= n; i++ {
		if dsu.Find(i) == i && !govRoot[i] {
			extra += int64(dsu.size[i])
		}
	}
	maxIdx := 0
	for i := range sizes {
		if sizes[i] > sizes[maxIdx] {
			maxIdx = i
		}
	}
	var total int64
	for i, sz := range sizes {
		if i == maxIdx {
			s := sz + extra
			total += s * (s - 1) / 2
		} else {
			total += sz * (sz - 1) / 2
		}
	}
	return total - int64(len(edges))
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(8) + 1
	k := rng.Intn(n) + 1
	nodes := rng.Perm(n)[:k]
	gov := make([]int, k)
	for i := 0; i < k; i++ {
		gov[i] = nodes[i] + 1
	}
	dsu := NewDSU(n)
	govComp := make(map[int]bool)
	for _, g := range gov {
		govComp[g] = true
	}
	edges := make([][2]int, 0)
	used := make(map[[2]int]bool)
	attempts := n * n
	for attempts > 0 && len(edges) < n*(n-1)/2 {
		attempts--
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if used[[2]int{a, b}] {
			continue
		}
		ru := dsu.Find(u)
		rv := dsu.Find(v)
		if ru == rv {
			continue
		}
		if govComp[ru] && govComp[rv] {
			continue
		}
		used[[2]int{a, b}] = true
		edges = append(edges, [2]int{u, v})
		dsu.Union(u, v)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), k))
	for i, g := range gov {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(g))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	ans := solveA(n, edges, gov)
	return sb.String(), ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, expect := genCase(rng)
		outStr, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(outStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output %q\n", t, outStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", t, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
