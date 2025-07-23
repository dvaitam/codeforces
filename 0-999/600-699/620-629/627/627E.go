package main

import (
	"bufio"
	"fmt"
	"os"
)

func countAtLeastK(arr []int, k int) int64 {
	left := 0
	sum := 0
	var res int64
	for right := 0; right < len(arr); right++ {
		sum += arr[right]
		for sum >= k {
			sum -= arr[left]
			left++
		}
		res += int64(left)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var r, c, n, k int
	if _, err := fmt.Fscan(reader, &r, &c, &n, &k); err != nil {
		return
	}

	grid := make([][]byte, r)
	for i := 0; i < r; i++ {
		grid[i] = make([]byte, c)
	}
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		grid[x-1][y-1] = 1
	}

	// ensure r <= c to reduce complexity
	if r > c {
		// transpose grid
		ng := make([][]byte, c)
		for i := 0; i < c; i++ {
			ng[i] = make([]byte, r)
			for j := 0; j < r; j++ {
				ng[i][j] = grid[j][i]
			}
		}
		grid = ng
		r, c = c, r
	}

	col := make([]int, c)
	var ans int64
	for top := 0; top < r; top++ {
		for j := 0; j < c; j++ {
			col[j] = 0
		}
		for bottom := top; bottom < r; bottom++ {
			for j := 0; j < c; j++ {
				if grid[bottom][j] == 1 {
					col[j]++
				}
			}
			ans += countAtLeastK(col, k)
		}
	}
	fmt.Fprintln(writer, ans)
}
