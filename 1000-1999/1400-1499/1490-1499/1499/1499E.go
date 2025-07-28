package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var x, y string
	fmt.Fscan(in, &x)
	fmt.Fscan(in, &y)
	n := len(x)
	m := len(y)

	dp0Only := make([][]int, n+1)
	dp1Only := make([][]int, n+1)
	dp0Both := make([][]int, n+1)
	dp1Both := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp0Only[i] = make([]int, m+1)
		dp1Only[i] = make([]int, m+1)
		dp0Both[i] = make([]int, m+1)
		dp1Both[i] = make([]int, m+1)
	}

	// initialize starting states for all substring pairs
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			dp0Only[i+1][j] = (dp0Only[i+1][j] + 1) % mod
			dp1Only[i][j+1] = (dp1Only[i][j+1] + 1) % mod
		}
	}

	ans := 0
	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			if i > 0 {
				last := x[i-1]
				val := dp0Only[i][j]
				if val != 0 {
					if i < n && x[i] != last {
						dp0Only[i+1][j] = (dp0Only[i+1][j] + val) % mod
					}
					if j < m && y[j] != last {
						dp1Both[i][j+1] = (dp1Both[i][j+1] + val) % mod
					}
				}
				val = dp0Both[i][j]
				if val != 0 {
					if i < n && x[i] != last {
						dp0Both[i+1][j] = (dp0Both[i+1][j] + val) % mod
					}
					if j < m && y[j] != last {
						dp1Both[i][j+1] = (dp1Both[i][j+1] + val) % mod
					}
				}
			}
			if j > 0 {
				last := y[j-1]
				val := dp1Only[i][j]
				if val != 0 {
					if j < m && y[j] != last {
						dp1Only[i][j+1] = (dp1Only[i][j+1] + val) % mod
					}
					if i < n && x[i] != last {
						dp0Both[i+1][j] = (dp0Both[i+1][j] + val) % mod
					}
				}
				val = dp1Both[i][j]
				if val != 0 {
					if j < m && y[j] != last {
						dp1Both[i][j+1] = (dp1Both[i][j+1] + val) % mod
					}
					if i < n && x[i] != last {
						dp0Both[i+1][j] = (dp0Both[i+1][j] + val) % mod
					}
				}
			}
			ans += dp0Both[i][j] + dp1Both[i][j]
			ans %= mod
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans%mod)
	out.Flush()
}
