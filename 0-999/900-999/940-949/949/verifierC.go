package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type test struct {
	n, m, h int
	w       []int
	edges   [][2]int
}

func genTests() []test {
	rand.Seed(3)
	tests := make([]test, 100)
	for i := range tests {
		n := rand.Intn(10) + 2
		h := rand.Intn(9) + 2
		m := rand.Intn(n*(n-1)/2 + 1)
		w := make([]int, n)
		for j := 0; j < n; j++ {
			w[j] = rand.Intn(h)
		}
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			for v == u {
				v = rand.Intn(n) + 1
			}
			edges[j] = [2]int{u, v}
		}
		tests[i] = test{n, m, h, w, edges}
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// reference solution from 949C.go
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveRef(t test) []int {
	n, h := t.n, t.h
	w := make([]int, n+1)
	for i := 1; i <= n; i++ {
		w[i] = t.w[i-1]
	}
	graph := make([][]int, n+1)
	for _, e := range t.edges {
		u, v := e[0], e[1]
		if (w[u]+1)%h == w[v] {
			graph[u] = append(graph[u], v)
		}
		if (w[v]+1)%h == w[u] {
			graph[v] = append(graph[v], u)
		}
	}
	dfn := make([]int, n+1)
	low := make([]int, n+1)
	onStack := make([]bool, n+1)
	stack := make([]int, 0, n)
	var index, sccCount int
	sccId := make([]int, n+1)
	size := make([]int, n+1)
	var dfs func(u int)
	dfs = func(u int) {
		index++
		dfn[u] = index
		low[u] = index
		stack = append(stack, u)
		onStack[u] = true
		for _, v := range graph[u] {
			if dfn[v] == 0 {
				dfs(v)
				low[u] = min(low[u], low[v])
			} else if onStack[v] {
				low[u] = min(low[u], dfn[v])
			}
		}
		if low[u] == dfn[u] {
			sccCount++
			for {
				sz := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[sz] = false
				sccId[sz] = sccCount
				size[sccCount]++
				if sz == u {
					break
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if dfn[i] == 0 {
			dfs(i)
		}
	}
	outDeg := make([]int, sccCount+1)
	for u := 1; u <= n; u++ {
		for _, v := range graph[u] {
			if sccId[u] != sccId[v] {
				outDeg[sccId[u]]++
			}
		}
	}
	bestSize := n + 1
	bestId := 0
	for i := 1; i <= sccCount; i++ {
		if outDeg[i] == 0 && size[i] < bestSize {
			bestSize = size[i]
			bestId = i
		}
	}
	res := []int{}
	for i := 1; i <= n; i++ {
		if sccId[i] == bestId {
			res = append(res, i)
		}
	}
	sort.Ints(res)
	return res
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", t.n, t.m, t.h))
	for i, v := range t.w {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func parseOutput(out string) ([]int, bool) {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return nil, false
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil || k < 0 {
		return nil, false
	}
	vals := make([]int, 0, k)
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, false
		}
		vals = append(vals, v)
	}
	if len(vals) != k {
		return nil, false
	}
	sort.Ints(vals)
	return vals, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, ok := parseOutput(out)
		if !ok {
			fmt.Printf("test %d: bad output\n", i+1)
			os.Exit(1)
		}
		exp := solveRef(t)
		if len(got) != len(exp) {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%v\noutput:%v\n", i+1, input, exp, got)
			os.Exit(1)
		}
		match := true
		for j := range got {
			if got[j] != exp[j] {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%v\noutput:%v\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
