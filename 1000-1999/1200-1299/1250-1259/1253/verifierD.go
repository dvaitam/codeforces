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

type dsu struct {
	parent []int
	right  []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	r := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		r[i] = i
	}
	return &dsu{p, r}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.right[ra] < d.right[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.right[rb] > d.right[ra] {
		d.right[ra] = d.right[rb]
	}
}

func expectedD(n int, edges [][2]int) int {
	dsu := newDSU(n)
	for _, e := range edges {
		dsu.union(e[0], e[1])
	}
	ans := 0
	i := 1
	for i <= n {
		r := dsu.right[dsu.find(i)]
		j := i + 1
		for j <= r {
			if dsu.find(j) != dsu.find(i) {
				dsu.union(i, j)
				ans++
			}
			if dsu.right[dsu.find(i)] > r {
				r = dsu.right[dsu.find(i)]
			}
			j++
		}
		i = r + 1
	}
	return ans
}

func generateCaseD(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(6) + 1
	m := rng.Intn(n*(n-1)/2 + 1)
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, key)
	}
	return n, edges
}

func runCaseD(bin string, n int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expect := expectedD(n, edges)
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, edges := generateCaseD(rng)
		if err := runCaseD(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
