package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func solveCase(n int, edges [][2]int) int64 {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}
	height := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		h := 0
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			if height[v]+1 > h {
				h = height[v] + 1
			}
		}
		height[u] = h
	}
	var total int64
	for i := 1; i <= n; i++ {
		total += int64(height[i] + 1)
	}
	pow := modPow(2, int64(n-1))
	ans := total % MOD * pow % MOD
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		u := i + 2
		v := rng.Intn(i+1) + 1
		edges[i] = [2]int{u, v}
	}
	ans := solveCase(n, edges)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %s got %s", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
