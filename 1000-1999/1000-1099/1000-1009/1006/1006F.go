package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for {
		var n, m int
		var k int64
		if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
			break
		}
		A := make([][]int64, n)
		for i := 0; i < n; i++ {
			A[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &A[i][j])
			}
		}
		// C maps for first half
		C := make([][]map[int64]int, n)
		for i := 0; i < n; i++ {
			C[i] = make([]map[int64]int, m)
			for j := 0; j < m; j++ {
				C[i][j] = make(map[int64]int)
			}
		}
		var ans int64
		d := min(n, m)
		total := n * m
		if total == 1 {
			if A[0][0] == k {
				ans = 1
			}
		} else {
			// first half dfs
			var dfs1 func(r, c, depth, inc, t int, sum int64)
			dfs1 = func(r, c, depth, inc, t int, sum int64) {
				if r < 0 || r >= n || c < 0 || c >= m {
					return
				}
				sum ^= A[r][c]
				if depth == t {
					C[r][c][sum]++
					return
				}
				dfs1(r+inc, c, depth+1, inc, t, sum)
				dfs1(r, c+inc, depth+1, inc, t, sum)
			}
			// second half dfs
			var dfs2 func(r, c, depth, inc, t int, sum int64)
			dfs2 = func(r, c, depth, inc, t int, sum int64) {
				if r < 0 || r >= n || c < 0 || c >= m {
					return
				}
				sum ^= A[r][c]
				if depth == t {
					// match with maps
					if r+inc >= 0 && r+inc < n {
						ans += int64(C[r+inc][c][sum^k])
					}
					if c+inc >= 0 && c+inc < m {
						ans += int64(C[r][c+inc][sum^k])
					}
					return
				}
				dfs2(r+inc, c, depth+1, inc, t, sum)
				dfs2(r, c+inc, depth+1, inc, t, sum)
			}
			dfs1(0, 0, 0, 1, d-1, 0)
			dfs2(n-1, m-1, 0, -1, n+m-2-d, 0)
		}
		fmt.Fprintln(writer, ans)
	}
}
