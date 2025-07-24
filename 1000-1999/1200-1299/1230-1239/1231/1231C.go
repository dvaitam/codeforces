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
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}

	const INF int = int(1e9)
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			bottom := INF
			right := INF
			if i+1 < n {
				bottom = a[i+1][j]
			}
			if j+1 < m {
				right = a[i][j+1]
			}
			if a[i][j] == 0 {
				val := bottom - 1
				if right-1 < val {
					val = right - 1
				}
				if val <= 0 {
					fmt.Fprintln(out, -1)
					return
				}
				a[i][j] = val
			} else {
				if a[i][j] >= bottom || a[i][j] >= right {
					fmt.Fprintln(out, -1)
					return
				}
			}
		}
	}

	var sum int64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j+1 < m && a[i][j] >= a[i][j+1] {
				fmt.Fprintln(out, -1)
				return
			}
			if i+1 < n && a[i][j] >= a[i+1][j] {
				fmt.Fprintln(out, -1)
				return
			}
			sum += int64(a[i][j])
		}
	}

	fmt.Fprintln(out, sum)
}
