package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveE(n int, edges [][2]int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	queue := make([]int, n)
	head, tail := 0, 0
	queue[tail] = 1
	tail++
	parent[1] = 0
	order := []int{1}
	for head < tail {
		u := queue[head]
		head++
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			queue[tail] = v
			tail++
			order = append(order, v)
		}
	}
	leafCount := make([]int, n+1)
	maxDepthLeaf := 0
	ans := 0
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		if u != 1 && len(adj[u]) == 1 {
			leafCount[u] = 1
			if depth[u] > maxDepthLeaf {
				maxDepthLeaf = depth[u]
			}
		} else {
			sum := 0
			for _, v := range adj[u] {
				if v == parent[u] {
					continue
				}
				sum += leafCount[v]
			}
			leafCount[u] = sum
		}
		if u != 1 && leafCount[u] > 1 {
			t := depth[u] + leafCount[u]
			if t > ans {
				ans = t
			}
		}
	}
	if maxDepthLeaf > ans {
		ans = maxDepthLeaf
	}
	return ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	type Test struct {
		n     int
		edges [][2]int
	}
	var tests []Test
	tests = append(tests, Test{2, [][2]int{{1, 2}}})
	tests = append(tests, Test{3, [][2]int{{1, 2}, {1, 3}}})
	for len(tests) < 100 {
		n := rand.Intn(20) + 2
		tests = append(tests, Test{n, randTree(n)})
	}
	for idx, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t.n))
		for _, e := range t.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
		expected := fmt.Sprintf("%d", solveE(t.n, t.edges))
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed:\ninput:\n%sexpected %s got %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
