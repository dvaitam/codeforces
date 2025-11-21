package main

import (
	"bufio"
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(tc.expect, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected: %s\nactual: %s\n", i+1, err, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(expect, output string) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(output, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not an integer: %v", err)
	}
	exp, _ := strconv.ParseInt(expect, 10, 64)
	if val != exp {
		return fmt.Errorf("expected %s but got %s", expect, output)
	}
	return nil
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

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest(3, [][2]int{{1, 2}, {2, 3}, {1, 3}}),
		makeTest(4, [][2]int{{1, 2}, {2, 3}, {3, 4}}),
		makeTest(5, [][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}, {4, 5}}),
	}
	for t := 0; t < 200; t++ {
		n := rand.Intn(8) + 3
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges + 1)
		edges := randomEdges(n, m)
		tests = append(tests, makeTest(n, edges))
	}
	return tests
}

func randomEdges(n, m int) [][2]int {
	used := make(map[int]struct{})
	res := make([][2]int, 0, m)
	for len(res) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := u*(n+1) + v
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		res = append(res, [2]int{u, v})
	}
	return res
}

func makeTest(n int, edges [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	input := sb.String()
	return testCase{input: input, expect: solveRef(input)}
}

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(reader, &n, &m)
	adj := make([][]bool, n+1)
	for i := range adj {
		adj[i] = make([]bool, n+1)
	}
	deg := make([]int, n+1)
	neigh := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if !adj[a][b] {
			adj[a][b] = true
			adj[b][a] = true
			deg[a]++
			deg[b]++
			neigh[a] = append(neigh[a], b)
			neigh[b] = append(neigh[b], a)
		}
	}
	const inf = int(1e9)
	ans := inf
	for u := 1; u <= n; u++ {
		l := len(neigh[u])
		for i := 0; i < l; i++ {
			v := neigh[u][i]
			for j := i + 1; j < l; j++ {
				w := neigh[u][j]
				if adj[v][w] {
					sum := deg[u] + deg[v] + deg[w] - 6
					if sum < ans {
						ans = sum
					}
				}
			}
		}
	}
	if ans == inf {
		return "-1"
	}
	return fmt.Sprintf("%d", ans)
}
