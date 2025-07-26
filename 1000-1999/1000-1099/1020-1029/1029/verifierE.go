package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	Input  string
	Output string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Solver struct {
	n   int
	adj [][]int
	f   [][]int
}

const INF = 0x7f7f7f7f

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min3(a, b, c int) int { return min(min(a, b), c) }

func (s *Solver) dp(u, fa int) {
	s.f[u][0], s.f[u][1], s.f[u][2] = 0, 0, 0
	mn := INF
	for _, v := range s.adj[u] {
		if v == fa {
			continue
		}
		s.dp(v, u)
		m012 := min3(s.f[v][0], s.f[v][1], s.f[v][2])
		s.f[u][0] += m012
		m02 := min(s.f[v][0], s.f[v][2])
		s.f[u][1] += m02
		s.f[u][2] += m02
		diff := s.f[v][0] - min(s.f[v][0], s.f[v][2])
		if diff < mn {
			mn = diff
		}
	}
	if fa != 1 {
		s.f[u][0]++
	}
	s.f[u][2] += mn
}

func solveE(n int, edges [][2]int) string {
	s := Solver{n: n, adj: make([][]int, n+1), f: make([][]int, n+1)}
	for i := 1; i <= n; i++ {
		s.f[i] = make([]int, 3)
	}
	for _, e := range edges {
		u, v := e[0], e[1]
		s.adj[u] = append(s.adj[u], v)
		s.adj[v] = append(s.adj[v], u)
	}
	ans := 0
	for _, c := range s.adj[1] {
		s.dp(c, 1)
		for _, v := range s.adj[c] {
			ans += min3(s.f[v][0], s.f[v][1], s.f[v][2])
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func generateTests() []TestCase {
	rand.Seed(46)
	tests := make([]TestCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(15) + 2
		edges := generateTree(n)
		inputBuilder := strings.Builder{}
		inputBuilder.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			inputBuilder.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		output := solveE(n, edges)
		tests[t] = TestCase{Input: inputBuilder.String(), Output: output}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.Input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.Output) {
			fmt.Fprintf(os.Stderr, "Test %d failed:\ninput:\n%s\nexpected:%s\n got:%s\n", i+1, tc.Input, tc.Output, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
