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
	w    int64
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	for x != d.parent[x] {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(a, b int, odd *int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	*odd -= d.size[ra] % 2
	*odd -= d.size[rb] % 2
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	*odd += d.size[ra] % 2
}

func solveCase(n, m int, edges []Edge) []int64 {
	if n%2 == 1 {
		res := make([]int64, m)
		for i := range res {
			res[i] = -1
		}
		return res
	}
	res := make([]int64, m)
	for i := 1; i <= m; i++ {
		subset := make([]Edge, i)
		copy(subset, edges[:i])
		sort.Slice(subset, func(a, b int) bool { return subset[a].w < subset[b].w })
		d := NewDSU(n)
		odd := n
		ans := int64(-1)
		for _, e := range subset {
			d.Union(e.u, e.v, &odd)
			if odd == 0 {
				ans = e.w
				break
			}
		}
		res[i-1] = ans
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	m := rng.Intn(8) + 1
	edges := make([]Edge, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		w := int64(rng.Intn(10) + 1)
		edges[i] = Edge{u: u, v: v, w: w}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
	}
	out := solveCase(n, m, edges)
	var exp strings.Builder
	for i, v := range out {
		if i > 0 {
			exp.WriteByte('\n')
		}
		exp.WriteString(fmt.Sprint(v))
	}
	exp.WriteByte('\n')
	return sb.String(), exp.String()
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
