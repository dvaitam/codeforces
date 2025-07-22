package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type opC struct {
	t int
	x int
	v int64
}

type testCaseC struct {
	n        int
	m        int
	values   []int64
	edges    [][2]int
	ops      []opC
	expected []string
}

func buildTree(n int, edges [][2]int) [][]int {
	g := make([][]int, n+1)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	return g
}

func dfsOrder(g [][]int, n int) (parent []int, depth []int, subtree [][]int) {
	parent = make([]int, n+1)
	depth = make([]int, n+1)
	subtree = make([][]int, n+1)
	stack := []int{1}
	order := []int{}
	parent[1] = 0
	depth[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			stack = append(stack, v)
		}
	}
	// build subtree lists reverse order
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		subtree[u] = append(subtree[u], u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			subtree[u] = append(subtree[u], subtree[v]...)
		}
	}
	return
}

func solveCase(tc testCaseC) []string {
	g := buildTree(tc.n, tc.edges)
	_, depth, subtree := dfsOrder(g, tc.n)
	vals := make([]int64, tc.n+1)
	copy(vals[1:], tc.values)
	res := make([]string, 0)
	for _, op := range tc.ops {
		if op.t == 1 {
			x := op.x
			v := op.v
			for _, node := range subtree[x] {
				if depth[node]%2 == depth[x]%2 {
					vals[node] += v
				} else {
					vals[node] -= v
				}
			}
		} else {
			x := op.x
			res = append(res, strconv.FormatInt(vals[x], 10))
		}
	}
	return res
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(7) + 2
		m := rng.Intn(5) + 1
		values := make([]int64, n)
		for j := 0; j < n; j++ {
			values[j] = int64(rng.Intn(10))
		}
		edges := make([][2]int, n-1)
		for j := 1; j < n; j++ {
			p := rng.Intn(j) + 1
			edges[j-1] = [2]int{p, j + 1}
		}
		ops := make([]opC, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				x := rng.Intn(n) + 1
				v := int64(rng.Intn(5)) - 2
				ops[j] = opC{t: 1, x: x, v: v}
			} else {
				x := rng.Intn(n) + 1
				ops[j] = opC{t: 2, x: x}
			}
		}
		tc := testCaseC{n: n, m: m, values: values, edges: edges, ops: ops}
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(values[j], 10))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for _, op := range ops {
			if op.t == 1 {
				sb.WriteString(fmt.Sprintf("1 %d %d\n", op.x, op.v))
			} else {
				sb.WriteString(fmt.Sprintf("2 %d\n", op.x))
			}
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\n", i+1, len(expected), len(lines))
			os.Exit(1)
		}
		for idx := range lines {
			if strings.TrimSpace(lines[idx]) != expected[idx] {
				fmt.Fprintf(os.Stderr, "case %d failed on line %d: expected %s got %s\n", i+1, idx+1, expected[idx], lines[idx])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
