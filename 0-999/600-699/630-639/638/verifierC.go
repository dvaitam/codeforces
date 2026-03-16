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
	const testcasesCRaw = `4 3 1 2 2 3 1 4
9 8 1 2 2 3 3 4 4 5 2 6 1 7 4 8 1 9
8 7 1 2 1 3 3 4 4 5 3 6 6 7 7 8
5 4 1 2 2 3 1 4 1 5
2 1 1 2
8 7 1 2 2 3 3 4 1 5 5 6 2 7 7 8
9 8 1 2 1 3 2 4 2 5 2 6 4 7 3 8 1 9
8 7 1 2 1 3 3 4 3 5 1 6 6 7 3 8
10 9 1 2 1 3 2 4 3 5 5 6 4 7 7 8 7 9 1 10
9 8 1 2 2 3 2 4 2 5 3 6 5 7 6 8 6 9
3 2 1 2 1 3
4 3 1 2 2 3 2 4
2 1 1 2
2 1 1 2
8 7 1 2 1 3 3 4 2 5 1 6 2 7 5 8
10 9 1 2 2 3 3 4 3 5 5 6 3 7 4 8 5 9 9 10
2 1 1 2
10 9 1 2 1 3 2 4 1 5 4 6 3 7 5 8 4 9 9 10
8 7 1 2 2 3 2 4 3 5 1 6 5 7 5 8
7 6 1 2 1 3 1 4 2 5 5 6 5 7
4 3 1 2 2 3 1 4
3 2 1 2 1 3
9 8 1 2 2 3 1 4 3 5 1 6 5 7 2 8 6 9
6 5 1 2 1 3 1 4 3 5 5 6
4 3 1 2 2 3 2 4
7 6 1 2 2 3 1 4 1 5 3 6 4 7
7 6 1 2 1 3 2 4 1 5 3 6 6 7
10 9 1 2 2 3 1 4 2 5 1 6 4 7 2 8 1 9 3 10
9 8 1 2 1 3 3 4 4 5 2 6 5 7 6 8 1 9
8 7 1 2 2 3 1 4 3 5 2 6 2 7 1 8
6 5 1 2 1 3 2 4 3 5 2 6
8 7 1 2 1 3 1 4 1 5 5 6 2 7 5 8
9 8 1 2 1 3 2 4 2 5 3 6 1 7 2 8 7 9
5 4 1 2 1 3 3 4 4 5
6 5 1 2 1 3 2 4 4 5 3 6
2 1 1 2
5 4 1 2 1 3 2 4 4 5
5 4 1 2 1 3 2 4 3 5
10 9 1 2 1 3 1 4 1 5 1 6 2 7 2 8 3 9 9 10
5 4 1 2 2 3 3 4 3 5
7 6 1 2 2 3 1 4 3 5 2 6 5 7
9 8 1 2 1 3 2 4 1 5 4 6 1 7 4 8 3 9
4 3 1 2 1 3 3 4
8 7 1 2 1 3 3 4 1 5 3 6 3 7 3 8
10 9 1 2 2 3 2 4 1 5 1 6 3 7 1 8 1 9 2 10
8 7 1 2 1 3 1 4 2 5 5 6 4 7 2 8
3 2 1 2 1 3
5 4 1 2 1 3 2 4 4 5
10 9 1 2 2 3 3 4 4 5 3 6 1 7 2 8 6 9 1 10
2 1 1 2
6 5 1 2 2 3 2 4 3 5 4 6
3 2 1 2 2 3
9 8 1 2 2 3 1 4 4 5 3 6 3 7 2 8 4 9
6 5 1 2 1 3 2 4 1 5 3 6
3 2 1 2 1 3
7 6 1 2 2 3 2 4 1 5 3 6 2 7
7 6 1 2 1 3 2 4 1 5 5 6 5 7
3 2 1 2 1 3
2 1 1 2
8 7 1 2 2 3 3 4 1 5 1 6 1 7 6 8
2 1 1 2
7 6 1 2 2 3 1 4 1 5 5 6 3 7
3 2 1 2 1 3
4 3 1 2 2 3 2 4
3 2 1 2 1 3
5 4 1 2 1 3 2 4 2 5
4 3 1 2 2 3 3 4
4 3 1 2 1 3 2 4
3 2 1 2 2 3
10 9 1 2 2 3 3 4 4 5 1 6 4 7 7 8 6 9 3 10
6 5 1 2 1 3 3 4 4 5 5 6
2 1 1 2
7 6 1 2 1 3 1 4 3 5 3 6 4 7
8 7 1 2 1 3 1 4 4 5 1 6 2 7 5 8
7 6 1 2 1 3 1 4 3 5 4 6 6 7
9 8 1 2 2 3 2 4 3 5 2 6 1 7 1 8 6 9
4 3 1 2 2 3 2 4
6 5 1 2 1 3 3 4 4 5 5 6
3 2 1 2 2 3
4 3 1 2 2 3 2 4
5 4 1 2 2 3 3 4 4 5
7 6 1 2 1 3 3 4 1 5 5 6 1 7
6 5 1 2 2 3 3 4 1 5 2 6
3 2 1 2 1 3
8 7 1 2 2 3 1 4 3 5 4 6 2 7 5 8
9 8 1 2 1 3 2 4 4 5 1 6 6 7 3 8 5 9
5 4 1 2 1 3 1 4 4 5
2 1 1 2
5 4 1 2 1 3 1 4 3 5
4 3 1 2 2 3 2 4
6 5 1 2 1 3 3 4 3 5 4 6
8 7 1 2 1 3 3 4 4 5 2 6 3 7 7 8
3 2 1 2 1 3
2 1 1 2
4 3 1 2 2 3 3 4
6 5 1 2 2 3 3 4 3 5 1 6
3 2 1 2 2 3
7 6 1 2 2 3 2 4 4 5 1 6 6 7
8 7 1 2 1 3 3 4 1 5 3 6 6 7 5 8
10 9 1 2 2 3 3 4 4 5 3 6 6 7 2 8 8 9 9 10`

	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
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
