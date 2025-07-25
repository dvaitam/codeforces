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

func runCandidate(bin, input string) (string, error) {
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

func expected(n int, colors []int, edges [][2]int) string {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n)
	comp := [2]int{}
	stack := make([]int, 0, n)
	for target := 0; target <= 1; target++ {
		for i := 0; i < n; i++ {
			visited[i] = false
		}
		for i := 0; i < n; i++ {
			if !visited[i] && colors[i] == target {
				comp[target]++
				stack = stack[:0]
				stack = append(stack, i)
				visited[i] = true
				for len(stack) > 0 {
					u := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					for _, v := range adj[u] {
						if !visited[v] && colors[v] == target {
							visited[v] = true
							stack = append(stack, v)
						}
					}
				}
			}
		}
	}
	ans := comp[0]
	if comp[1] < ans {
		ans = comp[1]
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	colors := make([]int, n)
	for i := range colors {
		colors[i] = rng.Intn(2)
	}
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges[i-1] = [2]int{i, p}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(c))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	input := sb.String()
	exp := expected(n, colors, edges)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
