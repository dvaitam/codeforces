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

type Edge struct{ u, v int }

func buildTree(n int, edges []Edge) (adj [][]int) {
	adj = make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	return
}

func dfs(u, p int, adj [][]int, parent, depth []int) {
	parent[u] = p
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		depth[v] = depth[u] + 1
		dfs(v, u, adj, parent, depth)
	}
}

func lca(a, b int, parent, depth []int) int {
	for depth[a] > depth[b] {
		a = parent[a]
	}
	for depth[b] > depth[a] {
		b = parent[b]
	}
	for a != b {
		a = parent[a]
		b = parent[b]
	}
	return a
}

func pathSum(x, y int, parent, depth, val []int) int64 {
	l := lca(x, y, parent, depth)
	idx := 0
	ans := int64(0)
	u := x
	for u != l {
		ans += int64(val[u] ^ idx)
		idx++
		u = parent[u]
	}
	ans += int64(val[l] ^ idx)
	idx++
	var stack []int
	v := y
	for v != l {
		stack = append(stack, v)
		v = parent[v]
	}
	for i := len(stack) - 1; i >= 0; i-- {
		ans += int64(val[stack[i]] ^ idx)
		idx++
	}
	return ans
}

func expected(n int, edges []Edge, vals []int, queries [][2]int) string {
	adj := buildTree(n, edges)
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	dfs(1, 1, adj, parent, depth)
	res := make([]string, len(queries))
	for i, q := range queries {
		ans := pathSum(q[0], q[1], parent, depth, vals)
		res[i] = fmt.Sprintf("%d", ans)
	}
	return strings.Join(res, "\n")
}

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesG.txt")
	if err != nil {
		fmt.Println("failed to open testcasesG.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		fmt.Println("invalid test count")
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= t; caseNum++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		edges := make([]Edge, 0, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			edges = append(edges, Edge{u, v})
		}
		vals := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &vals[i])
		}
		var q int
		fmt.Fscan(in, &q)
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &queries[i][0], &queries[i][1])
		}
		expect := expected(n, edges, vals, queries)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(vals[i]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for _, qu := range queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
		}
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed:\nexpected:\n%s\n got:\n%s\n", caseNum, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
