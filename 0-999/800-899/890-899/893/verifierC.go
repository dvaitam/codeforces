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

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(n, m int, cost []int64, edges [][2]int) int64 {
	g := make([][]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	visited := make([]bool, n)
	q := make([]int, 0)
	var ans int64
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		visited[i] = true
		q = append(q[:0], i)
		minCost := cost[i]
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			if cost[v] < minCost {
				minCost = cost[v]
			}
			for _, to := range g[v] {
				if !visited[to] {
					visited[to] = true
					q = append(q, to)
				}
			}
		}
		ans += minCost
	}
	return ans
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(20) + 1
	maxEdges := n * (n - 1) / 2
	m := r.Intn(maxEdges + 1)
	cost := make([]int64, n)
	for i := range cost {
		cost[i] = int64(r.Intn(100) + 1)
	}
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		a := r.Intn(n)
		b := r.Intn(n)
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, [2]int{a, b})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(cost[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	expect := fmt.Sprint(solveCase(n, m, cost, edges))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
