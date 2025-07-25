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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}

		minRw, maxRw, minCw, maxCw := n+1, 0, m+1, 0
		minRb, maxRb, minCb, maxCb := n+1, 0, m+1, 0
		countW, countB := 0, 0

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				ch := grid[i][j]
				if ch == 'W' {
					countW++
					if i+1 < minRw {
						minRw = i + 1
					}
					if i+1 > maxRw {
						maxRw = i + 1
					}
					if j+1 < minCw {
						minCw = j + 1
					}
					if j+1 > maxCw {
						maxCw = j + 1
					}
				} else if ch == 'B' {
					countB++
					if i+1 < minRb {
						minRb = i + 1
					}
					if i+1 > maxRb {
						maxRb = i + 1
					}
					if j+1 < minCb {
						minCb = j + 1
					}
					if j+1 > maxCb {
						maxCb = j + 1
					}
				}
			}
		}

		possible := false
		if countW > 0 && minRw == 1 && maxRw == n && minCw == 1 && maxCw == m {
			possible = true
		}
		if countB > 0 && minRb == 1 && maxRb == n && minCb == 1 && maxCb == m {
			possible = true
		}

		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
