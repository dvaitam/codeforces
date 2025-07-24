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
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &row[j])
		}
		matrix[i] = row
	}

	total := 0
	nm := n * m
	for col := 0; col < m; col++ {
		cnt := make([]int, n)
		for row := 0; row < n; row++ {
			val := matrix[row][col]
			if val >= 1 && val <= nm && (val-1)%m == col {
				targetRow := (val - 1) / m
				shift := (row - targetRow + n) % n
				cnt[shift]++
			}
		}
		best := n + 1
		for shift := 0; shift < n; shift++ {
			moves := shift + (n - cnt[shift])
			if moves < best {
				best = moves
			}
		}
		total += best
	}

	fmt.Fprintln(out, total)
}
