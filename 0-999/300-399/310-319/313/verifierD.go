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

// run executes the candidate binary with given input
func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// compute minimal cost expected answer for a test case
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solveExpected(n, m, k int, L, R []int, C []int64) int64 {
	const INF int64 = 1 << 60
	// M[l][r] = minimal cost to buy exact [l,r]
	mtx := make([][]int64, n+2)
	for i := range mtx {
		mtx[i] = make([]int64, n+2)
		for j := range mtx[i] {
			mtx[i][j] = INF
		}
	}
	for i := 0; i < m; i++ {
		l, r := L[i], R[i]
		if C[i] < mtx[l][r] {
			mtx[l][r] = C[i]
		}
	}
	// G[l][r] minimal cost to fully cover [l,r]
	g := make([][]int64, n+2)
	for i := range g {
		g[i] = make([]int64, n+2)
	}
	for r := 1; r <= n; r++ {
		// smin[l][x] = min_{t in [x..r]} mtx[l][t]
		smin := make([][]int64, n+2)
		for i := range smin {
			smin[i] = make([]int64, n+2)
		}
		for l := r; l >= 1; l-- {
			cur := INF
			for x := r; x >= l; x-- {
				if mtx[l][x] < cur {
					cur = mtx[l][x]
				}
				smin[l][x] = cur
			}
		}
		best := make([]int64, n+3)
		for i := range best {
			best[i] = INF
		}
		best[r+1] = 0
		for l := r; l >= 1; l-- {
			val := int64(1 << 60)
			for x := l; x <= r; x++ {
				sr := smin[l][x]
				if sr >= INF || best[x+1] >= INF {
					continue
				}
				v := sr + best[x+1]
				if v < val {
					val = v
				}
			}
			best[l] = val
		}
		for l := 1; l <= r; l++ {
			g[l][r] = best[l]
		}
	}
	// DP over disjoint segments to cover at least k positions
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, n+1)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		for t := 0; t <= n; t++ {
			dp[i][t] = dp[i-1][t]
		}
		for l := 1; l <= i; l++ {
			cost := g[l][i]
			if cost >= INF {
				continue
			}
			len := i - l + 1
			for t := len; t <= n; t++ {
				if dp[l-1][t-len] >= INF {
					continue
				}
				v := dp[l-1][t-len] + cost
				if v < dp[i][t] {
					dp[i][t] = v
				}
			}
		}
	}
	ans := int64(1 << 60)
	for t := k; t <= n; t++ {
		if dp[n][t] < ans {
			ans = dp[n][t]
		}
	}
	if ans >= INF {
		return -1
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("bad test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		Ls := make([]int, m)
		Rs := make([]int, m)
		Cs := make([]int64, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			l, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			r, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			c, _ := strconv.Atoi(scan.Text())
			input.WriteString(fmt.Sprintf("%d %d %d\n", l, r, c))
			Ls[i], Rs[i], Cs[i] = l, r, int64(c)
		}
		expected := solveExpected(n, m, k, Ls, Rs, Cs)
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(gotStr))
		if len(fields) == 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: empty output\n", caseNum)
			os.Exit(1)
		}
		got, perr := strconv.ParseInt(fields[0], 10, 64)
		if perr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: output is not an integer: %q\n", caseNum, gotStr)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
