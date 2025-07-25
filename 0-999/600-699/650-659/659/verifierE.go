package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	fields := strings.Fields(input)
	idx := 0
	n := 0
	m := 0
	fmt.Sscan(fields[idx], &n)
	idx++
	fmt.Sscan(fields[idx], &m)
	idx++
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Sscan(fields[idx], &x)
		idx++
		fmt.Sscan(fields[idx], &y)
		idx++
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	vis := make([]bool, n+1)
	ans := 0
	stack := make([]int, 0)
	for i := 1; i <= n; i++ {
		if !vis[i] {
			countV := 0
			countE := 0
			stack = append(stack[:0], i)
			vis[i] = true
			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				countV++
				countE += len(adj[v])
				for _, to := range adj[v] {
					if !vis[to] {
						vis[to] = true
						stack = append(stack, to)
					}
				}
			}
			if countE/2 == countV-1 {
				ans++
			}
		}
	}
	return fmt.Sprintln(ans)
}

func generateTests() []string {
	rand.Seed(46)
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(8) + 2
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges + 1)
		edges := make([][2]int, 0, m)
		exist := make(map[[2]int]bool)
		for len(edges) < m {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			if x == y {
				continue
			}
			if x > y {
				x, y = y, x
			}
			p := [2]int{x, y}
			if !exist[p] {
				exist[p] = true
				edges = append(edges, p)
			}
		}
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
		tests[t] = b.String()
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expect := strings.TrimSpace(solve(t))
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if expect != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
