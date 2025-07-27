package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

type op struct {
	typ byte
	val int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	ops := make([]op, n)
	for i := 0; i < n; i++ {
		var t string
		fmt.Fscan(in, &t)
		if t == "+" {
			var x int64
			fmt.Fscan(in, &x)
			ops[i] = op{typ: '+', val: x}
		} else {
			ops[i] = op{typ: '-'}
		}
	}

	ans := int64(0)
	for i := 0; i < n; i++ {
		if ops[i].typ != '+' {
			continue
		}
		x := ops[i].val
		dp := make([][2]int64, n+1)
		dp[0][0] = 1
		for j := 0; j < n; j++ {
			ndp := make([][2]int64, n+1)
			if j == i {
				for k := 0; k <= n; k++ {
					ndp[k][1] = (ndp[k][1] + dp[k][0] + dp[k][1]) % mod
				}
				dp = ndp
				continue
			}
			if ops[j].typ == '+' {
				val := ops[j].val
				if (j < i && val <= x) || (j > i && val < x) {
					for k := 0; k <= n; k++ {
						for b := 0; b < 2; b++ {
							v := dp[k][b]
							if v == 0 {
								continue
							}
							ndp[k][b] = (ndp[k][b] + v) % mod
							if k+1 <= n {
								ndp[k+1][b] = (ndp[k+1][b] + v) % mod
							}
						}
					}
				} else {
					for k := 0; k <= n; k++ {
						for b := 0; b < 2; b++ {
							v := dp[k][b]
							if v == 0 {
								continue
							}
							ndp[k][b] = (ndp[k][b] + 2*v) % mod
						}
					}
				}
			} else { // '-'
				for k := 0; k <= n; k++ {
					for b := 0; b < 2; b++ {
						v := dp[k][b]
						if v == 0 {
							continue
						}
						ndp[k][b] = (ndp[k][b] + v) % mod
						if k > 0 {
							ndp[k-1][b] = (ndp[k-1][b] + v) % mod
						} else {
							ndp[k][0] = (ndp[k][0] + v) % mod
						}
					}
				}
			}
			dp = ndp
		}
		sum := int64(0)
		for k := 0; k <= n; k++ {
			sum = (sum + dp[k][1]) % mod
		}
		ans = (ans + sum*x) % mod
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans%mod)
	out.Flush()
}
