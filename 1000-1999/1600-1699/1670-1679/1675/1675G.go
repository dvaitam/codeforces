package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	pref := make([]int, n)
	s := 0
	for i := 0; i < n; i++ {
		s += a[i]
		pref[i] = s
	}

	const INF int = int(1e9)
	dpPrev := make([][]int, m+1)
	dpCurr := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dpPrev[i] = make([]int, m+1)
		dpCurr[i] = make([]int, m+1)
		for j := 0; j <= m; j++ {
			dpPrev[i][j] = INF
			dpCurr[i][j] = INF
		}
	}

	if n == 1 {
		fmt.Fprintln(out, 0)
		return
	}

	for j := 0; j <= m; j++ {
		if j > m {
			continue
		}
		dpPrev[j][j] = abs(pref[0] - j)
	}

	best := make([]int, m+1)
	for i := 2; i <= n; i++ {
		for s := 0; s <= m; s++ {
			for j := 0; j <= m; j++ {
				dpCurr[s][j] = INF
			}
		}
		for sPrev := 0; sPrev <= m; sPrev++ {
			bestVal := INF
			for j := m; j >= 0; j-- {
				if dpPrev[sPrev][j] < bestVal {
					bestVal = dpPrev[sPrev][j]
				}
				best[j] = bestVal
			}
			for newJ := 0; newJ <= m-sPrev; newJ++ {
				val := best[newJ]
				if val >= INF {
					continue
				}
				newS := sPrev + newJ
				add := 0
				if i < n {
					add = abs(pref[i-1] - newS)
				}
				if dpCurr[newS][newJ] > val+add {
					dpCurr[newS][newJ] = val + add
				}
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}

	ans := INF
	for j := 0; j <= m; j++ {
		if dpPrev[m][j] < ans {
			ans = dpPrev[m][j]
		}
	}
	fmt.Fprintln(out, ans)
}
