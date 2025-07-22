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

func bfs(n int, G [][]int, start int) ([]int, []int64) {
	inf := 1<<31 - 1
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

func expectedAnswer(n int, edges [][2]int) string {
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
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+2*m {
			fmt.Printf("test %d: expected %d edge pairs got %d\n", idx, m, (len(fields)-2)/2)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			x, _ := strconv.Atoi(fields[2+2*i])
			y, _ := strconv.Atoi(fields[3+2*i])
			edges[i] = [2]int{x, y}
		}
		expect := expectedAnswer(n, edges)
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
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
