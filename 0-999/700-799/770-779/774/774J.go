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
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true

	for idx := 0; idx < n; idx++ {
		ndp := make([][]bool, n+1)
		for i := range ndp {
			ndp[i] = make([]bool, n+1)
		}
		for j := 0; j <= n; j++ {
			for m := 0; m <= n; m++ {
				if !dp[j][m] {
					continue
				}
				c := s[idx]
				if c != 'N' { // choose Y
					newJ := 0
					newM := m
					if j > newM {
						newM = j
					}
					ndp[newJ][newM] = true
				}
				if c != 'Y' { // choose N
					newJ := j + 1
					newM := m
					if newJ > newM {
						newM = newJ
					}
					if newJ <= n && newM <= n {
						ndp[newJ][newM] = true
					}
				}
			}
		}
		dp = ndp
	}

	ans := false
	for j := 0; j <= n; j++ {
		if dp[j][k] {
			ans = true
			break
		}
	}

	if ans {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
