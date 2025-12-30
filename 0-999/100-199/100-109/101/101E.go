package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 20005

var (
	n, m, mod int
	x, y      [N]int
	pos       [N]int
	f, g      [N]int
	ans       int
)

// calc computes the weight of cell (i, j)
func calc(i, j int) int {
	r := x[i] + y[j]
	if r >= mod {
		r -= mod
	}
	return r
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// solve uses divide and conquer to find the optimal path
// sx, sy: start x, start y
// tx, ty: target x, target y
func solve(sx, sy, tx, ty int) {
	pos[tx] = ty
	if sx == tx {
		return
	}
	mid := (sx + tx) >> 1

	// Clear f and g arrays for the current column range
	for k := sy; k <= ty; k++ {
		f[k] = 0
		g[k] = 0
	}

	// Calculate DP for the top half (sx to mid)
	for i := sx; i <= mid; i++ {
		for j := sy; j <= ty; j++ {
			if j != sy {
				f[j] = max(f[j], f[j-1])
			}
			f[j] += calc(i, j)
		}
	}

	// Calculate DP for the bottom half (tx down to mid+1)
	for i := tx; i > mid; i-- {
		for j := ty; j >= sy; j-- {
			if j != ty {
				g[j] = max(g[j], g[j+1])
			}
			g[j] += calc(i, j)
		}
	}

	// Find the best split column 'bestPos' in the middle row
	bestPos := sy
	maxVal := -1 // Initialize with a value lower than any possible sum (sums are >= 0)

	for i := sy; i <= ty; i++ {
		sum := f[i] + g[i]
		if sum > maxVal {
			maxVal = sum
			bestPos = i
		}
	}

	// Update global answer
	if maxVal > ans {
		ans = maxVal
	}

	// Recursively solve for the two halves
	solve(sx, sy, mid, bestPos)
	solve(mid+1, bestPos, tx, ty)
}

func main() {
	// Set up fast I/O
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Read input
	fmt.Fscan(reader, &n, &m, &mod)

	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &x[i])
		x[i] %= mod
	}
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &y[i])
		y[i] %= mod
	}

	// Special case for a single row
	if n == 1 {
		currentAns := 0
		for i := 1; i <= m; i++ {
			currentAns += (x[1] + y[i]) % mod
		}
		fmt.Fprintln(writer, currentAns)
		for i := 1; i < m; i++ {
			writer.WriteByte('S')
		}
		writer.WriteByte('\n')
		return
	}

	// Solve for the full grid
	solve(1, 1, n, m)

	// Output the maximum weight
	fmt.Fprintln(writer, ans)

	// Reconstruct and output path based on 'pos' array
	// pos[i] indicates the column index where we move down from row i to i+1
	j := 1
	for i := 1; i <= n; i++ {
		for j < pos[i] {
			writer.WriteByte('S')
			j++
		}
		if i != n {
			writer.WriteByte('C')
		}
	}
	writer.WriteByte('\n')
}
