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

type edge struct {
	u, v int
}

type Edge struct{ to, cost int }

func expectedAnswerD(n int, edges []edge) (int, []int) {
	adj := make([][]Edge, n+1)
	for _, e := range edges {
		u, v := e.u, e.v
		adj[u] = append(adj[u], Edge{v, 0})
		adj[v] = append(adj[v], Edge{u, 1})
	}
	dp := make([]int, n+1)
	parent := make([]int, n+1)
	costToParent := make([]int, n+1)
	visited := make([]bool, n+1)
	queue := []int{1}
	visited[1] = true
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		for _, e := range adj[u] {
			v, c := e.to, e.cost
			if !visited[v] {
				visited[v] = true
				dp[1] += c
				parent[v] = u
				costToParent[v] = c
				queue = append(queue, v)
			}
		}
	}
	for i := 1; i < len(queue); i++ {
		v := queue[i]
		u := parent[v]
		c := costToParent[v]
		dp[v] = dp[u] - c + (1 - c)
	}
	minCost := dp[1]
	for i := 2; i <= n; i++ {
		if dp[i] < minCost {
			minCost = dp[i]
		}
	}
	var caps []int
	for i := 1; i <= n; i++ {
		if dp[i] == minCost {
			caps = append(caps, i)
		}
	}
	return minCost, caps
}

func generateCaseD(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(10) + 2
	edges := make([]edge, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		if rng.Intn(2) == 0 {
			edges[i-2] = edge{p, i}
		} else {
			edges[i-2] = edge{i, p}
		}
	}
	return n, edges
}

func runCaseD(bin string, n int, edges []edge) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) < 1 {
		return fmt.Errorf("no output")
	}
	var gotMin int
	fmt.Sscan(fields[0], &gotMin)
	gotCap := make([]int, 0)
	for _, f := range fields[1:] {
		var v int
		fmt.Sscan(f, &v)
		gotCap = append(gotCap, v)
	}
	expMin, expCap := expectedAnswerD(n, edges)
	if gotMin != expMin || len(gotCap) != len(expCap) {
		return fmt.Errorf("expected %d %v got %d %v", expMin, expCap, gotMin, gotCap)
	}
	for i := range gotCap {
		if gotCap[i] != expCap[i] {
			return fmt.Errorf("expected %d %v got %d %v", expMin, expCap, gotMin, gotCap)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if err := runCaseD(bin, 2, []edge{{1, 2}}); err != nil {
		fmt.Fprintln(os.Stderr, "deterministic case failed:", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		n, edges := generateCaseD(rng)
		if err := runCaseD(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
