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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCaseD(n int, vals []int64, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}
	sum := make([]int64, n+1)
	f := make([]int64, n+1)
	const inf int64 = 1 << 60
	var ans int64
	hasAns := false
	for i := n - 1; i >= 0; i-- {
		u := order[i]
		sum[u] = vals[u-1]
		best1, best2 := -inf, -inf
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			sum[u] += sum[v]
			if f[v] > best1 {
				best2 = best1
				best1 = f[v]
			} else if f[v] > best2 {
				best2 = f[v]
			}
		}
		if best2 > -inf {
			cand := best1 + best2
			if !hasAns || cand > ans {
				ans = cand
				hasAns = true
			}
		}
		f[u] = sum[u]
		if best1 > f[u] {
			f[u] = best1
		}
	}
	if !hasAns {
		return "Impossible"
	}
	return fmt.Sprintf("%d", ans)
}

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 2
	vals := make([]int64, n)
	for i := range vals {
		vals[i] = int64(rng.Intn(21) - 10)
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	input := sb.String()
	expect := solveCaseD(n, vals, edges)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
