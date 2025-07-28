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

func lcaInit(n int, edges [][2]int) ([][]int, []int, int) {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	log := 0
	for (1 << log) <= n {
		log++
	}
	parent := make([][]int, log)
	for i := range parent {
		parent[i] = make([]int, n+1)
	}
	depth := make([]int, n+1)
	stack := []int{1}
	parent[0][1] = 0
	depth[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range g[v] {
			if to == parent[0][v] {
				continue
			}
			parent[0][to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	for k := 1; k < log; k++ {
		for i := 1; i <= n; i++ {
			parent[k][i] = parent[k-1][parent[k-1][i]]
		}
	}
	return parent, depth, log
}

func lcaQuery(a, b int, parent [][]int, depth []int, log int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := 0; diff > 0; k++ {
		if diff&1 == 1 {
			a = parent[k][a]
		}
		diff >>= 1
	}
	if a == b {
		return a
	}
	for k := log - 1; k >= 0; k-- {
		if parent[k][a] != parent[k][b] {
			a = parent[k][a]
			b = parent[k][b]
		}
	}
	return parent[0][a]
}

func solveH(n int, edges [][2]int, queries []struct {
	arr  []int
	u, v int
}) []int {
	parent, depth, log := lcaInit(n, edges)
	res := make([]int, len(queries))
	visited := make([]bool, n+1)
	used := make([]int, 0, n)
	for idx, q := range queries {
		l := lcaQuery(q.u, q.v, parent, depth, log)
		addVal := func(x int) {
			val := q.arr[x]
			if !visited[val] {
				visited[val] = true
				used = append(used, val)
			}
		}
		for x := q.u; x != l; x = parent[0][x] {
			addVal(x)
		}
		for x := q.v; x != l; x = parent[0][x] {
			addVal(x)
		}
		addVal(l)
		mex := 0
		for mex <= n && visited[mex] {
			mex++
		}
		res[idx] = mex
		for _, val := range used {
			visited[val] = false
		}
		used = used[:0]
	}
	return res
}

func generateCaseH(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	q := rng.Intn(3) + 1
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	queries := make([]struct {
		arr  []int
		u, v int
	}, q)
	for i := 0; i < q; i++ {
		arr := make([]int, n+1)
		for j := 1; j <= n; j++ {
			arr[j] = rng.Intn(n + 1)
		}
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		queries[i] = struct {
			arr  []int
			u, v int
		}{arr: arr, u: u, v: v}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for _, qu := range queries {
		for j := 1; j <= n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", qu.arr[j]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d %d\n", qu.u, qu.v))
	}
	input := sb.String()
	ans := solveH(n, edges, queries)
	var exp strings.Builder
	for i, v := range ans {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	exp.WriteByte('\n')
	return input, exp.String()
}

func runProg(path, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseH(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
