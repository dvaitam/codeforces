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

func expected(n, m int, edges [][2]int) string {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n)
	color := make([]int, n)
	dist := make([]int, n)
	queue := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		comp := []int{}
		head, tail := 0, 0
		queue[tail] = i
		tail++
		visited[i] = true
		color[i] = 0
		comp = append(comp, i)
		ok := true
		for head < tail && ok {
			u := queue[head]
			head++
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					color[v] = 1 - color[u]
					queue[tail] = v
					tail++
					comp = append(comp, v)
				} else if color[v] == color[u] {
					ok = false
					break
				}
			}
		}
		if !ok {
			return "-1"
		}
		for _, u := range comp {
			dist[u] = -1
		}
		start := comp[0]
		dist[start] = 0
		head, tail = 0, 0
		queue[tail] = start
		tail++
		far := start
		for head < tail {
			u := queue[head]
			head++
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue[tail] = v
					tail++
				}
			}
			if dist[u] > dist[far] {
				far = u
			}
		}
		for _, u := range comp {
			dist[u] = -1
		}
		dist[far] = 0
		head, tail = 0, 0
		queue[tail] = far
		tail++
		diam := 0
		for head < tail {
			u := queue[head]
			head++
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue[tail] = v
					tail++
				}
			}
			if dist[u] > diam {
				diam = dist[u]
			}
		}
		total += diam
	}
	return fmt.Sprintf("%d", total)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges + 1)
		edgeSet := make(map[[2]int]struct{})
		edges := make([][2]int, 0, m)
		for len(edges) < m {
			u := rng.Intn(n)
			v := rng.Intn(n)
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if _, ok := edgeSet[key]; ok {
				continue
			}
			edgeSet[key] = struct{}{}
			edges = append(edges, key)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		input := sb.String()
		exp := expected(n, m, edges)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
