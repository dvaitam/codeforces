package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseB struct {
	n     int
	edges [][2]int
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

func generateCase() testCaseB {
	n := rand.Intn(19) + 2 // 2..20
	edges := make([][2]int, 0, n)
	deg := make([]int, n+1)
	for i := 1; i < n; i++ {
		if rand.Intn(2) == 0 {
			u, v := i, i+1
			edges = append(edges, [2]int{u, v})
			deg[u]++
			deg[v]++
		}
	}
	if rand.Intn(2) == 0 {
		u, v := 1, n
		if deg[u] < 2 && deg[v] < 2 {
			edges = append(edges, [2]int{u, v})
			deg[u]++
			deg[v]++
		}
	}
	if len(edges) == 0 {
		edges = append(edges, [2]int{1, 2})
	}
	return testCaseB{n: n, edges: edges}
}

func compute(tc testCaseB) int {
	n := tc.n
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		u := e[0]
		v := e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n+1)
	bench := 0
	oddPath := 0
	for i := 1; i <= n; i++ {
		if visited[i] {
			continue
		}
		stack := []int{i}
		visited[i] = true
		nodes := 0
		degSum := 0
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nodes++
			degSum += len(adj[u])
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					stack = append(stack, v)
				}
			}
		}
		edgesCnt := degSum / 2
		if edgesCnt == nodes {
			if nodes%2 == 1 {
				bench++
			}
		} else {
			if nodes%2 == 1 {
				oddPath++
			}
		}
	}
	if oddPath%2 == 1 {
		bench++
	}
	return bench
}

func buildInput(tc testCaseB) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		tc := generateCase()
		input := buildInput(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", i)
			os.Exit(1)
		}
		exp := compute(tc)
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\n", i, exp, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
