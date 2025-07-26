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
	input  string
	output string
}

var (
	n, m     int
	edges    [][]int
	check    []bool
	c        []int
	pre      []int
	p        []int
	tot      int
	tag      []bool
	pos      []int
	flagArr  []bool
	sumArr   []int
	resCount int
)

func dfs1(x int) {
	if tot > 0 || check[x] {
		return
	}
	c[x] = -1
	for _, v := range edges[x] {
		if check[v] {
			continue
		}
		if c[v] == -1 {
			y := x
			for y != v {
				p = append(p, y)
				tot++
				y = pre[y]
			}
			p = append(p, v)
			tot++
			return
		}
		if c[v] == 0 {
			pre[v] = x
			dfs1(v)
			if tot > 0 {
				return
			}
		}
	}
	c[x] = 1
}

func dfs2(S, x int) {
	if check[x] {
		return
	}
	flagArr[x] = true
	for _, v := range edges[x] {
		if check[v] {
			continue
		}
		if tag[v] {
			if tag[x] && tag[v] {
				continue
			}
			if pos[S] >= pos[v] {
				continue
			}
			resCount++
			sumArr[1]++
			sumArr[pos[S]+1]--
			sumArr[pos[v]]++
			sumArr[tot+1]--
		} else if !flagArr[v] {
			dfs2(S, v)
		}
	}
}

func dfs3(S, x int) {
	if check[x] {
		return
	}
	flagArr[x] = true
	for _, v := range edges[x] {
		if check[v] {
			continue
		}
		if tag[v] {
			if tag[x] && tag[v] {
				continue
			}
			if pos[S] < pos[v] {
				continue
			}
			resCount++
			sumArr[pos[v]]++
			sumArr[pos[S]+1]--
		} else if !flagArr[v] {
			dfs3(S, v)
		}
	}
}

func solveInstance() int {
	c = make([]int, n+1)
	pre = make([]int, n+1)
	p = p[:0]
	tot = 0
	resCount = 0
	for i := 1; i <= n && tot == 0; i++ {
		if c[i] == 0 {
			dfs1(i)
		}
	}
	if tot == 0 {
		return 0
	}
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
	tag = make([]bool, n+1)
	pos = make([]int, n+1)
	sumArr = make([]int, tot+2)
	for i, v := range p {
		idx := i + 1
		tag[v] = true
		pos[v] = idx
	}
	flagArr = make([]bool, n+1)
	for _, v := range p {
		dfs2(v, v)
		for i := range flagArr {
			flagArr[i] = false
		}
	}
	for i := len(p) - 1; i >= 0; i-- {
		v := p[i]
		dfs3(v, v)
		for j := range flagArr {
			flagArr[j] = false
		}
	}
	for i := 1; i <= tot; i++ {
		sumArr[i] += sumArr[i-1]
		if sumArr[i] == resCount {
			po := p[i-1]
			check[po] = true
			return po
		}
	}
	return 0
}

func solve(nv int, edgesInput [][2]int) string {
	n = nv
	m = len(edgesInput)
	edges = make([][]int, n+1)
	check = make([]bool, n+1)
	for _, e := range edgesInput {
		a, b := e[0], e[1]
		edges[a] = append(edges[a], b)
	}
	ans := solveInstance()
	if ans == 0 {
		return "-1"
	}
	if solveInstance() == 0 {
		return fmt.Sprintf("%d", ans)
	}
	return "-1"
}

func generateTests() []testCase {
	rand.Seed(6)
	var tests []testCase
	tests = append(tests, testCase{input: "1 0\n", output: "1"})
	for len(tests) < 120 {
		n := rand.Intn(5) + 1
		maxEdges := n * n
		m := rand.Intn(maxEdges + 1)
		edgeMap := map[[2]int]bool{}
		edges := make([][2]int, 0, m)
		for len(edges) < m {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			if a == b {
				continue
			}
			e := [2]int{a, b}
			if edgeMap[e] {
				continue
			}
			edgeMap[e] = true
			edges = append(edges, e)
		}
		var bld strings.Builder
		bld.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
		for _, e := range edges {
			fmt.Fprintf(&bld, "%d %d\n", e[0], e[1])
		}
		tests = append(tests, testCase{input: bld.String(), output: solve(n, edges)})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.output {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, tc.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
