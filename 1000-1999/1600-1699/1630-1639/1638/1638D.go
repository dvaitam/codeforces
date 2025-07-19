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
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}
	vis := make([][]bool, n)
	for i := range vis {
		vis[i] = make([]bool, m)
	}
	type op struct{ r, c, col int }
	var ans []op
	dx := []int{1, 0, -1, 0, 1, 1, -1, -1}
	dy := []int{0, 1, 0, -1, 1, -1, 1, -1}

	// count of matching + zero in 2x2 square or -1 if invalid
	countSq := func(r, c int) int {
		if r < 0 || c < 0 || r+1 >= n || c+1 >= m {
			return -1
		}
		t0 := a[r][c]
		t1 := a[r][c+1]
		t2 := a[r+1][c]
		t3 := a[r+1][c+1]
		t := [4]int{t0, t1, t2, t3}
		cnt, zr := 0, 0
		for _, v := range t {
			if v == 0 {
				zr++
			} else {
				x := 0
				for _, u := range t {
					if u == v {
						x++
					}
				}
				if x > cnt {
					cnt = x
				}
			}
		}
		return cnt + zr
	}

	// BFS queue
	var qr, qc []int
	for i := 0; i < n-1; i++ {
		for j := 0; j < m-1; j++ {
			if countSq(i, j) == 4 {
				vis[i][j] = true
				qr = append(qr, i)
				qc = append(qc, j)
			}
		}
	}
	for head := 0; head < len(qr); head++ {
		r, c := qr[head], qc[head]
		// determine color
		col := 0
		if a[r][c] != 0 {
			col = a[r][c]
		}
		if a[r+1][c] != 0 {
			col = a[r+1][c]
		}
		if a[r][c+1] != 0 {
			col = a[r][c+1]
		}
		if a[r+1][c+1] != 0 {
			col = a[r+1][c+1]
		}
		if col == 0 {
			continue
		}
		ans = append(ans, op{r, c, col})
		// paint to zero
		a[r][c], a[r][c+1], a[r+1][c], a[r+1][c+1] = 0, 0, 0, 0
		// check neighbors
		for k := 0; k < 8; k++ {
			nr, nc := r+dx[k], c+dy[k]
			if nr >= 0 && nc >= 0 && nr < n-1 && nc < m-1 && !vis[nr][nc] && countSq(nr, nc) == 4 {
				vis[nr][nc] = true
				qr = append(qr, nr)
				qc = append(qc, nc)
			}
		}
	}
	// check all zero
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if a[i][j] != 0 {
				fmt.Fprintln(out, -1)
				return
			}
		}
	}
	// reverse ans
	for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
		ans[i], ans[j] = ans[j], ans[i]
	}
	fmt.Fprintln(out, len(ans))
	for _, v := range ans {
		// output 1-based
		fmt.Fprintf(out, "%d %d %d\n", v.r+1, v.c+1, v.col)
	}
}
