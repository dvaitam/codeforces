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

func runExe(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseEdges(out string, n int) ([][2]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n-1 {
		return nil, fmt.Errorf("expected %d edges", n-1)
	}
	edges := make([][2]int, 0, n-1)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("bad edge")
		}
		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil || u < 1 || u > n || v < 1 || v > n {
			return nil, fmt.Errorf("invalid edge")
		}
		edges = append(edges, [2]int{u, v})
	}
	return edges, nil
}

func isTree(n int, edges [][2]int) bool {
	if len(edges) != n-1 {
		return false
	}
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	vis := make([]bool, n+1)
	stack := []int{1}
	vis[1] = true
	cnt := 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		cnt++
		for _, to := range g[v] {
			if !vis[to] {
				vis[to] = true
				stack = append(stack, to)
			}
		}
	}
	return cnt == n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 3
		input := fmt.Sprintf("%d\n", n)
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		edges, err := parseEdges(out, n)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
		if !isTree(n, edges) {
			fmt.Printf("case %d failed: output is not a tree\ninput:\n%soutput:\n%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
