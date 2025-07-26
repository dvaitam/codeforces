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
	to int
	id int
}

func solve(n int, edges [][2]int) []int {
	m := len(edges)
	g := make([][]Edge, n+1)
	for i, e := range edges {
		u, v := e[0], e[1]
		id := i + 1
		g[u] = append(g[u], Edge{v, id})
		g[v] = append(g[v], Edge{u, id})
	}
	disc := make([]int, n+1)
	low := make([]int, n+1)
	isBridge := make([]bool, m+1)
	timer := 0
	var dfs func(u, parentEdge int)
	dfs = func(u, parentEdge int) {
		timer++
		disc[u] = timer
		low[u] = timer
		for _, e := range g[u] {
			if e.id == parentEdge {
				continue
			}
			v := e.to
			if disc[v] == 0 {
				dfs(v, e.id)
				if low[v] < low[u] {
					low[u] = low[v]
				}
				if low[v] > disc[u] {
					isBridge[e.id] = true
				}
			} else {
				if disc[v] < low[u] {
					low[u] = disc[v]
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if disc[i] == 0 {
			dfs(i, -1)
		}
	}
	visited := make([]bool, n+1)
	edgeVisited := make([]bool, m+1)
	var ans []int
	var dfs2 func(int, *int, *[]int)
	dfs2 = func(u int, vertexCount *int, edges *[]int) {
		visited[u] = true
		*vertexCount = *vertexCount + 1
		for _, e := range g[u] {
			if isBridge[e.id] {
				continue
			}
			if !edgeVisited[e.id] {
				edgeVisited[e.id] = true
				*edges = append(*edges, e.id)
			}
			if !visited[e.to] {
				dfs2(e.to, vertexCount, edges)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			tempEdges := []int{}
			count := 0
			dfs2(i, &count, &tempEdges)
			if len(tempEdges) == count {
				ans = append(ans, tempEdges...)
			}
		}
	}
	sort.Ints(ans)
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		e := [2]int{u, v}
		if used[e] {
			continue
		}
		used[e] = true
		edges = append(edges, e)
	}
	input := fmt.Sprintf("%d %d\n", n, m)
	for _, e := range edges {
		input += fmt.Sprintf("%d %d\n", e[0], e[1])
	}
	ans := solve(n, edges)
	expected := fmt.Sprintf("%d\n", len(ans))
	for i, id := range ans {
		if i > 0 {
			expected += " "
		}
		expected += fmt.Sprintf("%d", id)
	}
	if len(ans) > 0 {
		expected += "\n"
	}
	return input, expected
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, out.String())
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
