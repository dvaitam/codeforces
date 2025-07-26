package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n    int
	adj  [][]int
	path []int
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveC(tc testCase) string {
	n := tc.n
	adj := tc.adj
	m := len(tc.path)
	// compute distances via BFS for each node
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = -1
		}
		queue := []int{i}
		dist[i][i] = 0
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, v := range adj[u] {
				if dist[i][v] == -1 {
					dist[i][v] = dist[i][u] + 1
					queue = append(queue, v)
				}
			}
		}
	}
	ans := []int{tc.path[0]}
	last := 0
	for i := 1; i < m-1; i++ {
		if dist[tc.path[last]][tc.path[i+1]] < i+1-last {
			ans = append(ans, tc.path[i])
			last = i
		}
	}
	ans = append(ans, tc.path[m-1])
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	return sb.String()
}

func generateGraph(rng *rand.Rand, n int) [][]int {
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if rng.Intn(2) == 0 {
				adj[i] = append(adj[i], j)
			}
		}
		if len(adj[i]) == 0 {
			j := (i + 1) % n
			adj[i] = append(adj[i], j)
		}
	}
	return adj
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(3))
	tests := make([]testCase, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(5) + 2
		adj := generateGraph(rng, n)
		m := rng.Intn(10) + 2
		path := make([]int, m)
		path[0] = rng.Intn(n)
		for i := 1; i < m; i++ {
			options := adj[path[i-1]]
			path[i] = options[rng.Intn(len(options))]
		}
		tests = append(tests, testCase{n: n, adj: adj, path: path})
	}
	return tests
}

func formatInput(tc testCase) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		row := make([]byte, tc.n)
		for j := 0; j < tc.n; j++ {
			row[j] = '0'
		}
		for _, v := range tc.adj[i] {
			row[v] = '1'
		}
		sb.WriteString(fmt.Sprintf("%s\n", string(row)))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.path)))
	for i, v := range tc.path {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := formatInput(t)
		expect := solveC(t)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
