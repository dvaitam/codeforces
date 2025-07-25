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

func compute(n int, k int, special []int, edges [][2]int) int64 {
	total := 2 * k
	specialMark := make([]int, n+1)
	for _, v := range special {
		specialMark[v] = 1
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
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
	cnt := make([]int, n+1)
	var result int64
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		cnt[u] += specialMark[u]
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			cnt[u] += cnt[v]
		}
		if u != 1 {
			x := cnt[u]
			if x > total-x {
				x = total - x
			}
			result += int64(x)
		}
	}
	return result
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(18) + 2
	k := rng.Intn(n/2) + 1
	total := 2 * k
	perm := rng.Perm(n)
	special := make([]int, total)
	for i := 0; i < total; i++ {
		special[i] = perm[i] + 1
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		x := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{x, i}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < total; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", special[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	expected := compute(n, k, special, edges)
	return sb.String(), expected
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(outStr, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
