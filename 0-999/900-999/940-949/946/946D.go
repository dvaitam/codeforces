package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = int(1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	dp := make([]int, k+1)
	for i := 1; i <= k; i++ {
		dp[i] = inf
	}

	for day := 0; day < n; day++ {
		var s string
		fmt.Fscan(in, &s)
		positions := make([]int, 0, m)
		for i := 0; i < m; i++ {
			if s[i] == '1' {
				positions = append(positions, i)
			}
		}
		r := len(positions)
		arr := make([]int, r+1)
		for i := 0; i <= r; i++ {
			arr[i] = inf
		}
		if r == 0 {
			arr[0] = 0
		} else {
			arr[r] = 0
			for l := 0; l < r; l++ {
				for rr := l; rr < r; rr++ {
					removed := l + (r - 1 - rr)
					length := positions[rr] - positions[l] + 1
					if length < arr[removed] {
						arr[removed] = length
					}
				}
			}
		}

		newdp := make([]int, k+1)
		for i := 0; i <= k; i++ {
			newdp[i] = inf
		}
		for used := 0; used <= k; used++ {
			if dp[used] == inf {
				continue
			}
			for rem := 0; rem <= r && used+rem <= k; rem++ {
				cost := arr[rem]
				val := dp[used] + cost
				if val < newdp[used+rem] {
					newdp[used+rem] = val
				}
			}
		}
		dp = newdp
	}

	ans := inf
	for i := 0; i <= k; i++ {
		if dp[i] < ans {
			ans = dp[i]
		}
	}
	fmt.Println(ans)
}
