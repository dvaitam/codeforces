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

const MOD = 1000000007

type edge struct{ u, v int }

func add(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}
func mul(a, b int) int { return int((int64(a) * int64(b)) % MOD) }

func solveE(n, k int, edges []edge) int {
	g := make([][]int, n)
	for _, e := range edges {
		g[e.u] = append(g[e.u], e.v)
		g[e.v] = append(g[e.v], e.u)
	}
	dp := make([][]int, n)
	var dfs func(int, int)
	dfs = func(v, p int) {
		dp[v] = make([]int, k+2)
		dp0 := make([]int, k+2)
		dp0[k+1] = 1
		prodAll := 1
		for _, u := range g[v] {
			if u == p {
				continue
			}
			dfs(u, v)
			new0 := make([]int, k+2)
			for d0 := 0; d0 <= k+1; d0++ {
				if dp0[d0] == 0 {
					continue
				}
				for du := 0; du <= k+1; du++ {
					if dp[u][du] == 0 {
						continue
					}
					nd := du + 1
					if nd > k+1 {
						nd = k + 1
					}
					d := d0
					if nd < d {
						d = nd
					}
					new0[d] = (new0[d] + dp0[d0]*dp[u][du]) % MOD
				}
			}
			dp0 = new0
			s := 0
			for du := 0; du <= k+1; du++ {
				s = (s + dp[u][du]) % MOD
			}
			prodAll = prodAll * s % MOD
		}
		dp[v][0] = prodAll
		for d := 0; d <= k+1; d++ {
			dp[v][d] = add(dp[v][d], dp0[d])
		}
	}
	dfs(0, -1)
	res := 0
	for d := 0; d <= k; d++ {
		res = add(res, dp[0][d])
	}
	return res
}

func genTestE() (int, int, []edge) {
	n := rand.Intn(10) + 1
	k := rand.Intn(min(20, n-1) + 1)
	edges := make([]edge, 0, n-1)
	for i := 1; i < n; i++ {
		p := rand.Intn(i)
		edges = append(edges, edge{p, i})
	}
	return n, k, edges
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runBinary(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	path := os.Args[1]
	for i := 0; i < 100; i++ {
		n, k, edges := genTestE()
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, k)
		for _, e := range edges {
			fmt.Fprintf(&b, "%d %d\n", e.u+1, e.v+1)
		}
		input := b.String()
		expected := solveE(n, k, edges)
		gotStr, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var got int
		_, err = fmt.Sscanf(gotStr, "%d", &got)
		if err != nil {
			fmt.Printf("test %d: parse output error: %v\ninput:\n%soutput:%s\n", i+1, err, input, gotStr)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%d\ngot:%d\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
