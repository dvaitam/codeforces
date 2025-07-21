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

func countCycles(n int, edges [][2]int) int64 {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	vis := make([]bool, n)
	var ans int64
	var dfs func(start, u, depth int)
	dfs = func(start, u, depth int) {
		if depth == 4 {
			for _, w := range adj[u] {
				if w == start {
					ans++
					break
				}
			}
			return
		}
		for _, w := range adj[u] {
			if w > start && !vis[w] {
				vis[w] = true
				dfs(start, w, depth+1)
				vis[w] = false
			}
		}
	}
	for s := 0; s < n; s++ {
		vis[s] = true
		for _, v := range adj[s] {
			if v > s {
				vis[v] = true
				dfs(s, v, 1)
				vis[v] = false
			}
		}
		vis[s] = false
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 5
	edges := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float32() < 0.3 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	m := len(edges)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	ans := countCycles(n, edges)
	return sb.String(), fmt.Sprintf("%d", ans)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
