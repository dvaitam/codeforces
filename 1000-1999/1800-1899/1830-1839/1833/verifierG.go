package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type pair struct{ u, v int }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveCase(n int, edges []pair) string {
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	size := make([]int, n+1)
	visited := make([]bool, n+1)
	mp2 := make(map[pair]int)
	for i := 1; i <= n; i++ {
		size[i] = 1
	}
	for idx, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
		deg[e.u]++
		deg[e.v]++
		u, v := min(e.u, e.v), max(e.u, e.v)
		mp2[pair{u, v}] = idx + 1
	}
	if n%3 != 0 {
		return "-1\n"
	}
	queue := make([]int, 0, n)
	head := 0
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			queue = append(queue, i)
		}
	}
	st := make(map[pair]bool)
	check := false
	for head < len(queue) {
		it := queue[head]
		head++
		if size[it] > 3 {
			check = true
			break
		}
		visited[it] = true
		for _, child := range adj[it] {
			if !visited[child] {
				deg[child]--
				if size[it] == 3 {
					u, v := min(it, child), max(it, child)
					st[pair{u, v}] = true
				} else {
					size[child] += size[it]
				}
				if deg[child] == 1 {
					queue = append(queue, child)
				}
			}
		}
	}
	if check {
		return "-1\n"
	}
	keys := make([]pair, 0, len(st))
	for k := range st {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].u != keys[j].u {
			return keys[i].u < keys[j].u
		}
		return keys[i].v < keys[j].v
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(keys)))
	for i, k := range keys {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", mp2[k]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTree(rng *rand.Rand, n int) []pair {
	edges := make([]pair, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, pair{p, i})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2
	edges := generateTree(rng, n)
	input := fmt.Sprintf("1\n%d\n", n)
	for _, e := range edges {
		input += fmt.Sprintf("%d %d\n", e.u, e.v)
	}
	return input, solveCase(n, edges)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
