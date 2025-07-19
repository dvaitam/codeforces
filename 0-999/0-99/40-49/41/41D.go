package main

import (
	"fmt"
)

func main() {
	var n, m, k int
	if _, err := fmt.Scan(&n, &m, &k); err != nil {
		return
	}
	K := k + 1
	// mp[i][j]: peas at row i, col j (1-based)
	var mp [101][101]int
	for i := 1; i <= n; i++ {
		var s string
		fmt.Scan(&s)
		for j := 1; j <= m; j++ {
			mp[i][j] = int(s[j-1] - '0')
		}
	}
	// f[i][j][u]: max peas sum from row i..n at (i,j) with total mod K == u
	var f [101][101][12]int
	// dir for reconstruction: 'L' or 'R'
	var dir [101][101][12]byte
	// init to -1
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			for u := 0; u < K; u++ {
				f[i][j][u] = -1
			}
		}
	}
	// base: bottom row i=n
	for j := 1; j <= m; j++ {
		r := mp[n][j] % K
		f[n][j][r] = mp[n][j]
	}
	// dp from bottom-1 up to top
	for i := n - 1; i >= 1; i-- {
		for j := 1; j <= m; j++ {
			for u := 0; u < K; u++ {
				u0 := (u + mp[i][j]) % K
				// move down-left (pawn up-right)
				if j > 1 {
					if f[i+1][j-1][u] >= 0 {
						v := f[i+1][j-1][u] + mp[i][j]
						if v > f[i][j][u0] {
							f[i][j][u0] = v
							dir[i][j][u0] = 'R'
						}
					}
				}
				// move down-right (pawn up-left)
				if j < m {
					if f[i+1][j+1][u] >= 0 {
						v := f[i+1][j+1][u] + mp[i][j]
						if v > f[i][j][u0] {
							f[i][j][u0] = v
							dir[i][j][u0] = 'L'
						}
					}
				}
			}
		}
	}
	// find best at top row i=1 with remainder 0
	ans := -1
	startColAtTop := 1
	for j := 1; j <= m; j++ {
		if f[1][j][0] > ans {
			ans = f[1][j][0]
			startColAtTop = j
		}
	}
	if ans < 0 {
		fmt.Println(-1)
		return
	}
	// reconstruct path
	var path []byte
	var startCol int
	// recursive find from top to bottom
	var find func(i, j, u int)
	find = func(i, j, u int) {
		if i == n {
			startCol = j
			return
		}
		sum := f[i][j][u]
		u0 := (sum - mp[i][j]) % K
		if u0 < 0 {
			u0 += K
		}
		d := dir[i][j][u]
		if d == 'R' {
			find(i+1, j-1, u0)
		} else {
			find(i+1, j+1, u0)
		}
		path = append(path, d)
	}
	find(1, startColAtTop, 0)
	// output
	fmt.Println(ans)
	fmt.Println(startCol)
	fmt.Println(string(path))
}
