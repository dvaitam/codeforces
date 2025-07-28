package main

import (
	"bytes"
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTree(n int) []edge {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	return edges
}

func key(a, b int) string {
	if a > b {
		a, b = b, a
	}
	return fmt.Sprintf("%d-%d", a, b)
}

func checkCase(n int, edges []edge, out string) error {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil || k < 0 {
		return fmt.Errorf("invalid k")
	}
	if len(tokens) != 1+4*k {
		return fmt.Errorf("expected %d operations", k)
	}
	idx := 1
	g := make(map[string]bool)
	deg := make([]int, n+1)
	for _, e := range edges {
		g[key(e.u, e.v)] = true
		deg[e.u]++
		deg[e.v]++
	}
	for op := 0; op < k; op++ {
		x1, _ := strconv.Atoi(tokens[idx])
		idx++
		y1, _ := strconv.Atoi(tokens[idx])
		idx++
		x2, _ := strconv.Atoi(tokens[idx])
		idx++
		y2, _ := strconv.Atoi(tokens[idx])
		idx++
		if x1 < 1 || x1 > n || y1 < 1 || y1 > n || x2 < 1 || x2 > n || y2 < 1 || y2 > n || x1 == y1 || x2 == y2 {
			return fmt.Errorf("bad indices in op %d", op+1)
		}
		k1 := key(x1, y1)
		if !g[k1] {
			return fmt.Errorf("edge to remove not present in op %d", op+1)
		}
		// remove edge
		delete(g, k1)
		deg[x1]--
		deg[y1]--
		k2 := key(x2, y2)
		if g[k2] {
			return fmt.Errorf("edge already exists in op %d", op+1)
		}
		g[k2] = true
		deg[x2]++
		deg[y2]++
	}
	if len(g) != n-1 {
		return fmt.Errorf("final edge count not n-1")
	}
	// check connected via BFS
	visited := make([]bool, n+1)
	queue := []int{1}
	visited[1] = true
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for i := 1; i <= n; i++ {
			if i == x {
				continue
			}
			if g[key(i, x)] && !visited[i] {
				visited[i] = true
				queue = append(queue, i)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return fmt.Errorf("graph not connected")
		}
	}
	oneCnt := 0
	for i := 1; i <= n; i++ {
		if deg[i] > 2 {
			return fmt.Errorf("degree>2")
		}
		if deg[i] == 1 {
			oneCnt++
		}
	}
	if !(oneCnt == 2 || (n == 1 && oneCnt == 0)) {
		return fmt.Errorf("not a bamboo")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		tr := generateTree(n)
		input := fmt.Sprintf("1\n%d\n", n)
		for _, e := range tr {
			input += fmt.Sprintf("%d %d\n", e.u, e.v)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", t, err)
			os.Exit(1)
		}
		if err := checkCase(n, tr, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%soutput:%s\n", t, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
