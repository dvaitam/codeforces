package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD = 1000000007

// enumeration DFS of simple paths without repeating edges
func countPaths(n int, adj [][]int, start, target int) int {
	used := make(map[[2]int]bool)
	var dfs func(u int) int
	dfs = func(u int) int {
		if u == target {
			return 1
		}
		total := 0
		for _, v := range adj[u] {
			e := edge(u, v)
			if used[e] {
				continue
			}
			used[e] = true
			total = (total + dfs(v)) % MOD
			used[e] = false
		}
		return total
	}
	return dfs(start) % MOD
}

func edge(a, b int) [2]int {
	if a < b {
		return [2]int{a, b}
	}
	return [2]int{b, a}
}

func solveCaseE(n int, edges [][2]int, queries [][2]int) []int {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	res := make([]int, len(queries))
	for i, q := range queries {
		res[i] = countPaths(n, adj, q[0]-1, q[1]-1)
	}
	return res
}

type TestCaseE struct {
	n, m    int
	edges   [][2]int
	k       int
	queries [][2]int
	ans     []int
}

func readCasesE(path string) ([]TestCaseE, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]TestCaseE, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			edges[j] = [2]int{u, v}
		}
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		queries := make([][2]int, k)
		for j := 0; j < k; j++ {
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			y, _ := strconv.Atoi(scan.Text())
			queries[j] = [2]int{x, y}
		}
		ans := solveCaseE(n, edges, queries)
		cases[i] = TestCaseE{n, m, edges, k, queries, ans}
	}
	return cases, nil
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := readCasesE("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		fmt.Fprintf(&sb, "%d\n", tc.k)
		for _, q := range tc.queries {
			fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
		}
		expectedLines := make([]string, len(tc.ans))
		for j, v := range tc.ans {
			expectedLines[j] = fmt.Sprintf("%d", v%MOD)
		}
		expected := strings.Join(expectedLines, "\n")
		got, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != expected {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
