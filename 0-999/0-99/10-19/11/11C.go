package main

import (
	"bufio"
	"fmt"
	"os"
)

var dr = [8]int{0, 1, 0, -1, 1, 1, -1, -1}
var dc = [8]int{1, 0, -1, 0, 1, -1, -1, 1}
var n, m int
var a [][]byte

// chk checks for a square side starting at (r,c) in direction d.
// Returns the side length if a valid frame side is found, otherwise 0.
func chk(r, c, d int) int {
	d1 := (d + 1) % 4
	d2 := (d + 3) % 4
	if d >= 4 {
		d1 += 4
		d2 += 4
	}
	for i := 1; ; i++ {
		r += dr[d]
		c += dc[d]
		if a[r][c] == '0' {
			return 0
		}
		if d >= 4 {
			// no foreign 1s adjacent for diagonal squares
			for j := 0; j < 4; j++ {
				if a[r+dr[j]][c+dc[j]] == '1' {
					return 0
				}
			}
		}
		// check forbidden overlap with next side
		nr := r + dr[d2]
		nc := c + dc[d2]
		if a[nr][nc] == '1' || (d < 4 && a[r+dr[d2]+dr[d]][c+dc[d2]+dc[d]] == '1') {
			return 0
		}
		// corner detection: next side starts
		if a[r+dr[d1]][c+dc[d1]] == '1' {
			// ensure no forward overlap
			if a[r+dr[d]][c+dc[d]] == '1' {
				return 0
			}
			return i
		}
		// ensure continuous side
		if a[r+dr[d]][c+dc[d]] == '0' {
			return 0
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		fmt.Fscan(in, &n, &m)
		// initialize with border of zeros
		a = make([][]byte, n+2)
		for i := range a {
			a[i] = make([]byte, m+2)
			for j := range a[i] {
				a[i][j] = '0'
			}
		}
		// read matrix
		for i := 1; i <= n; i++ {
			var s string
			fmt.Fscan(in, &s)
			for j := 1; j <= m; j++ {
				a[i][j] = s[j-1]
			}
		}
		var res int
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				// axis-aligned squares
				if a[i-1][j-1] == '0' && a[i-1][j] == '0' && a[i-1][j+1] == '0' &&
					a[i][j-1] == '0' && a[i][j] == '1' && a[i][j+1] == '1' &&
					a[i+1][j-1] == '0' && a[i+1][j] == '1' {
					t0 := chk(i, j, 0)
					if t0 > 0 && chk(i, j+t0, 1) == t0 && chk(i+t0, j+t0, 2) == t0 && chk(i+t0, j, 3) == t0 {
						res++
					}
				}
				// diagonal squares
				if a[i-1][j-1] == '0' && a[i-1][j] == '0' && a[i-1][j+1] == '0' &&
					a[i][j-1] == '0' && a[i][j] == '1' && a[i][j+1] == '0' &&
					a[i+1][j-1] == '1' && a[i+1][j] == '0' && a[i+1][j+1] == '1' {
					t0 := chk(i, j, 4)
					if t0 > 0 && chk(i+t0, j+t0, 5) == t0 && chk(i+2*t0, j, 6) == t0 && chk(i+t0, j-t0, 7) == t0 {
						res++
					}
				}
			}
		}
		fmt.Fprintln(out, res)
	}
}
