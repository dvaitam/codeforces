package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = 0x3f3f3f3f3f3f3f3f

func main() {
	// Faster I/O setup
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, t, m int
	if _, err := fmt.Fscan(reader, &n, &t, &m); err != nil {
		return
	}

	a := make([]int, n+2)
	b := make([]int, n+2)
	c := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i], &b[i], &c[i])
	}

	// Boundary condition for the extra "infinite" item
	n++
	a[n] = 1000000000000000000
	c[n] = 1000000000000000000

	pre := make([][2]int, n+1)
	for i := 1; i <= n; i++ {
		pre[i][0] = pre[i-1][0] + a[i]
		pre[i][1] = pre[i-1][1] + b[i]
	}

	// Initialize DP tables with infinity
	f := make([][][2]int, n+1)
	g := make([][][2]int, n+1)
	for i := range f {
		f[i] = make([][2]int, t+1)
		g[i] = make([][2]int, t+1)
		for j := range f[i] {
			f[i][j][0], f[i][j][1] = inf, inf
			g[i][j][0], g[i][j][1] = inf, inf
		}
	}

	// Base cases
	for i := 0; i <= t; i++ {
		f[0][i][0], f[0][i][1] = 0, 0
		g[0][i][0], g[0][i][1] = 0, 0
	}

	// DP Transitions
	for i := 1; i <= n; i++ {
		for j := 0; j <= t; j++ {
			for p := 0; p < 2; p++ {
				// Transition 1: Current item fits
				if a[i]*p+b[i]*j <= c[i] && f[i-1][j][p] < inf {
					f[i][j][p] = f[i-1][j][p]
					cef := (pre[i-1][0]*p + pre[i-1][1]*j + m - 1) / m
					if cef*m <= pre[i][0]*p+pre[i][1]*j {
						g[i][j][p] = cef
					}
				}

				// Transition 2: Splitting based on k
				for k := 0; k < j; k++ {
					if g[i][k][p] < inf {
						r := pre[i][0]*p + pre[i][1]*k - g[i][k][p]*m
						
						// Calculate ceil((max(0, r + (j-k)*b[i] - c[i])) / m)
						val := r + (j-k)*b[i] - c[i]
						if val < 0 {
							val = 0
						}
						cef := (val + m - 1) / m

						if r-cef*m >= 0 && f[i-1][j-k][0] < inf {
							f[i][j][p] = min(f[i][j][p], cef+g[i][k][p]+f[i-1][j-k][0])
							
							ceff := ((j-k)*pre[i-1][1] + m - 1) / m
							if ceff*m <= (j-k)*pre[i][1]+r-cef*m {
								g[i][j][p] = min(g[i][j][p], cef+ceff+g[i][k][p])
							}
						}
					}
				}
			}
		}
	}

	fmt.Fprintln(writer, f[n][t][1])
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
