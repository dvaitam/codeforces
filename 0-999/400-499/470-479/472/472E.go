package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MAXN = 32

var n, m int
var sArr [MAXN][MAXN]int
var tArr [MAXN][MAXN]int
var fixArr [MAXN][MAXN]bool
var used [MAXN][MAXN]int
var prvX [MAXN][MAXN]int
var prvY [MAXN][MAXN]int
var iterID int

// pair represents a position (row, col)
type pair struct{ r, c int }

var res []pair
var curx, cury int

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

// solve1d tries to transform s to t by adjacent swaps
func solve1d(s, t []int) []int {
	n := len(s)
	tmp1 := make([]int, n)
	tmp2 := make([]int, n)
	copy(tmp1, s)
	copy(tmp2, t)
	sort.Ints(tmp1)
	sort.Ints(tmp2)
	for i := 0; i < n; i++ {
		if tmp1[i] != tmp2[i] {
			return []int{-1}
		}
	}
	for from := 0; from < n; from++ {
		for to := 0; to < n; to++ {
			if from == to {
				continue
			}
			a := make([]int, n)
			copy(a, s)
			step := sign(to - from)
			for i := from; i != to; i += step {
				a[i], a[i+step] = a[i+step], a[i]
			}
			ok := true
			for i := 0; i < n; i++ {
				if a[i] != t[i] {
					ok = false
					break
				}
			}
			if ok {
				var seq []int
				for i := from; i != to; i += step {
					seq = append(seq, i)
				}
				seq = append(seq, to)
				return seq
			}
		}
	}
	return []int{-1}
}

// movePos swaps current position with (nx,ny)
func movePos(nx, ny int) {
	// swap
	sArr[curx][cury], sArr[nx][ny] = sArr[nx][ny], sArr[curx][cury]
	res = append(res, pair{nx, ny})
	curx, cury = nx, ny
}

// goFar moves cursor to (tx,ty) via BFS avoiding fixed cells
func goFar(tx, ty int) {
	iterID++
	queue := []pair{{tx, ty}}
	used[tx][ty] = iterID
	for used[curx][cury] != iterID {
		p := queue[0]
		queue = queue[1:]
		x, y := p.r, p.c
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				nx, ny := x+dx, y+dy
				if !fixArr[nx][ny] && used[nx][ny] != iterID {
					used[nx][ny] = iterID
					prvX[nx][ny] = x
					prvY[nx][ny] = y
					queue = append(queue, pair{nx, ny})
				}
			}
		}
	}
	for curx != tx || cury != ty {
		px := prvX[curx][cury]
		py := prvY[curx][cury]
		movePos(px, py)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &n, &m)
	// 1D case
	if n == 1 || m == 1 {
		total := n * m
		s1 := make([]int, total)
		t1 := make([]int, total)
		for i := 0; i < total; i++ {
			fmt.Fscan(reader, &s1[i])
		}
		for i := 0; i < total; i++ {
			fmt.Fscan(reader, &t1[i])
		}
		cr := solve1d(s1, t1)
		if len(cr) == 1 && cr[0] == -1 {
			fmt.Fprintln(writer, -1)
			return
		}
		// record moves
		for _, i := range cr {
			if n == 1 {
				res = append(res, pair{1, i + 1})
			} else {
				res = append(res, pair{i + 1, 1})
			}
		}
		// output
		fmt.Fprintln(writer, len(res))
		for _, p := range res {
			fmt.Fprintf(writer, "%d %d\n", p.r, p.c)
		}
		return
	}
	// 2D case
	// read arrays with 1-based indexing
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			fmt.Fscan(reader, &sArr[i][j])
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			fmt.Fscan(reader, &tArr[i][j])
		}
	}
	// set borders fixed
	for j := 0; j <= m+1; j++ {
		fixArr[0][j] = true
		fixArr[n+1][j] = true
	}
	for i := 0; i <= n+1; i++ {
		fixArr[i][0] = true
		fixArr[i][m+1] = true
	}
	// check multisets
	cnt := make([]int, 901)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cnt[sArr[i][j]]++
			cnt[tArr[i][j]]--
		}
	}
	ok := true
	for v := 1; v <= 900; v++ {
		if cnt[v] != 0 {
			ok = false
			break
		}
	}
	if !ok {
		fmt.Fprintln(writer, -1)
		return
	}
	// start from position of target bottom-right value
	target := tArr[n][m]
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if sArr[i][j] == target {
				curx, cury = i, j
			}
		}
	}
	res = append(res, pair{curx, cury})
	// process rows 1..n-2
	for i := 1; i <= n-2; i++ {
		for j := 1; j <= m; j++ {
			// find position with value tArr[i][j]
			var tx, ty int
			for k := 1; k <= n; k++ {
				for l := 1; l <= m; l++ {
					if !fixArr[k][l] && !(k == curx && l == cury) && sArr[k][l] == tArr[i][j] && tx == 0 {
						tx, ty = k, l
					}
				}
			}
			// bring to column j
			for ty != j {
				fixArr[tx][ty] = true
				goFar(tx, ty+sign(j-ty))
				fixArr[tx][ty] = false
				movePos(tx, ty)
				ty += sign(j - ty)
			}
			// bring to row i
			for tx != i {
				fixArr[tx][ty] = true
				goFar(tx+sign(i-tx), ty)
				fixArr[tx][ty] = false
				movePos(tx, ty)
				tx += sign(i - tx)
			}
			fixArr[i][j] = true
		}
	}
	// process last two rows
	for j := 1; j <= m; j++ {
		for _, i2 := range []int{n - 1, n} {
			if i2 == n && j == m {
				continue
			}
			var tx, ty int
			for k := 1; k <= n; k++ {
				for l := 1; l <= m; l++ {
					if !fixArr[k][l] && !(k == curx && l == cury) && sArr[k][l] == tArr[i2][j] && tx == 0 {
						tx, ty = k, l
					}
				}
			}
			for tx != i2 {
				fixArr[tx][ty] = true
				goFar(tx+sign(i2-tx), ty)
				fixArr[tx][ty] = false
				movePos(tx, ty)
				tx += sign(i2 - tx)
			}
			for ty != j {
				fixArr[tx][ty] = true
				goFar(tx, ty+sign(j-ty))
				fixArr[tx][ty] = false
				movePos(tx, ty)
				ty += sign(j - ty)
			}
			fixArr[i2][j] = true
		}
	}
	// output
	// number of moves = len(res)-1
	fmt.Fprintln(writer, len(res)-1)
	for _, p := range res {
		fmt.Fprintf(writer, "%d %d\n", p.r, p.c)
	}
}
