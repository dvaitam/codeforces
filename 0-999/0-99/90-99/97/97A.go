package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, m  int
	a     []string
	b     [][]int
	ans   [][]int
	isans bool
	res   uint64
	dx    = [4]int{0, 0, 1, -1}
	dy    = [4]int{1, -1, 0, 0}
)

func rec(maxc, pos int, mask uint64) {
	// skip assigned or empty
	for pos < n*m {
		x, y := pos/m, pos%m
		if a[x][y] != '.' && b[x][y] == -1 {
			break
		}
		pos++
	}
	if pos == n*m {
		res++
		if !isans {
			// copy b to ans
			for i := 0; i < n; i++ {
				for j := 0; j < m; j++ {
					ans[i][j] = b[i][j]
				}
			}
			isans = true
		}
		return
	}
	x, y := pos/m, pos%m
	// try colors
	for clr := 0; clr <= max(maxc+1, 6); clr++ {
		if clr > 6 {
			break
		}
		// check 2x2 block
		if x+1 >= n || y+1 >= m {
			continue
		}
		ok := true
		for i := 0; i < 2 && ok; i++ {
			for j := 0; j < 2; j++ {
				if a[x+i][y+j] == '.' || b[x+i][y+j] != -1 {
					ok = false
					break
				}
			}
		}
		if !ok {
			continue
		}
		// assign
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				b[x+i][y+j] = clr
			}
		}
		// check adjacency for domino uniqueness
		nmask := mask
		conflict := false
		for i := 0; i < 2 && !conflict; i++ {
			for j := 0; j < 2 && !conflict; j++ {
				xx, yy := x+i, y+j
				for dir := 0; dir < 4; dir++ {
					xxx, yyy := xx+dx[dir], yy+dy[dir]
					if xxx >= 0 && xxx < n && yyy >= 0 && yyy < m && a[xxx][yyy] != '.' && b[xxx][yyy] != -1 {
						// skip within same block twice
						if xxx >= x && xxx <= x+1 && yyy >= y && yyy <= y+1 {
							if xx > xxx || (xx == xxx && yy > yyy) {
								continue
							}
						}
						c1 := b[xx][yy]*7 + b[xxx][yyy]
						c2 := b[xx][yy] + 7*b[xxx][yyy]
						if (nmask>>c1)&1 == 1 {
							conflict = true
							break
						}
						nmask |= 1 << c1
						nmask |= 1 << c2
					}
				}
			}
		}
		if !conflict {
			rec(max(maxc, clr), pos+1, nmask)
		}
		// unassign
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				b[x+i][y+j] = -1
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &n, &m)
	a = make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b = make([][]int, n)
	ans = make([][]int, n)
	for i := 0; i < n; i++ {
		b[i] = make([]int, m)
		ans[i] = make([]int, m)
		for j := 0; j < m; j++ {
			b[i][j] = -1
			ans[i][j] = -1
		}
	}
	isans = false
	res = 0
	rec(-1, 0, 0)
	// multiply by factorials 1..7
	for i := 1; i <= 7; i++ {
		res *= uint64(i)
	}
	fmt.Fprintln(writer, res)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if a[i][j] == '.' {
				writer.WriteByte('.')
			} else {
				writer.WriteByte(byte('0' + ans[i][j]))
			}
		}
		writer.WriteByte('\n')
	}
}
