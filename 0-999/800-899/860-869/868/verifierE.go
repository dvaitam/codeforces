package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Edge struct {
	to int
	w  int
}

type Test struct {
	in  string
	out string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func oracle(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	g := make([][]Edge, n+1)
	total := 0
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
		total += w
	}
	var s int
	fmt.Fscan(in, &s)
	var m int
	fmt.Fscan(in, &m)
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(in, &x)
		_ = x
	}
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{s}
	dist[s] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g[v] {
			if dist[e.to] == -1 {
				dist[e.to] = dist[v] + e.w
				q = append(q, e.to)
			}
		}
	}
	maxd := 0
	for i := 1; i <= n; i++ {
		if dist[i] > maxd {
			maxd = dist[i]
		}
	}
	ans := 2*total - maxd
	return fmt.Sprintf("%d", ans)
}

func genTree(rng *rand.Rand, n int) []struct{ u, v, w int } {
	edges := make([]struct{ u, v, w int }, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := rng.Intn(i-1) + 1
		w := rng.Intn(5) + 1
		edges = append(edges, struct{ u, v, w int }{parent, i, w})
	}
	return edges
}

func genCase(rng *rand.Rand) Test {
	n := rng.Intn(5) + 1
	edges := genTree(rng, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	s := rng.Intn(n) + 1
	fmt.Fprintf(&sb, "%d\n", s)
	m := rng.Intn(n) + 1
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n) + 1
		for x == s {
			x = rng.Intn(n) + 1
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", x)
	}
	sb.WriteByte('\n')
	input := sb.String()
	out := oracle(input)
	return Test{input, out}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
