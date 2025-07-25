package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const maxN = 300005

var (
	g  [maxN][]int
	a  [maxN]int64
	dp [maxN][20]int64
)

func dfs(x, p int) {
	for i := 1; i < 20; i++ {
		dp[x][i] = int64(i) * a[x]
	}
	for _, u := range g[x] {
		if u == p {
			continue
		}
		dfs(u, x)
		for j := 1; j < 20; j++ {
			sum := int64(1<<63 - 1)
			for k := 1; k < 20; k++ {
				if j != k && dp[u][k] < sum {
					sum = dp[u][k]
				}
			}
			dp[x][j] += sum
		}
	}
}

func solve(n int) int64 {
	dfs(1, 0)
	ans := dp[1][1]
	for i := 2; i < 20; i++ {
		if dp[1][i] < ans {
			ans = dp[1][i]
		}
	}
	return ans
}

func generate() (string, string) {
	const T = 100
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	rand.Seed(4)
	for t := 0; t < T; t++ {
		n := rand.Intn(6) + 1
		fmt.Fprintf(&in, "%d\n", n)
		for i := 1; i <= n; i++ {
			a[i] = int64(rand.Intn(10) + 1)
			fmt.Fprintf(&in, "%d ", a[i])
			g[i] = g[i][:0]
		}
		fmt.Fprintf(&in, "\n")
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			g[p] = append(g[p], i)
			g[i] = append(g[i], p)
			fmt.Fprintf(&in, "%d %d\n", p, i)
		}
		res := solve(n)
		fmt.Fprintf(&out, "%d\n", res)
	}
	return in.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}
	got := buf.String()
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		fmt.Println("wrong answer")
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
