package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const mod int = 998244353

var (
	n, k int
	adj  [][]int
)

func dfs(v, p int) []int {
	dp := make([]int, k+1)
	dp[0] = 1
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		child := dfs(to, v)
		sumChild := 0
		for _, x := range child {
			sumChild += x
			if sumChild >= mod {
				sumChild -= mod
			}
		}
		newDp := make([]int, k+1)
		for i := 0; i <= k; i++ {
			if dp[i] == 0 {
				continue
			}
			val := dp[i] * sumChild % mod
			newDp[i] = (newDp[i] + val) % mod
		}
		for i := 0; i <= k; i++ {
			if dp[i] == 0 {
				continue
			}
			for j := 0; j <= k; j++ {
				if child[j] == 0 {
					continue
				}
				if i+j+1 > k {
					continue
				}
				nd := i
				if j+1 > nd {
					nd = j + 1
				}
				val := (dp[i] * child[j]) % mod
				newDp[nd] += val
				if newDp[nd] >= mod {
					newDp[nd] -= mod
				}
			}
		}
		dp = newDp
	}
	return dp
}

func solveCase(nv, kv int, edges [][2]int) string {
	n = nv
	k = kv
	adj = make([][]int, n)
	for _, e := range edges {
		u := e[0] - 1
		v := e[1] - 1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dp := dfs(0, -1)
	ans := 0
	for _, x := range dp {
		ans += x
		if ans >= mod {
			ans -= mod
		}
	}
	return fmt.Sprintf("%d", ans%mod)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		nv, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		kv, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, nv-1)
		for j := 0; j < nv-1; j++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			edges[j] = [2]int{u, v}
		}
		expected[i] = solveCase(nv, kv, edges)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
