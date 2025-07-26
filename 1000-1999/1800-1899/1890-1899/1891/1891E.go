package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		const INF = int(1e9)
		dp0 := make([]int, k+1)
		dp1 := make([]int, k+1)
		for j := 0; j <= k; j++ {
			dp0[j] = INF
			dp1[j] = INF
		}
		dp0[0] = 0
		if k >= 1 {
			dp1[1] = 0
		}
		for i := 2; i <= n; i++ {
			n0 := make([]int, k+1)
			n1 := make([]int, k+1)
			for j := 0; j <= k; j++ {
				n0[j] = INF
				n1[j] = INF
			}
			for c := 0; c <= k; c++ {
				if dp0[c] != INF {
					add := 0
					if gcd(a[i-1], a[i]) == 1 {
						add = 1
					}
					if dp0[c]+add < n0[c] {
						n0[c] = dp0[c] + add
					}
					if c+1 <= k {
						add2 := 0
						if a[i-1] == 1 {
							add2 = 1
						}
						if dp0[c]+add2 < n1[c+1] {
							n1[c+1] = dp0[c] + add2
						}
					}
				}
				if dp1[c] != INF {
					add := 0
					if a[i] == 1 {
						add = 1
					}
					if dp1[c]+add < n0[c] {
						n0[c] = dp1[c] + add
					}
					if c+1 <= k {
						add2 := 0 // pair(i-1,i) after zero i both 0 -> gcd=0 -> not 1
						if dp1[c]+add2 < n1[c+1] {
							n1[c+1] = dp1[c] + add2
						}
					}
				}
			}
			dp0, dp1 = n0, n1
		}
		ans := INF
		for c := 0; c <= k; c++ {
			if dp0[c] < ans {
				ans = dp0[c]
			}
			if dp1[c] < ans {
				ans = dp1[c]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
