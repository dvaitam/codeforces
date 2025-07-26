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

type edge struct{ to, w int }

func bfs(start int, adj [][]edge) (node, dist int) {
	n := len(adj)
	distArr := make([]int, n)
	for i := range distArr {
		distArr[i] = -1
	}
	q := []int{start}
	distArr[start] = 0
	node = start
	for idx := 0; idx < len(q); idx++ {
		u := q[idx]
		for _, e := range adj[u] {
			if distArr[e.to] == -1 {
				distArr[e.to] = distArr[u] + e.w
				q = append(q, e.to)
				if distArr[e.to] > distArr[node] {
					node = e.to
				}
			}
		}
	}
	return node, distArr[node]
}

func expected(n int, edges [][3]int) int {
	adj := make([][]edge, n+1)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adj[u] = append(adj[u], edge{v, w})
		adj[v] = append(adj[v], edge{u, w})
	}
	u, _ := bfs(1, adj)
	_, d := bfs(u, adj)
	return d
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	edges := make([][3]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := rng.Intn(10) + 1
		edges[i-2] = [3]int{p, i, w}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	return sb.String(), expected(n, edges)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
