package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		base := int64(0)
		type pair struct{ A, B int64 }
		contested := make([]pair, 0)
		for i := 0; i < n; i++ {
			if a[i] > 0 && b[i] > 0 {
				contested = append(contested, pair{a[i] - 1, -(b[i] - 1)})
			} else if a[i] > 0 {
				base += a[i]
			} else if b[i] > 0 {
				base -= b[i]
			}
		}

		k := len(contested)
		const maxMask = 1 << 6
		var memo [maxMask][2]int64
		var vis [maxMask][2]bool

		var dfs func(mask int, turn int) int64
		dfs = func(mask int, turn int) int64 {
			if mask == 0 {
				return 0
			}
			if vis[mask][turn] {
				return memo[mask][turn]
			}
			vis[mask][turn] = true
			if turn == 0 { // Alice
				best := int64(-1 << 63)
				for i := 0; i < k; i++ {
					if mask&(1<<i) != 0 {
						val := contested[i].A + dfs(mask^(1<<i), 1)
						if val > best {
							best = val
						}
					}
				}
				memo[mask][turn] = best
			} else { // Bob
				best := int64(1<<63 - 1)
				for i := 0; i < k; i++ {
					if mask&(1<<i) != 0 {
						val := contested[i].B + dfs(mask^(1<<i), 0)
						if val < best {
							best = val
						}
					}
				}
				memo[mask][turn] = best
			}
			return memo[mask][turn]
		}

		res := base + dfs((1<<k)-1, 0)
		fmt.Fprintln(writer, res)
	}
}
