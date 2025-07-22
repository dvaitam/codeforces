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

type edge struct{ u, v int }

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n int, weights []int64, edges []edge) int64 {
	tree := make([][]int, n+1)
	for _, e := range edges {
		tree[e.u] = append(tree[e.u], e.v)
		tree[e.v] = append(tree[e.v], e.u)
	}
	childs := make([][]int, n+1)
	parent := make([]int, n+1)
	parent[1] = 0
	queue := []int{1}
	for i := 0; i < len(queue); i++ {
		v := queue[i]
		for _, u := range tree[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			childs[v] = append(childs[v], u)
			queue = append(queue, u)
		}
	}
	var dfs func(int) (int64, int64)
	dfs = func(v int) (int64, int64) {
		if len(childs[v]) == 0 {
			return 0, weights[v]
		}
		var totalRem int64
		var sumW int64
		minW := int64(1<<62 - 1)
		for _, u := range childs[v] {
			r, w := dfs(u)
			totalRem += r
			sumW += w
			if w < minW {
				minW = w
			}
		}
		cnt := int64(len(childs[v]))
		totalRem += sumW - minW*cnt
		return totalRem, minW * cnt
	}
	rem, _ := dfs(1)
	return rem
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	edges := make([]edge, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = edge{p, i}
	}
	deg := make([]int, n+1)
	for _, e := range edges {
		deg[e.u]++
		deg[e.v]++
	}
	weights := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 && i != 1 {
			weights[i] = rng.Int63n(20) + 1
		} else {
			weights[i] = 0
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", weights[i]))
	}
	sb.WriteString("\n")
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	expect := fmt.Sprintf("%d", solveCase(n, weights, edges))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
