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

type Edge struct {
	to, idx int
}

var (
	adj [][]Edge
	ans [][]int
	mx  int
)

func dfs(u, last, parent int) {
	tot := 0
	for _, e := range adj[u] {
		v := e.to
		if v == parent {
			continue
		}
		tot++
		if tot == last {
			tot++
		}
		if tot > mx {
			mx = tot
		}
		ans[tot] = append(ans[tot], e.idx)
		dfs(v, tot, u)
	}
}

func solve(n int, edges [][2]int) (int, [][]int) {
	adj = make([][]Edge, n+1)
	ans = make([][]int, n+2)
	mx = 0
	for i, e := range edges {
		u := e[0]
		v := e[1]
		idx := i + 1
		adj[u] = append(adj[u], Edge{v, idx})
		adj[v] = append(adj[v], Edge{u, idx})
	}
	dfs(1, 0, 0)
	res := make([][]int, mx)
	for i := 1; i <= mx; i++ {
		tmp := make([]int, len(ans[i]))
		copy(tmp, ans[i])
		res[i-1] = tmp
	}
	return mx, res
}

func formatOutput(k int, days [][]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", k))
	for _, d := range days {
		sb.WriteString(fmt.Sprintf("%d", len(d)))
		for _, id := range d {
			sb.WriteString(fmt.Sprintf(" %d", id))
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
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
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			x, _ := strconv.Atoi(parts[2+2*i])
			y, _ := strconv.Atoi(parts[2+2*i+1])
			edges[i] = [2]int{x, y}
		}
		k, days := solve(n, edges)
		expect := formatOutput(k, days)

		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("Test %d failed:\nexpected:\n%s\n\ngot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
