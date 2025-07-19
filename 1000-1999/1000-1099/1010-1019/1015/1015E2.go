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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	grid := make([][]byte, n+2)
	for i := 1; i <= n; i++ {
		var line string
		fmt.Fscan(reader, &line)
		// pad to 1-based
		grid[i] = make([]byte, m+2)
		for j := 1; j <= m; j++ {
			grid[i][j] = line[j-1]
		}
	}
	// direction counts
	gore := make([][]int, n+2)
	dole := make([][]int, n+2)
	levo := make([][]int, n+2)
	desno := make([][]int, n+2)
	pox := make([][]int, n+2)
	poy := make([][]int, n+2)
	res := make([][]int, n+2)
	covered := make([][]bool, n+2)
	for i := 0; i <= n+1; i++ {
		gore[i] = make([]int, m+2)
		dole[i] = make([]int, m+2)
		levo[i] = make([]int, m+2)
		desno[i] = make([]int, m+2)
		pox[i] = make([]int, m+2)
		poy[i] = make([]int, m+2)
		res[i] = make([]int, m+2)
		covered[i] = make([]bool, m+2)
	}
	// compute left and up
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i][j] == '*' {
				levo[i][j] = levo[i][j-1] + 1
				gore[i][j] = gore[i-1][j] + 1
			}
		}
	}
	// compute right and down
	for i := n; i >= 1; i-- {
		for j := m; j >= 1; j-- {
			if grid[i][j] == '*' {
				desno[i][j] = desno[i][j+1] + 1
				dole[i][j] = dole[i+1][j] + 1
			}
		}
	}
	// find stars
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i][j] != '*' {
				continue
			}
			d := min(min(gore[i][j], dole[i][j]), min(levo[i][j], desno[i][j]))
			if d > 1 {
				res[i][j] = d - 1
				// mark ranges
				left := max(1, j-d+1)
				right := min(m+1, j+d)
				pox[i][left]++
				pox[i][right]--
				up := max(1, i-d+1)
				down := min(n+1, i+d)
				poy[up][j]++
				poy[down][j]--
			}
		}
	}
	// apply horizontal coverage
	for i := 1; i <= n; i++ {
		acc := 0
		for j := 1; j <= m; j++ {
			acc += pox[i][j]
			if acc > 0 {
				covered[i][j] = true
			}
		}
	}
	// apply vertical coverage
	for j := 1; j <= m; j++ {
		acc := 0
		for i := 1; i <= n; i++ {
			acc += poy[i][j]
			if acc > 0 {
				covered[i][j] = true
			}
		}
	}
	// check all '*' covered
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i][j] == '*' && !covered[i][j] {
				fmt.Fprintln(writer, -1)
				return
			}
		}
	}
	// collect and output
	stars := make([][3]int, 0, n*m)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if res[i][j] > 0 {
				stars = append(stars, [3]int{i, j, res[i][j]})
			}
		}
	}
	fmt.Fprintln(writer, len(stars))
	for _, st := range stars {
		fmt.Fprintf(writer, "%d %d %d\n", st[0], st[1], st[2])
	}
}
