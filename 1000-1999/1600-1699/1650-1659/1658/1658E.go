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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	// Positions of values 1..n*n
	// p[v] stores coordinates (r, c)
	// Using 0-based indexing for coordinates in storage, but logic is coordinate-system agnostic
	// as long as we are consistent.
	// The C++ code used 1-based indexing. I'll use 0-based.
	type pt struct{ r, c int }
	p := make([]pt, n*n+1)

	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			var val int
			fmt.Fscan(in, &val)
			grid[i][j] = val
			p[val] = pt{i, j}
		}
	}

	// ans[i][j] defaults to false ('G'), set to true for 'M'
	ans := make([][]bool, n)
	for i := range ans {
		ans[i] = make([]bool, n)
	}

	// Initialize bounds to a very small number
	const negInf = -1000000000
	a, b, c, d := negInf, negInf, negInf, negInf

	max := func(x, y int) int {
		if x > y {
			return x
		}
		return y
	}

	for v := n * n; v >= 1; v-- {
		r, col := p[v].r, p[v].c
		
		// Check if (r, col) is "blocked" by any previous winning cell.
		// Blocked means dist > k.
		// We want to know if dist <= k for ALL previous winning cells.
		// If so, we cannot move to a better cell, so we win.
		
		// The conditions check if (r,c) is within the "safe" box defined by previous winners.
		// Logic from C++:
		// if(a <= x+y+k && b <= x-y+k && c <= -x+y+k && d <= -x-y+k)
		
		// My r, c are 0-based. x, y in C++ were 1-based.
		// The relative difference doesn't change for Manhattan distance.
		// Constants like bounds just shift.
		// So logic holds for 0-based too.
		
		if a <= r+col+k && b <= r-col+k && c <= -r+col+k && d <= -r-col+k {
			ans[r][col] = true // Marin wins
			a = max(a, r+col)
			b = max(b, r-col)
			c = max(c, -r+col)
			d = max(d, -r-col)
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if ans[i][j] {
				fmt.Fprint(out, "M")
			} else {
				fmt.Fprint(out, "G")
			}
		}
		fmt.Fprintln(out)
	}
}