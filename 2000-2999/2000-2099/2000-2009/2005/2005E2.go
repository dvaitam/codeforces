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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, n, m int
		fmt.Fscan(in, &l, &n, &m)

		a := make([]int, l)
		for i := 0; i < l; i++ {
			fmt.Fscan(in, &a[i])
		}

		b := make([][]int, n)
		for i := 0; i < n; i++ {
			b[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &b[i][j])
			}
		}

		if l == 0 {
			fmt.Fprintln(out, "N")
			continue
		}

		// hasWin[r][c] = true if in submatrix [r..n-1][c..m-1] there exists a WINNING cell
		// for the NEXT index (idx+1). Size (n+1) x (m+1) to handle out-of-bounds.
		// Initialize for after idx=l-1 (no next index exists).
		hasWin := make([][]bool, n+1)
		for i := range hasWin {
			hasWin[i] = make([]bool, m+1)
		}

		for idx := l - 1; idx >= 0; idx-- {
			// For each cell matching a[idx], compute win(idx, r, c):
			//   win = true if idx==l-1, OR submatrix(r+1,c+1) has no winning cell for idx+1
			//   win = (idx == l-1) || (r+1>=n || c+1>=m || !hasWin[r+1][c+1])
			// Then build new hasWin = suffix-OR of winning cells for idx.

			// First pass: mark winning cells
			winCell := make([][]bool, n)
			for r := 0; r < n; r++ {
				winCell[r] = make([]bool, m)
				for c := 0; c < m; c++ {
					if b[r][c] == a[idx] {
						if idx == l-1 {
							winCell[r][c] = true
						} else if r+1 >= n || c+1 >= m {
							winCell[r][c] = true
						} else {
							winCell[r][c] = !hasWin[r+1][c+1]
						}
					}
				}
			}

			// Build suffix-OR: hasWin[r][c] = OR of winCell[r'][c'] for r'>=r, c'>=c
			newHW := make([][]bool, n+1)
			for i := range newHW {
				newHW[i] = make([]bool, m+1)
			}
			for r := n - 1; r >= 0; r-- {
				for c := m - 1; c >= 0; c-- {
					newHW[r][c] = winCell[r][c] || newHW[r+1][c] || newHW[r][c+1]
					// newHW[r+1][c+1] is already covered by both newHW[r+1][c] and newHW[r][c+1]
					// but OR is idempotent so no issue
				}
			}

			hasWin = newHW
		}

		if hasWin[0][0] {
			fmt.Fprintln(out, "T")
		} else {
			fmt.Fprintln(out, "N")
		}
	}
}
