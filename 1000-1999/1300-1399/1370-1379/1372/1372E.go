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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	intervals := make([][2]int, 0)
	for i := 0; i < n; i++ {
		var k int
		fmt.Fscan(in, &k)
		for j := 0; j < k; j++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			intervals = append(intervals, [2]int{l, r})
		}
	}

	m1 := m + 2
	arr := make([][][]int, m1)
	pre := make([][][]int, m1)
	for k := 0; k <= m; k++ {
		arr[k] = make([][]int, m1)
		pre[k] = make([][]int, m1)
		for i := 0; i <= m; i++ {
			arr[k][i] = make([]int, m1)
			pre[k][i] = make([]int, m1)
		}
	}

	for _, it := range intervals {
		L, R := it[0], it[1]
		for k := L; k <= R; k++ {
			arr[k][L][R]++
		}
	}

	for k := 1; k <= m; k++ {
		for l := m; l >= 1; l-- {
			for r := l; r <= m; r++ {
				v := arr[k][l][r]
				if l+1 <= m {
					v += pre[k][l+1][r]
				}
				if r-1 >= 1 {
					v += pre[k][l][r-1]
				}
				if l+1 <= m && r-1 >= 1 {
					v -= pre[k][l+1][r-1]
				}
				pre[k][l][r] = v
			}
		}
	}

	dp := make([][]int, m1)
	for i := range dp {
		dp[i] = make([]int, m1)
	}

	for length := 1; length <= m; length++ {
		for l := 1; l+length-1 <= m; l++ {
			r := l + length - 1
			best := 0
			for k := l; k <= r; k++ {
				cnt := pre[k][l][r]
				v := cnt * cnt
				if k > l {
					v += dp[l][k-1]
				}
				if k < r {
					v += dp[k+1][r]
				}
				if v > best {
					best = v
				}
			}
			dp[l][r] = best
		}
	}

	fmt.Fprintln(out, dp[1][m])
}
