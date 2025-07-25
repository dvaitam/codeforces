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

type testCase struct {
	input    string
	expected string
}

func solveCase(n, m int, cats []int, edges [][2]int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	type node struct{ v, cnt int }
	stack := []node{{1, cats[1]}}
	visited := make([]bool, n+1)
	visited[1] = true
	ans := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.cnt > m {
			continue
		}
		isLeaf := true
		for _, to := range adj[cur.v] {
			if !visited[to] {
				visited[to] = true
				nextCnt := 0
				if cats[to] == 1 {
					nextCnt = cur.cnt + 1
				}
				stack = append(stack, node{to, nextCnt})
				isLeaf = false
			}
		}
		if isLeaf {
			ans++
		}
	}
	return ans
}

func buildCase(n, m int, cats []int, edges [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", cats[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	ans := solveCase(n, m, cats, edges)
	return testCase{input: sb.String(), expected: fmt.Sprintf("%d\n", ans)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 2
	m := rng.Intn(n) + 1
	cats := make([]int, n+1)
	for i := 1; i <= n; i++ {
		cats[i] = rng.Intn(2)
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return buildCase(n, m, cats, edges)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cats := []int{0, 1, 1}
	edges := [][2]int{{1, 2}, {1, 3}}
	cases := []testCase{
		buildCase(3, 1, cats, edges),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
