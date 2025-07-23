package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000007

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
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(reader, &line)
		grid[i] = []byte(line)
	}
	if grid[0][0] != grid[n-1][m-1] {
		fmt.Println(0)
		return
	}
	tot := n + m - 2
	steps := tot / 2
	dpCurr := make([][]int, n)
	dpNext := make([][]int, n)
	zeroRow := make([]int, n)
	for i := 0; i < n; i++ {
		dpCurr[i] = make([]int, n)
		dpNext[i] = make([]int, n)
	}
	dpCurr[0][n-1] = 1
	for t := 0; t < steps; t++ {
		for i := 0; i < n; i++ {
			copy(dpNext[i], zeroRow)
		}
		r1min := max(0, t-(m-1))
		r1max := min(n-1, t)
		r2min := max(0, tot-t-(m-1))
		r2max := min(n-1, tot-t)
		for r1 := r1min; r1 <= r1max; r1++ {
			c1 := t - r1
			for r2 := r2min; r2 <= r2max; r2++ {
				val := dpCurr[r1][r2]
				if val == 0 {
					continue
				}
				c2 := tot - t - r2
				if grid[r1][c1] != grid[r2][c2] {
					continue
				}
				// move start down
				if r1+1 < n {
					nr1 := r1 + 1
					nc1 := c1
					if r2-1 >= 0 {
						nr2 := r2 - 1
						nc2 := c2
						if grid[nr1][nc1] == grid[nr2][nc2] {
							dpNext[nr1][nr2] += val
							if dpNext[nr1][nr2] >= MOD {
								dpNext[nr1][nr2] -= MOD
							}
						}
					}
					if c2-1 >= 0 {
						nr2 := r2
						nc2 := c2 - 1
						if grid[nr1][nc1] == grid[nr2][nc2] {
							dpNext[nr1][nr2] += val
							if dpNext[nr1][nr2] >= MOD {
								dpNext[nr1][nr2] -= MOD
							}
						}
					}
				}
				// move start right
				if c1+1 < m {
					nr1 := r1
					nc1 := c1 + 1
					if r2-1 >= 0 {
						nr2 := r2 - 1
						nc2 := c2
						if grid[nr1][nc1] == grid[nr2][nc2] {
							dpNext[nr1][nr2] += val
							if dpNext[nr1][nr2] >= MOD {
								dpNext[nr1][nr2] -= MOD
							}
						}
					}
					if c2-1 >= 0 {
						nr2 := r2
						nc2 := c2 - 1
						if grid[nr1][nc1] == grid[nr2][nc2] {
							dpNext[nr1][nr2] += val
							if dpNext[nr1][nr2] >= MOD {
								dpNext[nr1][nr2] -= MOD
							}
						}
					}
				}
			}
		}
		dpCurr, dpNext = dpNext, dpCurr
	}
	ans := 0
	t := steps
	r1min := max(0, t-(m-1))
	r1max := min(n-1, t)
	r2min := max(0, tot-t-(m-1))
	r2max := min(n-1, tot-t)
	for r1 := r1min; r1 <= r1max; r1++ {
		c1 := t - r1
		for r2 := r2min; r2 <= r2max; r2++ {
			val := dpCurr[r1][r2]
			if val == 0 {
				continue
			}
			c2 := tot - t - r2
			if grid[r1][c1] != grid[r2][c2] {
				continue
			}
			if tot%2 == 0 {
				if r1 == r2 && c1 == c2 {
					ans += val
					if ans >= MOD {
						ans -= MOD
					}
				}
			} else {
				if abs(r1-r2)+abs(c1-c2) == 1 {
					ans += val
					if ans >= MOD {
						ans -= MOD
					}
				}
			}
		}
	}
	fmt.Println(ans)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
