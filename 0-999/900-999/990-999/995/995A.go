package main

import (
	"fmt"
)

type move struct {
	id, r, c int
}

func main() {
	var n, k int
	if fmt.Scan(&n, &k); n == 0 && k == 0 {
		// no input
	}
	a := make([][]int, 4)
	for i := 0; i < 4; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Scan(&a[i][j])
			a[i][j]--
		}
	}
	var ans []move
	// direct removals
	for i := 0; i < n; i++ {
		if a[0][i] != -1 && a[0][i] == a[1][i] {
			ans = append(ans, move{a[0][i], 0, i})
			a[0][i], a[1][i] = -1, -1
			k--
		}
	}
	for i := 0; i < n; i++ {
		if a[2][i] != -1 && a[2][i] == a[3][i] {
			ans = append(ans, move{a[2][i], 3, i})
			a[2][i], a[3][i] = -1, -1
			k--
		}
	}
	// check at least one empty in inner rows
	ok := false
	for i := 0; i < n; i++ {
		if a[1][i] == -1 || a[2][i] == -1 {
			ok = true
			break
		}
	}
	if !ok {
		fmt.Println(-1)
		return
	}
	// build cycle of positions
	var vv [][2]int
	for i := 0; i < n; i++ {
		vv = append(vv, [2]int{1, i})
	}
	for i := n - 1; i >= 0; i-- {
		vv = append(vv, [2]int{2, i})
	}
	// repeat start to simplify movement
	vv = append(vv, vv[0])
	// simulate
	for k > 0 {
		for i := 1; i < len(vv) && k > 0; i++ {
			x, y := vv[i][0], vv[i][1]
			if a[x][y] == -1 {
				continue
			}
			// if can remove directly
			nx := x ^ 1
			if a[nx][y] == a[x][y] {
				ans = append(ans, move{a[x][y], nx, y})
				a[x][y], a[nx][y] = -1, -1
				k--
			} else {
				// move to previous if empty
				px, py := vv[i-1][0], vv[i-1][1]
				if a[px][py] == -1 {
					ans = append(ans, move{a[x][y], px, py})
					a[px][py] = a[x][y]
					a[x][y] = -1
				}
			}
		}
	}
	// output
	fmt.Println(len(ans))
	for _, mv := range ans {
		fmt.Printf("%d %d %d\n", mv.id+1, mv.r+1, mv.c+1)
	}
}
