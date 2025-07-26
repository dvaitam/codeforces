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

const modD int64 = 1000000007

func countPaths(n int, edges [][2]int) int64 {
	adj := make([][]int, n)
	for i := 2; i <= n; i++ {
		p := i / 2
		adj[i-1] = append(adj[i-1], p-1)
		adj[p-1] = append(adj[p-1], i-1)
	}
	for _, e := range edges {
		u := e[0] - 1
		v := e[1] - 1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var res int64
	visited := make([]bool, n)
	var dfs func(int)
	dfs = func(u int) {
		res++
		if res >= modD {
			res -= modD
		}
		visited[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v)
			}
		}
		visited[u] = false
	}
	for i := 0; i < n; i++ {
		dfs(i)
	}
	return res % modD
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		m := rng.Intn(5)
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for v == u {
				v = rng.Intn(n) + 1
			}
			edges[j] = [2]int{u, v}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
		expected := fmt.Sprintf("%d", countPaths(n, edges))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
