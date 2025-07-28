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

type edge struct{ u, v int }

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func construct(n, d12, d23, d13 int) (bool, [][2]int) {
	edges := make([][2]int, 0)
	if d23 == d12+d13 {
		if d12+d13+1 > n {
			return false, nil
		}
		node := 4
		second := 1
		for i := 1; i <= d12-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 2})
		second = 1
		for i := 1; i <= d13-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 3})
		for node <= n {
			edges = append(edges, [2]int{1, node})
			node++
		}
		return true, edges
	}
	if d12 == d23+d13 {
		if d23+d13+1 > n {
			return false, nil
		}
		node := 4
		second := 3
		for i := 1; i <= d13-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 1})
		second = 3
		for i := 1; i <= d23-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 2})
		for node <= n {
			edges = append(edges, [2]int{1, node})
			node++
		}
		return true, edges
	}
	if d13 == d12+d23 {
		if d12+d23+1 > n {
			return false, nil
		}
		node := 4
		second := 2
		for i := 1; i <= d12-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 1})
		second = 2
		for i := 1; i <= d23-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 3})
		for node <= n {
			edges = append(edges, [2]int{1, node})
			node++
		}
		return true, edges
	}
	for j := 1; j < d12; j++ {
		d41 := j
		d42 := d12 - j
		d43 := d13 - j
		if d43 <= 0 || d42 <= 0 || d42+d43 != d23 || d41+d42+d43+1 > n {
			continue
		}
		node := 5
		second := 4
		edges = edges[:0]
		for i := 1; i <= d41-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 1})
		second = 4
		for i := 1; i <= d43-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 3})
		second = 4
		for i := 1; i <= d42-1; i++ {
			first := second
			second = node
			edges = append(edges, [2]int{first, second})
			node++
		}
		edges = append(edges, [2]int{second, 2})
		for node <= n {
			edges = append(edges, [2]int{4, node})
			node++
		}
		return true, edges
	}
	return false, nil
}

func bfs(n int, adj [][]int, start int) []int {
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func verify(out string, n, d12, d23, d13 int, possible bool) error {
	fields := strings.Fields(out)
	if !possible {
		if len(fields) != 1 || strings.ToUpper(fields[0]) != "NO" {
			return fmt.Errorf("expected NO")
		}
		return nil
	}
	if len(fields) != 1+2*(n-1) {
		return fmt.Errorf("expected %d numbers", 1+2*(n-1))
	}
	if strings.ToUpper(fields[0]) != "YES" {
		return fmt.Errorf("expected YES")
	}
	edges := make([][2]int, n-1)
	idx := 1
	for i := 0; i < n-1; i++ {
		u, err1 := strconv.Atoi(fields[idx])
		v, err2 := strconv.Atoi(fields[idx+1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("bad integers")
		}
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("bad node")
		}
		edges[i] = [2]int{u, v}
		idx += 2
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	// check tree property
	visited := make([]bool, n+1)
	q := []int{1}
	visited[1] = true
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range adj[v] {
			if !visited[to] {
				visited[to] = true
				q = append(q, to)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return fmt.Errorf("graph not connected")
		}
	}
	// check distances
	d1 := bfs(n, adj, 1)
	d2 := bfs(n, adj, 2)
	if d1[2] != d12 || d2[3] != d23 || d1[3] != d13 {
		return fmt.Errorf("wrong distances")
	}
	return nil
}

func genCase(rng *rand.Rand) (string, int, int, int, int, bool) {
	n := rng.Intn(7) + 3
	d12 := rng.Intn(n-1) + 1
	d23 := rng.Intn(n-1) + 1
	d13 := rng.Intn(n-1) + 1
	possible, _ := construct(n, d12, d23, d13)
	input := fmt.Sprintf("1\n%d %d %d %d\n", n, d12, d23, d13)
	return input, n, d12, d23, d13, possible
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, d12, d23, d13, possible := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verify(out, n, d12, d23, d13, possible); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
