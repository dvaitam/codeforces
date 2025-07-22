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

func expected(n int, edges [][2]int) float64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	q := []int{1}
	depth[1] = 1
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range adj[u] {
			if v != parent[u] {
				parent[v] = u
				depth[v] = depth[u] + 1
				q = append(q, v)
			}
		}
	}
	ans := 0.0
	for i := 1; i <= n; i++ {
		ans += 1.0 / float64(depth[i])
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(15) + 1
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{i, p}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), expected(n, edges)
}

func runCase(bin, input string, exp float64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	diff := got - exp
	if diff < 0 {
		diff = -diff
	}
	if diff > 1e-6*exp+1e-6 {
		return fmt.Errorf("expected %.10f got %.10f", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
