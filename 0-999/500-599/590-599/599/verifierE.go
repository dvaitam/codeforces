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

type edge struct{ u, v int }

type query struct{ a, b, c int }

func lca(parent []int, a, b int) int {
	seen := make(map[int]bool)
	for a != 0 {
		seen[a] = true
		a = parent[a]
	}
	for {
		if seen[b] {
			return b
		}
		b = parent[b]
	}
}

func countTrees(n int, edges []edge, qs []query) int {
	parent := make([]int, n+1)
	var dfs func(int)
	cnt := 0
	dfs = func(i int) {
		if i > n {
			// check edges
			for _, e := range edges {
				if !(parent[e.u] == e.v || parent[e.v] == e.u) {
					return
				}
			}
			// connectivity from root 1
			vis := make([]bool, n+1)
			stack := []int{1}
			vis[1] = true
			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for j := 2; j <= n; j++ {
					if parent[j] == v && !vis[j] {
						vis[j] = true
						stack = append(stack, j)
					}
				}
			}
			for j := 1; j <= n; j++ {
				if !vis[j] {
					return
				}
			}
			// check queries
			for _, qq := range qs {
				if lca(parent, qq.a, qq.b) != qq.c {
					return
				}
			}
			cnt++
			return
		}
		for p := 1; p <= n; p++ {
			if p == i {
				continue
			}
			parent[i] = p
			dfs(i + 1)
		}
	}
	parent[1] = 0
	dfs(2)
	return cnt
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		q, _ := strconv.Atoi(scan.Text())
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			edges[i].u, _ = strconv.Atoi(scan.Text())
			scan.Scan()
			edges[i].v, _ = strconv.Atoi(scan.Text())
		}
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			scan.Scan()
			qs[i].a, _ = strconv.Atoi(scan.Text())
			scan.Scan()
			qs[i].b, _ = strconv.Atoi(scan.Text())
			scan.Scan()
			qs[i].c, _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		for _, qq := range qs {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", qq.a, qq.b, qq.c))
		}
		input := sb.String()
		cnt := countTrees(n, edges, qs)
		exp := fmt.Sprintf("%d\n", cnt)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
