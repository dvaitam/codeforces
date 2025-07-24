package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1_000_000_007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		ng := make([]int, n)
		stack := make([]int, 0)
		for i := n - 1; i >= 0; i-- {
			for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				ng[i] = n
			} else {
				ng[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}

		children := make([][]int, n+1)
		for i := 0; i < n; i++ {
			p := ng[i]
			children[p] = append(children[p], i)
		}
		for i := 0; i <= n; i++ {
			sort.Ints(children[i])
		}

		dp := make([][]int64, n+1)
		after := make([][]int64, n+1)
		for i := 0; i <= n; i++ {
			arr := make([]int64, m+1)
			for v := 0; v <= m; v++ {
				arr[v] = 1
			}
			for idx := len(children[i]) - 1; idx >= 0; idx-- {
				c := children[i][idx]
				tmp := make([]int64, m+1)
				for v := 1; v <= m; v++ {
					tmp[v] = dp[c][v] * arr[v] % MOD
				}
				arr[0] = 0
				var s int64
				for v := 1; v <= m; v++ {
					s += tmp[v]
					if s >= MOD {
						s -= MOD
					}
					arr[v] = s
				}
			}
			after[i] = arr
			dp[i] = make([]int64, m+1)
			for v := 1; v <= m; v++ {
				dp[i][v] = arr[v-1]
			}
		}
		fmt.Fprintln(out, after[n][m]%MOD)
	}
}
