package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func dfs(u, p int, adj [][]int) ([]*big.Int, int) {
	dp := make([]*big.Int, 2)
	dp[1] = big.NewInt(1)
	sz := 1
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		dpv, sv := dfs(v, u, adj)
		newDp := make([]*big.Int, sz+sv+1)
		for su := 1; su <= sz; su++ {
			if dp[su] == nil {
				continue
			}
			for svv := 1; svv < len(dpv); svv++ {
				if dpv[svv] == nil {
					continue
				}
				nk := su + svv
				prod1 := new(big.Int).Mul(dp[su], dpv[svv])
				if newDp[nk] == nil || prod1.Cmp(newDp[nk]) > 0 {
					newDp[nk] = prod1
				}
				prod2 := new(big.Int).Mul(dp[su], dpv[svv])
				prod2.Mul(prod2, big.NewInt(int64(svv)))
				if newDp[su] == nil || prod2.Cmp(newDp[su]) > 0 {
					newDp[su] = prod2
				}
			}
		}
		sz += sv
		dp = newDp
	}
	return dp, sz
}

func solveE(n int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dp, _ := dfs(1, 0, adj)
	ans := big.NewInt(0)
	for k, prod := range dp {
		if prod == nil {
			continue
		}
		total := new(big.Int).Mul(prod, big.NewInt(int64(k)))
		if total.Cmp(ans) > 0 {
			ans = total
		}
	}
	return ans.String()
}

func genTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	edges := genTree(rng, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String(), solveE(n, edges)
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
