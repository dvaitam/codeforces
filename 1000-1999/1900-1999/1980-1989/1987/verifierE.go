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

func dfs(v int, g [][]int, a []int) (sum int, ops int) {
	if len(g[v]) == 0 {
		return a[v], 0
	}
	total := 0
	cost := 0
	for _, u := range g[v] {
		s, o := dfs(u, g, a)
		total += s
		cost += o
	}
	if total < a[v] {
		diff := a[v] - total
		cost += diff
		total += diff
	}
	return total + a[v], cost
}

func expected(n int, a []int, parent []int) int {
	g := make([][]int, n)
	for i := 1; i < n; i++ {
		g[parent[i]] = append(g[parent[i]], i)
	}
	_, ops := dfs(0, g, a)
	return ops
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(5)
	}
	parent := make([]int, n)
	for i := 1; i < n; i++ {
		parent[i] = rng.Intn(i)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 1; i < n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", parent[i]+1))
	}
	if n > 1 {
		sb.WriteByte('\n')
	}
	return sb.String(), expected(n, a, parent)
}

func runCase(bin string, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
