package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp := strings.ToUpper(strings.TrimSpace(expect))
	act := strings.ToUpper(strings.Split(actual, "\n")[0])
	if act != exp {
		return fmt.Errorf("expected %s but got %s", exp, act)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTestCase(3, [][]int{{1, 2}, {2, 3}}, []int{2}, []int{3}),
		makeTestCase(2, [][]int{{1, 2}}, []int{2}, []int{}),
		makeTestCase(1, [][]int{}, []int{1}, []int{}),
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTestCase())
	}
	return tests
}

func makeTestCase(n int, edges [][]int, tokens []int, bonus []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, len(edges))
	fmt.Fprintf(&sb, "%d %d\n", len(tokens), len(bonus))
	for i, v := range tokens {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range bonus {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	expect := "NO"
	if solveReference(n, edges, tokens, bonus) {
		expect = "YES"
	}
	return testCase{
		input:  sb.String(),
		expect: expect,
	}
}

func randomTestCase() testCase {
	n := rand.Intn(6) + 1
	m := rand.Intn(n*(n-1)/2 + 1)
	edges := make([][]int, 0, m)
	used := map[[2]int]bool{}
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{min(u, v), max(u, v)}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, []int{u, v})
	}
	ensureConnected(n, &edges)
	tokens := randomSet(n, rand.Intn(n)+1)
	bonus := randomSet(n, rand.Intn(n+1))
	return makeTestCase(n, edges, tokens, bonus)
}

func ensureConnected(n int, edges *[][]int) {
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		pa, pb := find(a), find(b)
		if pa != pb {
			parent[pa] = pb
			*edges = append(*edges, []int{a, b})
		}
	}
	for i := 2; i <= n; i++ {
		if find(i) != find(1) && rand.Intn(2) == 0 {
			union(1, i)
		}
	}
	for i := 2; i <= n; i++ {
		if find(i) != find(1) {
			union(i, rand.Intn(i-1)+1)
		}
	}
}

func randomSet(n, size int) []int {
	if size > n {
		size = n
	}
	perm := rand.Perm(n)
	set := make([]int, size)
	for i := 0; i < size; i++ {
		set[i] = perm[i] + 1
	}
	return set
}

func solveReference(n int, edges [][]int, tokens []int, bonus []int) bool {
	graph := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}
	bonusSet := make([]bool, n+1)
	for _, b := range bonus {
		bonusSet[b] = true
	}
	for _, t := range tokens {
		if t == 1 {
			return true
		}
	}
	if len(tokens) == 0 {
		return false
	}
	adj1 := make(map[int]bool)
	for _, to := range graph[1] {
		adj1[to] = true
	}
	for _, t := range tokens {
		if adj1[t] {
			return true
		}
	}
	if len(tokens) < 2 {
		return false
	}
	visited := make([]bool, n+1)
	queue := []int{1}
	visited[1] = true
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range graph[v] {
			if !visited[to] && bonusSet[to] {
				visited[to] = true
				queue = append(queue, to)
			}
		}
	}
	for _, t := range tokens {
		for _, to := range graph[t] {
			if visited[to] {
				return true
			}
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
