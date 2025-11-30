package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const rawTestcasesData = `
2 1 1 2
4 2 2 3 1 4
6 3 2 4 3 6 1 6
6 7 1 2 1 3 1 4 1 5 1 6 2 4 5 6
6 4 1 2 1 5 4 5 1 6
5 2 1 2 1 5
5 4 1 2 3 4 3 5 1 5
6 6 1 4 1 5 2 4 2 6 4 6 5 6
3 1 1 3
2 1 1 2
4 3 1 2 1 4 2 4
2 1 1 2
3 2 1 3 2 3
2 1 1 2
2 1 1 2
3 3 1 2 1 3 2 3
6 6 1 4 1 6 2 5 2 6 3 6 4 6
2 1 1 2
3 1 1 3
4 3 1 2 1 4 3 4
6 9 1 4 1 5 1 6 2 4 3 4 3 5 4 5 4 6 5 6
6 9 1 3 1 6 2 3 2 5 3 4 3 5 3 6 4 5 4 6
5 4 1 4 2 4 3 5 1 5
5 4 2 4 3 4 3 5 1 5
4 3 1 2 2 3 2 4
2 1 1 2
2 1 1 2
5 5 1 5 2 3 2 4 3 4 4 5
4 4 1 2 2 3 2 4 3 4
2 1 1 2
4 2 2 3 1 4
5 3 1 2 1 4 1 5
4 3 1 2 1 3 1 4
2 1 1 2
6 8 1 2 1 3 1 5 2 3 2 4 3 5 3 6 4 5
5 6 1 2 1 3 1 4 2 5 3 4 4 5
3 1 1 3
4 2 1 4 2 4
3 3 1 2 1 3 2 3
4 2 3 4 1 4
2 1 1 2
3 1 1 3
5 3 1 4 3 5 1 5
6 5 1 5 2 4 3 4 3 6 5 6
6 3 1 2 1 6 2 6
2 1 1 2
4 2 1 2 1 4
3 2 1 2 1 3
4 4 1 2 1 3 2 3 2 4
6 4 1 4 2 5 3 4 4 6
2 1 1 2
3 2 1 2 1 3
6 10 1 3 1 4 1 6 2 3 2 5 3 4 3 5 3 6 4 5 5 6
5 4 1 4 2 4 2 5 3 4
4 5 1 2 1 3 1 4 2 3 3 4
4 3 1 2 3 4 1 4
5 4 2 4 3 5 4 5 1 5
3 2 1 2 1 3
2 1 1 2
3 2 1 2 1 3
2 1 1 2
4 3 1 2 2 4 3 4
6 12 1 2 1 3 1 5 1 6 2 4 2 5 2 6 3 4 3 6 4 5 4 6 5 6
6 6 1 2 2 3 2 4 2 6 3 4 3 6
3 1 1 3
6 4 1 5 3 4 3 5 1 6
2 1 1 2
6 9 1 3 1 4 1 6 2 3 2 4 2 5 2 6 3 5 4 6
3 3 1 2 1 3 2 3
6 7 1 3 1 5 2 3 2 5 2 6 3 4 3 6
5 3 1 2 1 4 4 5
3 3 1 2 1 3 2 3
3 1 1 3
2 1 1 2
5 4 1 2 1 5 2 4 3 4
3 2 1 2 1 3
2 1 1 2
5 4 1 5 2 3 2 4 3 5
3 1 1 3
5 3 1 4 1 5 3 4
2 1 1 2
2 1 1 2
4 3 1 3 2 3 1 4
3 2 1 3 2 3
5 2 1 4 4 5
2 1 1 2
6 5 2 3 3 5 4 5 5 6 1 6
6 7 1 2 1 4 2 4 2 6 3 4 3 6 4 5
2 1 1 2
3 2 1 2 1 3
3 2 1 2 1 3
6 6 1 3 2 3 2 4 3 5 4 5 5 6
3 2 1 2 1 3
2 1 1 2
6 5 1 3 2 3 4 5 4 6 1 6
3 2 1 3 2 3
2 1 1 2
4 3 1 3 2 3 1 4
3 2 1 2 2 3
2 1 1 2
`

func loadTestcases() []string {
	lines := strings.Split(strings.TrimSpace(rawTestcasesData), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

func bfs(n int, G [][]int, start int) ([]int, []int64) {
	inf := int(1e9)
	dist := make([]int, n+1)
	memo := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf
	}
	queue := make([]int, 0, n)
	head := 0
	dist[start] = 0
	memo[start] = 1
	queue = append(queue, start)
	for head < len(queue) {
		u := queue[head]
		head++
		for _, v := range G[u] {
			if dist[v] == inf {
				dist[v] = dist[u] + 1
				queue = append(queue, v)
				memo[v] += memo[u]
			} else if dist[v] == dist[u]+1 {
				memo[v] += memo[u]
			}
		}
	}
	return dist, memo
}

func solveInstance(n int, edges [][2]int) string {
	G := make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		G[x] = append(G[x], y)
		G[y] = append(G[y], x)
	}
	dist1, memo1 := bfs(n, G, 1)
	dist2, memo2 := bfs(n, G, n)
	total := memo1[n]
	best := 0.0
	shortest := dist1[n]
	for i := 1; i <= n; i++ {
		if dist1[i]+dist2[i] != shortest {
			continue
		}
		safe := memo1[i] * memo2[i]
		cur := float64(safe) / float64(total)
		if i != 1 && i != n {
			cur *= 2.0
		}
		if cur > best {
			best = cur
		}
	}
	return fmt.Sprintf("%.9f", best)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range loadTestcases() {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 {
			fmt.Printf("test %d: invalid line\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+2*m {
			fmt.Printf("test %d: expected %d edge pairs got %d\n", idx+1, m, (len(fields)-2)/2)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			x, _ := strconv.Atoi(fields[2+2*i])
			y, _ := strconv.Atoi(fields[3+2*i])
			edges[i] = [2]int{x, y}
		}
		expect := solveInstance(n, edges)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
