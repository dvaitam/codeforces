package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 1_000_000_007

func powmod(a, e int64) int64 {
	a %= MOD
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func expected(n, k int, edges [][2]int) string {
	if k == 1 || k == 3 {
		return "1"
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	sz := make([]int, n+1)
	var sumDist int64
	var dfs func(int, int)
	dfs = func(v, p int) {
		sz[v] = 1
		for _, to := range adj[v] {
			if to == p {
				continue
			}
			dfs(to, v)
			s := sz[to]
			sumDist += int64(s) * int64(n-s)
			sz[v] += s
		}
	}
	dfs(1, 0)
	numerator := (2 * (sumDist % MOD)) % MOD
	denom := int64(n) * int64(n-1) % MOD
	invDenom := powmod(denom, MOD-2)
	ans := (1 + numerator*invDenom%MOD) % MOD
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesB1.txt")
	if err != nil {
		fmt.Println("could not read testcasesB1.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			edges[j] = [2]int{u, v}
		}

		var sb bytes.Buffer
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		exp := expected(n, k, edges)
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(sb.Bytes())
		res, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx+1, err, res)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(res))
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", caseIdx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
