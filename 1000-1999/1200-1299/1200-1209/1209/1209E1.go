package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([][]int, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}

		type col struct{ idx, mx int }
		cols := make([]col, m)
		for j := 0; j < m; j++ {
			mx := 0
			for i := 0; i < n; i++ {
				if a[i][j] > mx {
					mx = a[i][j]
				}
			}
			cols[j] = col{j, mx}
		}
		sort.Slice(cols, func(i, j int) bool { return cols[i].mx > cols[j].mx })
		if m > n {
			cols = cols[:n]
		}

		k := len(cols)
		size := 1 << n
		dp := make([]int, size)
		for i := 0; i < k; i++ {
			best := make([]int, size)
			j := cols[i].idx
			for rot := 0; rot < n; rot++ {
				vals := make([]int, n)
				for r := 0; r < n; r++ {
					vals[r] = a[(r+rot)%n][j]
				}
				for mask := 0; mask < size; mask++ {
					sum := 0
					for b := 0; b < n; b++ {
						if (mask>>b)&1 == 1 {
							sum += vals[b]
						}
					}
					if sum > best[mask] {
						best[mask] = sum
					}
				}
			}
			newDP := make([]int, size)
			copy(newDP, dp)
			for mask := 0; mask < size; mask++ {
				for sub := 1; sub < size; sub++ {
					nm := mask | sub
					if val := dp[mask] + best[sub]; val > newDP[nm] {
						newDP[nm] = val
					}
				}
			}
			dp = newDP
		}
		fmt.Fprintln(out, dp[size-1])
	}
}
