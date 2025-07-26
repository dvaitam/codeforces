package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func inv(x int64) int64 { return modPow(x, mod-2) }

func hasMatch(adj [][]bool) bool {
	n := len(adj)
	match := make([]int, n)
	for i := range match {
		match[i] = -1
	}
	var dfs func(int, []bool) bool
	dfs = func(v int, used []bool) bool {
		for j := 0; j < n; j++ {
			if adj[v][j] && !used[j] {
				used[j] = true
				if match[j] == -1 || dfs(match[j], used) {
					match[j] = v
					return true
				}
			}
		}
		return false
	}
	cnt := 0
	for v := 0; v < n; v++ {
		used := make([]bool, n)
		if dfs(v, used) {
			cnt++
		}
	}
	return cnt == n
}

func bruteProb(p [][]int) int64 {
	n := len(p)
	inv100 := inv(100)
	q := make([][]int64, n)
	r := make([][]int64, n)
	for i := 0; i < n; i++ {
		q[i] = make([]int64, n)
		r[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			q[i][j] = int64(p[i][j]) * inv100 % mod
			r[i][j] = (1 - q[i][j] + mod) % mod
		}
	}
	edges := n * n
	ans := int64(0)
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for mask := 0; mask < (1 << edges); mask++ {
		prob := int64(1)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				bit := (mask >> (i*n + j)) & 1
				if bit == 1 {
					prob = prob * q[i][j] % mod
					adj[i][j] = true
				} else {
					prob = prob * r[i][j] % mod
					adj[i][j] = false
				}
			}
		}
		if hasMatch(adj) {
			ans = (ans + prob) % mod
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	p := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		p[i] = make([]int, n)
		for j := 0; j < n; j++ {
			val := rng.Intn(101)
			p[i][j] = val
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	ans := bruteProb(p)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, want := genCase(rng)
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
