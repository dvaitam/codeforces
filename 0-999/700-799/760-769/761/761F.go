package main

import (
	"bufio"
	"fmt"
	"os"
)

func make2DInt(n, m int) [][]int {
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, m)
	}
	return a
}

func make2DInt64(n, m int) [][]int64 {
	a := make([][]int64, n)
	for i := range a {
		a[i] = make([]int64, m)
	}
	return a
}

func query(ps [][]int64, a, b, c, d int) int64 {
	return ps[c][d] - ps[a-1][d] - ps[c][b-1] + ps[a-1][b-1]
}

type Photo struct {
	a, b, c, d int
	e          int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	base := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		base[i] = []byte(s)
	}

	photos := make([]Photo, k)

	diff := make([][][]int, 26)
	for l := 0; l < 26; l++ {
		diff[l] = make2DInt(n+2, m+2)
	}
	cover := make2DInt(n+2, m+2)

	for i := 0; i < k; i++ {
		var a, b, c, d int
		var e string
		fmt.Fscan(in, &a, &b, &c, &d, &e)
		le := int(e[0] - 'a')
		photos[i] = Photo{a, b, c, d, le}

		diff[le][a][b]++
		diff[le][a][d+1]--
		diff[le][c+1][b]--
		diff[le][c+1][d+1]++

		cover[a][b]++
		cover[a][d+1]--
		cover[c+1][b]--
		cover[c+1][d+1]++
	}

	for l := 0; l < 26; l++ {
		for i := 1; i <= n; i++ {
			row := diff[l][i]
			rowPrev := diff[l][i-1]
			for j := 1; j <= m; j++ {
				row[j] += rowPrev[j] + row[j-1] - rowPrev[j-1]
			}
		}
	}
	for i := 1; i <= n; i++ {
		row := cover[i]
		rowPrev := cover[i-1]
		for j := 1; j <= m; j++ {
			row[j] += rowPrev[j] + row[j-1] - rowPrev[j-1]
		}
	}

	contrib := make([][][]int64, 26)
	for l := 0; l < 26; l++ {
		contrib[l] = make2DInt64(n+1, m+1)
	}
	basePrefix := make2DInt64(n+1, m+1)

	counts := make([]int, 26)
	preCnt := make([]int, 26)
	preVal := make([]int, 26)
	contributions := make([]int64, 26)

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			painted := cover[i][j]
			for l := 0; l < 26; l++ {
				counts[l] = diff[l][i][j]
			}
			baseLetter := int(base[i-1][j-1] - 'a')
			counts[baseLetter] += k - painted

			sCnt := 0
			sVal := 0
			for l := 0; l < 26; l++ {
				sCnt += counts[l]
				sVal += counts[l] * l
				preCnt[l] = sCnt
				preVal[l] = sVal
			}
			totalCnt := preCnt[25]
			totalVal := preVal[25]
			for l := 0; l < 26; l++ {
				leftCnt := 0
				leftVal := 0
				if l > 0 {
					leftCnt = preCnt[l-1]
					leftVal = preVal[l-1]
				}
				rightCnt := totalCnt - preCnt[l]
				rightVal := totalVal - preVal[l]
				contributions[l] = int64(l*leftCnt - leftVal + rightVal - l*rightCnt)
			}
			baseC := contributions[baseLetter]
			for l := 0; l < 26; l++ {
				contrib[l][i][j] = contrib[l][i-1][j] + contrib[l][i][j-1] - contrib[l][i-1][j-1] + contributions[l]
			}
			basePrefix[i][j] = basePrefix[i-1][j] + basePrefix[i][j-1] - basePrefix[i-1][j-1] + baseC
		}
	}

	totalBase := basePrefix[n][m]
	minVal := int64(1<<63 - 1)
	for _, p := range photos {
		rectBase := query(basePrefix, p.a, p.b, p.c, p.d)
		rectNew := query(contrib[p.e], p.a, p.b, p.c, p.d)
		val := totalBase - rectBase + rectNew
		if val < minVal {
			minVal = val
		}
	}
	fmt.Fprintln(out, minVal)
}
