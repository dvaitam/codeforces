package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var S, T string
	if _, err := fmt.Fscan(in, &S); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	n := len(S)
	m := len(T)
	s := []byte(S)
	t := []byte(T)

	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
	}

	if m <= n {
		for i := 0; i < n; i++ {
			if i >= m || t[i] == s[0] {
				dp[i][i] = 2
			}
		}

		for length := 1; length < n; length++ {
			c := s[length]
			for l := 0; l+length-1 < n; l++ {
				r := l + length - 1
				val := dp[l][r]
				if val == 0 {
					continue
				}
				if l > 0 {
					if l-1 >= m || t[l-1] == c {
						dp[l-1][r] = (dp[l-1][r] + val) % MOD
					}
				}
				if r+1 < n {
					if r+1 >= m || t[r+1] == c {
						dp[l][r+1] = (dp[l][r+1] + val) % MOD
					}
				}
			}
		}

		ans := 0
		for i := m - 1; i < n; i++ {
			ans = (ans + dp[0][i]) % MOD
		}
		fmt.Println(ans)
	}
}
