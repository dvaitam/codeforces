package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	x := make([][]int, n)
	y := make([][]int, n)
	for i := 0; i < n; i++ {
		x[i] = make([]int, n)
		y[i] = make([]int, n)
		for j := 0; j < n; j++ {
			var xi, yi int
			fmt.Fscan(reader, &xi, &yi)
			if xi > 0 {
				xi--
			}
			if yi > 0 {
				yi--
			}
			x[i][j] = xi
			y[i][j] = yi
		}
	}
	ans := make([][]byte, n)
	for i := range ans {
		ans[i] = make([]byte, n)
	}
	// directions: R, D, L, U
	dx := []int{0, 1, 0, -1}
	dy := []int{1, 0, -1, 0}
	pv := []byte{'R', 'D', 'L', 'U'}
	dp := []byte{'L', 'U', 'R', 'D'}

	// associate -1 regions
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if x[i][j] == -1 {
				for t := 0; t < 4; t++ {
					ni, nj := i+dx[t], j+dy[t]
					if ni >= 0 && ni < n && nj >= 0 && nj < n && x[ni][nj] == -1 {
						ans[i][j] = pv[t]
						break
					}
				}
			}
		}
	}
	// BFS from roots
	type pair struct{ i, j int }
	var q []pair
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if x[i][j] == i && y[i][j] == j {
				ans[i][j] = 'X'
				q = append(q, pair{i, j})
				// process queue
				for head := 0; head < len(q); head++ {
					u := q[head]
					for t := 0; t < 4; t++ {
						ni, nj := u.i+dx[t], u.j+dy[t]
						if ni < 0 || ni >= n || nj < 0 || nj >= n {
							continue
						}
						if ans[ni][nj] != 0 {
							continue
						}
						if x[ni][nj] == i && y[ni][nj] == j {
							ans[ni][nj] = dp[t]
							q = append(q, pair{ni, nj})
						}
					}
				}
				// reset queue for next root
				q = q[:0]
			}
		}
	}
	// validate
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if ans[i][j] == 0 {
				fmt.Fprintln(writer, "INVALID")
				return
			}
		}
	}
	fmt.Fprintln(writer, "VALID")
	for i := 0; i < n; i++ {
		writer.Write(ans[i])
		writer.WriteByte('\n')
	}
}
