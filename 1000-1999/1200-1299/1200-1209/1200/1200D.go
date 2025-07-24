package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		grid[i] = []byte(s)
	}

	base := 0
	size := n - k + 1
	if size < 0 {
		size = 0
	}
	diff := make([][]int, size+2)
	for i := range diff {
		diff[i] = make([]int, size+2)
	}

	// process rows
	for r := 0; r < n; r++ {
		left, right := -1, -1
		for c := 0; c < n; c++ {
			if grid[r][c] == 'B' {
				if left == -1 {
					left = c
				}
				right = c
			}
		}
		if left == -1 {
			base++
			continue
		}
		if right-left+1 > k {
			continue
		}
		x1 := max(0, r-k+1)
		x2 := min(r, n-k)
		y1 := max(0, right-k+1)
		y2 := min(left, n-k)
		if x1 <= x2 && y1 <= y2 {
			diff[x1][y1]++
			diff[x1][y2+1]--
			diff[x2+1][y1]--
			diff[x2+1][y2+1]++
		}
	}

	// process columns
	for c := 0; c < n; c++ {
		top, bottom := -1, -1
		for r := 0; r < n; r++ {
			if grid[r][c] == 'B' {
				if top == -1 {
					top = r
				}
				bottom = r
			}
		}
		if top == -1 {
			base++
			continue
		}
		if bottom-top+1 > k {
			continue
		}
		x1 := max(0, bottom-k+1)
		x2 := min(top, n-k)
		y1 := max(0, c-k+1)
		y2 := min(c, n-k)
		if x1 <= x2 && y1 <= y2 {
			diff[x1][y1]++
			diff[x1][y2+1]--
			diff[x2+1][y1]--
			diff[x2+1][y2+1]++
		}
	}

	ans := 0
	for i := 0; i <= size; i++ {
		for j := 1; j <= size; j++ {
			diff[i][j] += diff[i][j-1]
		}
	}
	for i := 1; i <= size; i++ {
		for j := 0; j <= size; j++ {
			diff[i][j] += diff[i-1][j]
		}
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if diff[i][j] > ans {
				ans = diff[i][j]
			}
		}
	}

	fmt.Fprintln(writer, base+ans)
}
