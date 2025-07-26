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
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func expected(n, m, c int, initEdges [][3]int, events []string) string {
	dsu := NewDSU(n)
	adj := make([]map[int][]int, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = make(map[int][]int)
	}
	addEdge := func(x, y, z int) {
		if lst, ok := adj[x][z]; ok {
			for _, v := range lst {
				dsu.Union(v, y)
			}
		}
		if lst, ok := adj[y][z]; ok {
			for _, u := range lst {
				dsu.Union(u, x)
			}
		}
		adj[x][z] = append(adj[x][z], y)
		adj[y][z] = append(adj[y][z], x)
	}
	for _, e := range initEdges {
		addEdge(e[0], e[1], e[2])
	}
	var out strings.Builder
	for _, line := range events {
		parts := strings.Fields(line)
		if parts[0] == "+" {
			x := atoi(parts[1])
			y := atoi(parts[2])
			z := atoi(parts[3])
			addEdge(x, y, z)
		} else {
			x := atoi(parts[1])
			y := atoi(parts[2])
			if dsu.Find(x) == dsu.Find(y) {
				out.WriteString("Yes\n")
			} else {
				out.WriteString("No\n")
			}
		}
	}
	res := out.String()
	if len(res) > 0 && res[len(res)-1] == '\n' {
		res = res[:len(res)-1]
	}
	return res
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	c := rng.Intn(3) + 1
	m := rng.Intn(3)
	q := rng.Intn(5) + 1
	used := make(map[[2]int]bool)
	var initEdges [][3]int
	for len(initEdges) < m {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		if x == y {
			continue
		}
		if x > y {
			x, y = y, x
		}
		if used[[2]int{x, y}] {
			continue
		}
		used[[2]int{x, y}] = true
		z := rng.Intn(c) + 1
		initEdges = append(initEdges, [3]int{x, y, z})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, c, q))
	for _, e := range initEdges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	events := make([]string, q)
	hasQuery := false
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			// add edge
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			for x == y {
				y = rng.Intn(n) + 1
			}
			z := rng.Intn(c) + 1
			events[i] = fmt.Sprintf("+ %d %d %d", x, y, z)
		} else {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			for x == y {
				y = rng.Intn(n) + 1
			}
			events[i] = fmt.Sprintf("? %d %d", x, y)
			hasQuery = true
		}
	}
	if !hasQuery {
		// ensure at least one query
		events[q-1] = fmt.Sprintf("? %d %d", 1, n)
		hasQuery = true
	}
	for _, line := range events {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String(), expected(n, m, c, initEdges, events)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
