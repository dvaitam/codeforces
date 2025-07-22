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

type edge struct{ u, v int }

func solveCase(n int, edges []edge) string {
	adj := make([][]int, n+1)
	degree := make([]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
		degree[e.u]++
		degree[e.v]++
	}
	inCycle := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		inCycle[i] = true
	}
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if degree[i] == 1 {
			queue = append(queue, i)
		}
	}
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		inCycle[u] = false
		for _, v := range adj[u] {
			if inCycle[v] {
				degree[v]--
				if degree[v] == 1 {
					queue = append(queue, v)
				}
			}
		}
	}
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	bfs := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if inCycle[i] {
			dist[i] = 0
			bfs = append(bfs, i)
		}
	}
	for head := 0; head < len(bfs); head++ {
		u := bfs[head]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				bfs = append(bfs, v)
			}
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(dist[i]))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 3 // 3..9
	edges := make([]edge, 0, n)
	cycleLen := rng.Intn(n-1) + 2
	for i := 1; i <= cycleLen; i++ {
		u := i
		v := i%cycleLen + 1
		edges = append(edges, edge{u, v})
	}
	for v := cycleLen + 1; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, edge{v, p})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	input := sb.String()
	expected := solveCase(n, edges)
	return input, expected
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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
