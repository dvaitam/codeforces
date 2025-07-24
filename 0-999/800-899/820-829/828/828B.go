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

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &grid[i])
	}

	minR, maxR := n, -1
	minC, maxC := m, -1
	count := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'B' {
				count++
				if i < minR {
					minR = i
				}
				if i > maxR {
					maxR = i
				}
				if j < minC {
					minC = j
				}
				if j > maxC {
					maxC = j
				}
			}
		}
	}

	if count == 0 {
		fmt.Fprintln(out, 1)
		return
	}

	height := maxR - minR + 1
	width := maxC - minC + 1
	side := height
	if width > side {
		side = width
	}

	if side > n || side > m {
		fmt.Fprintln(out, -1)
		return
	}

	fmt.Fprintln(out, side*side-count)
}
