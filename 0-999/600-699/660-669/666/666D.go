package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var perms [][4]int
var corners = [4][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}

func init() {
	used := [4]bool{}
	var p [4]int
	var dfs func(int)
	dfs = func(pos int) {
		if pos == 4 {
			var pp [4]int
			copy(pp[:], p[:])
			perms = append(perms, pp)
			return
		}
		for i := 0; i < 4; i++ {
			if !used[i] {
				used[i] = true
				p[pos] = i
				dfs(pos + 1)
				used[i] = false
			}
		}
	}
	dfs(0)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// uniqueInts sorts and deduplicates slice
func uniqueInts(a []int) []int {
	sort.Ints(a)
	j := 0
	for i := 0; i < len(a); i++ {
		if i == 0 || a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	const INF = 1000000000
	for ; t > 0; t-- {
		var A [4][2]int
		for i := 0; i < 4; i++ {
			fmt.Fscan(reader, &A[i][0], &A[i][1])
		}
		// collect candidate side lengths
		D := make([]int, 0, 12)
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				D = append(D, abs(A[i][0]-A[j][0]), abs(A[i][1]-A[j][1]))
			}
		}
		D = uniqueInts(D)
		ans := INF
		var finalPos [4][2]int

		// try each side length
		for _, d := range D {
			// candidates for x and y
			X := make([]int, 0, 40)
			Y := make([]int, 0, 40)
			for j := 0; j < 4; j++ {
				X = append(X, A[j][0], A[j][0]-d, A[j][0]+d)
				Y = append(Y, A[j][1], A[j][1]-d, A[j][1]+d)
			}
			// midpoints from permutations
			for _, p := range perms {
				lx, rx := INF, -INF
				ly, ry := INF, -INF
				for k := 0; k < 4; k++ {
					xx := A[k][0] - corners[p[k]][0]*d
					yy := A[k][1] - corners[p[k]][1]*d
					if xx < lx {
						lx = xx
					}
					if xx > rx {
						rx = xx
					}
					if yy < ly {
						ly = yy
					}
					if yy > ry {
						ry = yy
					}
				}
				X = append(X, (lx+rx)/2)
				Y = append(Y, (ly+ry)/2)
			}
			X = uniqueInts(X)
			Y = uniqueInts(Y)
			// try each placement
			for _, x := range X {
				for _, y := range Y {
					for _, p := range perms {
						tmax := 0
						ok := true
						var tmp [4][2]int
						for j := 0; j < 4; j++ {
							xx := x + corners[p[j]][0]*d
							yy := y + corners[p[j]][1]*d
							if xx != A[j][0] && yy != A[j][1] {
								ok = false
								break
							}
							move := abs(xx-A[j][0]) + abs(yy-A[j][1])
							if move > tmax {
								tmax = move
							}
							tmp[j][0], tmp[j][1] = xx, yy
						}
						if ok && tmax < ans {
							ans = tmax
							finalPos = tmp
						}
					}
				}
			}
		}
		if ans < INF {
			fmt.Fprintln(writer, ans)
			for i := 0; i < 4; i++ {
				fmt.Fprintln(writer, finalPos[i][0], finalPos[i][1])
			}
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
