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

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func checkBFS(n int, edges [][2]int, order []int) string {
	if len(order) != n || order[0] != 1 {
		return "No"
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n+1)
	queue := []int{1}
	visited[1] = true
	idx := 1
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		// gather unvisited neighbors
		var neigh []int
		for _, v := range adj[u] {
			if !visited[v] {
				neigh = append(neigh, v)
				visited[v] = true
			}
		}
		// check that next len(neigh) numbers in order are exactly these neighbors
		set := make(map[int]bool)
		for _, v := range neigh {
			set[v] = true
		}
		for i := 0; i < len(neigh); i++ {
			if idx >= len(order) || !set[order[idx]] {
				return "No"
			}
			queue = append(queue, order[idx])
			idx++
		}
	}
	if idx != len(order) {
		return "No"
	}
	return "Yes"
}

func genCase(r *rand.Rand) (string, string) {
	n := r.Intn(9) + 2 // at least 2
	edges := make([][2]int, n-1)
	// build random tree using union
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] == x {
			return x
		}
		parent[x] = find(parent[x])
		return parent[x]
	}
	union := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false
		}
		parent[ra] = rb
		return true
	}
	idx := 0
	for idx < n-1 {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		if union(u, v) {
			edges[idx] = [2]int{u, v}
			idx++
		}
	}
	// compute bfs order
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n+1)
	order := make([]int, 0, n)
	q := []int{1}
	visited[1] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				q = append(q, v)
			}
		}
	}
	// maybe make invalid
	if r.Intn(2) == 0 {
		// shuffle order randomly to make invalid
		i := r.Intn(n-1) + 1
		j := r.Intn(n-1) + 1
		order[i], order[j] = order[j], order[i]
	}
	// compute expected
	expect := checkBFS(n, edges, order)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i, v := range order {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
