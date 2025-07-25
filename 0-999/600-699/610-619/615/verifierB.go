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

func expected(n int, edges [][2]int) int {
	deg := make([]int, n+1)
	prev := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		deg[u]++
		deg[v]++
		if u < v {
			prev[v] = append(prev[v], u)
		} else {
			prev[u] = append(prev[u], v)
		}
	}
	dp := make([]int, n+1)
	ans := 0
	for i := 1; i <= n; i++ {
		dp[i] = 1
		for _, p := range prev[i] {
			if dp[p]+1 > dp[i] {
				dp[i] = dp[p] + 1
			}
		}
		beauty := dp[i] * deg[i]
		if beauty > ans {
			ans = beauty
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges) + 1
	edges := make([][2]int, m)
	used := make(map[[2]int]bool)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		for {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			a, b := u, v
			if a > b {
				a, b = b, a
			}
			if used[[2]int{a, b}] {
				continue
			}
			used[[2]int{a, b}] = true
			edges[i] = [2]int{u, v}
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
			break
		}
	}
	return sb.String(), expected(n, edges)
}

func runCase(bin, input string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := []string{
		"2 1\n1 2\n",
	}
	exps := []int{2}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}

	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
