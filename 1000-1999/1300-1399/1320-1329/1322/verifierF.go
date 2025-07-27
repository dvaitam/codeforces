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

type testCase struct {
	n     int
	m     int
	edges [][2]int
	lines [][2]int
}

func pathBetween(n int, edges [][2]int, a, b int) []int {
	// adjacency
	g := make([][]int, n+1)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	// BFS to get parent
	q := []int{a}
	parent := make([]int, n+1)
	parent[a] = -1
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if v == b {
			break
		}
		for _, to := range g[v] {
			if parent[to] == 0 && to != a {
				parent[to] = v
				q = append(q, to)
			}
		}
	}
	// reconstruct path
	path := []int{}
	cur := b
	for cur != -1 && cur != 0 {
		path = append([]int{cur}, path...)
		if cur == a {
			break
		}
		cur = parent[cur]
	}
	return path
}

func isMono(vals []int, path []int) bool {
	inc := true
	for i := 1; i < len(path); i++ {
		if vals[path[i]-1] <= vals[path[i-1]-1] {
			inc = false
			break
		}
	}
	if inc {
		return true
	}
	dec := true
	for i := 1; i < len(path); i++ {
		if vals[path[i]-1] >= vals[path[i-1]-1] {
			dec = false
			break
		}
	}
	return dec
}

func checkAssign(tc testCase, vals []int) bool {
	for _, line := range tc.lines {
		p := pathBetween(tc.n, tc.edges, line[0], line[1])
		if len(p) == 0 {
			return false
		}
		if !isMono(vals, p) {
			return false
		}
	}
	return true
}

func dfsAssign(tc testCase, idx int, k int, vals []int) bool {
	if idx == tc.n {
		return checkAssign(tc, vals)
	}
	for v := 1; v <= k; v++ {
		vals[idx] = v
		if dfsAssign(tc, idx+1, k, vals) {
			return true
		}
	}
	return false
}

func solve(tc testCase) (int, []int, bool) {
	for k := 1; k <= tc.n; k++ {
		vals := make([]int, tc.n)
		if dfsAssign(tc, 0, k, vals) {
			return k, vals, true
		}
	}
	return 0, nil, false
}

func run(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for _, l := range tc.lines {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[0], l[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		pos := 2
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u, _ := strconv.Atoi(parts[pos])
			v, _ := strconv.Atoi(parts[pos+1])
			edges[i] = [2]int{u, v}
			pos += 2
		}
		lines := make([][2]int, m)
		for i := 0; i < m; i++ {
			a, _ := strconv.Atoi(parts[pos])
			b, _ := strconv.Atoi(parts[pos+1])
			lines[i] = [2]int{a, b}
			pos += 2
		}
		tc := testCase{n: n, m: m, edges: edges, lines: lines}
		k, _, ok := solve(tc)
		if !ok {
			continue
		}
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		fields := strings.Fields(gotStr)
		if len(fields) != n+1 {
			fmt.Fprintf(os.Stderr, "case %d bad output\n", idx)
			os.Exit(1)
		}
		gk, _ := strconv.Atoi(fields[0])
		if gk != k {
			fmt.Fprintf(os.Stderr, "case %d failed: expected k=%d got %d\n", idx, k, gk)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
