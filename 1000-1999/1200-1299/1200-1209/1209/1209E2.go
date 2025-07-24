package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([][]int, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}
		type Col struct {
			idx int
			mx  int
		}
		cols := make([]Col, m)
		for j := 0; j < m; j++ {
			mx := 0
			for i := 0; i < n; i++ {
				if a[i][j] > mx {
					mx = a[i][j]
				}
			}
			cols[j] = Col{idx: j, mx: mx}
		}
		sort.Slice(cols, func(i, j int) bool { return cols[i].mx > cols[j].mx })
		k := n
		if k > m {
			k = m
		}
		cols = cols[:k]
		size := 1 << n
		dp := make([]int, size)
		for _, c := range cols {
			idx := c.idx
			colVals := make([]int, n)
			for i := 0; i < n; i++ {
				colVals[i] = a[i][idx]
			}
			best := make([]int, size)
			for r := 0; r < n; r++ {
				rotated := make([]int, n)
				for i := 0; i < n; i++ {
					rotated[i] = colVals[(i+r)%n]
				}
				for mask := 1; mask < size; mask++ {
					sum := 0
					for bit := 0; bit < n; bit++ {
						if mask&(1<<bit) != 0 {
							sum += rotated[bit]
						}
					}
					if sum > best[mask] {
						best[mask] = sum
					}
				}
			}
			ndp := make([]int, size)
			copy(ndp, dp)
			for mask := 0; mask < size; mask++ {
				remain := (size - 1) ^ mask
				sub := remain
				for sub > 0 {
					val := dp[mask] + best[sub]
					if val > ndp[mask|sub] {
						ndp[mask|sub] = val
					}
					sub = (sub - 1) & remain
				}
			}
			dp = ndp
		}
		fmt.Println(dp[size-1])
	}
}
