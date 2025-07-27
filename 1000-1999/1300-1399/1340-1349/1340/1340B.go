package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	s := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}

	pat := []string{
		"1110111",
		"0010010",
		"1011101",
		"1011011",
		"0111010",
		"1101011",
		"1101111",
		"1010010",
		"1111111",
		"1111011",
	}

	cost := make([][10]int, n)
	for i := 0; i < n; i++ {
		for d := 0; d < 10; d++ {
			diff := 0
			ok := true
			for j := 0; j < 7; j++ {
				if s[i][j] == '1' && pat[d][j] == '0' {
					ok = false
					break
				}
				if s[i][j] == '0' && pat[d][j] == '1' {
					diff++
				}
			}
			if ok {
				cost[i][d] = diff
			} else {
				cost[i][d] = -1
			}
		}
	}

	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, k+1)
	}
	dp[n][0] = true
	for i := n - 1; i >= 0; i-- {
		for used := 0; used <= k; used++ {
			for d := 0; d < 10; d++ {
				c := cost[i][d]
				if c >= 0 && used >= c && dp[i+1][used-c] {
					dp[i][used] = true
					break
				}
			}
		}
	}

	if !dp[0][k] {
		fmt.Fprintln(out, -1)
		return
	}

	remaining := k
	result := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		for d := 9; d >= 0; d-- {
			c := cost[i][d]
			if c >= 0 && remaining >= c && dp[i+1][remaining-c] {
				result = append(result, byte('0'+d))
				remaining -= c
				break
			}
		}
	}
	fmt.Fprintln(out, string(result))
}
