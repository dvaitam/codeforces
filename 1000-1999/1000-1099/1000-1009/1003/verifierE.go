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

const testcasesRaw = `100
8 3 1
13 8 2
3 1 1
13 9 3
2 1 5
18 12 3
6 1 3
7 1 3
9 4 2
10 5 3
3 2 4
17 8 2
8 4 3
3 2 1
10 5 5
7 4 4
20 10 4
15 3 2
10 5 1
3 1 4
9 8 3
5 2 1
14 4 4
9 3 3
14 12 5
11 9 2
11 2 1
8 3 5
20 8 1
11 3 3
15 1 1
12 2 3
11 1 3
10 6 2
14 10 1
10 4 4
10 3 3
13 10 2
11 10 1
12 1 4
6 3 3
10 2 4
7 4 2
4 1 1
2 1 5
5 1 5
16 10 2
11 1 1
17 10 4
7 4 2
8 4 4
16 1 2
14 8 2
14 4 4
7 1 1
9 5 2
17 7 2
14 5 2
11 1 3
19 4 5
13 11 1
16 7 1
14 4 5
6 3 3
16 13 3
14 9 2
9 6 4
16 2 3
7 1 4
20 5 3
2 1 4
19 16 4
13 4 1
7 2 1
20 9 1
13 7 2
18 2 2
6 5 3
18 16 5
15 1 1
2 1 4
18 9 5
5 1 3
3 1 3
12 2 1
18 15 4
7 3 4
8 7 4
13 2 1
4 3 3
17 14 4
15 2 2
10 8 4
4 3 2
12 3 2
5 3 4
11 5 5
1 1 1
10 2 5
4 2 5
`

type testCase struct {
	n int
	d int
	k int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// referenceSolution embeds the construction algorithm from 1003E.go so the
// verifier does not need any external oracle or file.
func referenceSolution(N, d, k int) (bool, [][2]int) {
	if d >= N || (k == 1 && N > 2) {
		return false, nil
	}
	if N == 1 && d == 1 {
		return true, [][2]int{}
	}
	count := make([]int, N+2)
	depth := make([]int, N+2)
	edges := make([][2]int, 0, N-1)
	for i := 1; i <= d; i++ {
		u := i
		v := i + 1
		edges = append(edges, [2]int{u, v})
		count[u]++
		count[v]++
		depth[u] = min(u-1, d-u+1)
		depth[v] = min(u, d-u)
	}
	i := d + 2
	j := 2
	for i <= N {
		for j < i && (count[j] == k || depth[j] == 0) {
			j++
		}
		if i == j {
			return false, nil
		}
		edges = append(edges, [2]int{i, j})
		count[i]++
		count[j]++
		depth[i] = depth[j] - 1
		i++
	}
	if len(edges) != N-1 {
		return false, nil
	}
	return true, edges
}

func diameter(n int, edges [][2]int) int {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	bfs := func(start int) (int, int) {
		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = -1
		}
		q := []int{start}
		dist[start] = 0
		idx := 0
		for idx < len(q) {
			v := q[idx]
			idx++
			for _, to := range g[v] {
				if dist[to] == -1 {
					dist[to] = dist[v] + 1
					q = append(q, to)
				}
			}
		}
		far := start
		for i, dv := range dist {
			if dv > dist[far] {
				far = i
			}
		}
		return far, dist[far]
	}
	v, _ := bfs(1)
	w, dist := bfs(v)
	_ = w
	return dist
}

func validateTree(n, d, k int, edges [][2]int) error {
	if len(edges) != n-1 {
		return fmt.Errorf("expected %d edges got %d", n-1, len(edges))
	}
	deg := make([]int, n+1)
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
	components := n
	for _, e := range edges {
		u, v := e[0], e[1]
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("invalid vertex")
		}
		deg[u]++
		deg[v]++
		if deg[u] > k || deg[v] > k {
			return fmt.Errorf("degree limit exceeded")
		}
		ru := find(u)
		rv := find(v)
		if ru == rv {
			return fmt.Errorf("cycle detected")
		}
		parent[ru] = rv
		components--
	}
	if components != 1 {
		return fmt.Errorf("graph not connected")
	}
	if diameter(n, edges) != d {
		return fmt.Errorf("diameter mismatch")
	}
	return nil
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d %d\n", tc.n, tc.d, tc.k)
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	possible, _ := referenceSolution(tc.n, tc.d, tc.k)
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := scanner.Text()
	switch first {
	case "NO":
		if possible {
			return fmt.Errorf("expected YES but got NO")
		}
		return nil
	case "YES":
		if !possible {
			return fmt.Errorf("expected NO but got YES")
		}
	default:
		return fmt.Errorf("invalid output %q", first)
	}
	edges := make([][2]int, 0, tc.n-1)
	for scanner.Scan() {
		uStr := scanner.Text()
		if !scanner.Scan() {
			return fmt.Errorf("incomplete edge")
		}
		vStr := scanner.Text()
		u, err := strconv.Atoi(uStr)
		if err != nil {
			return fmt.Errorf("invalid vertex %q", uStr)
		}
		v, err := strconv.Atoi(vStr)
		if err != nil {
			return fmt.Errorf("invalid vertex %q", vStr)
		}
		edges = append(edges, [2]int{u, v})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return validateTree(tc.n, tc.d, tc.k, edges)
}

func parseTestcases(raw string) ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(raw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("invalid test file")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		vals := make([]int, 3)
		for j := 0; j < 3; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("invalid test file")
			}
			vals[j], err = strconv.Atoi(scan.Text())
			if err != nil {
				return nil, err
			}
		}
		tests = append(tests, testCase{n: vals[0], d: vals[1], k: vals[2]})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
