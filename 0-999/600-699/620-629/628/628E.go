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
	fmt.Fscan(in, &n, &m)
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}

	right := make([][]int, n)
	for i := 0; i < n; i++ {
		right[i] = make([]int, m)
		cnt := 0
		for j := m - 1; j >= 0; j-- {
			if grid[i][j] == 'z' {
				cnt++
			} else {
				cnt = 0
			}
			right[i][j] = cnt
		}
	}

	diag := make([][]int, n)
	for i := range diag {
		diag[i] = make([]int, m)
	}
	for i := n - 1; i >= 0; i-- {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'z' {
				if i+1 < n && j-1 >= 0 {
					diag[i][j] = diag[i+1][j-1] + 1
				} else {
					diag[i][j] = 1
				}
			} else {
				diag[i][j] = 0
			}
		}
	}

	ans := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 'z' {
				continue
			}
			maxL := right[i][j]
			if i+1 < maxL {
				maxL = i + 1
			}
			for l := 1; l <= maxL; l++ {
				top := i - l + 1
				if right[top][j] < l {
					break
				}
				if j+l-1 >= m {
					break
				}
				if diag[top][j+l-1] < l {
					break
				}
				ans++
			}
		}
	}

	fmt.Fprintln(out, ans)
}
