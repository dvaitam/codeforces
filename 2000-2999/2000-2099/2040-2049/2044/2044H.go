package main

import (
	"bufio"
	"fmt"
	"os"
)

func rectSum(pref [][]int64, x1, y1, x2, y2 int) int64 {
	res := pref[x2][y2] - pref[x1-1][y2] - pref[x2][y1-1] + pref[x1-1][y1-1]
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		pref00 := make([][]int64, n+1)
		pref10 := make([][]int64, n+1)
		pref01 := make([][]int64, n+1)
		for i := 0; i <= n; i++ {
			pref00[i] = make([]int64, n+1)
			pref10[i] = make([]int64, n+1)
			pref01[i] = make([]int64, n+1)
		}

		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				var val int64
				fmt.Fscan(reader, &val)
				pref00[i][j] = pref00[i-1][j] + pref00[i][j-1] - pref00[i-1][j-1] + val
				pref10[i][j] = pref10[i-1][j] + pref10[i][j-1] - pref10[i-1][j-1] + int64(i)*val
				pref01[i][j] = pref01[i-1][j] + pref01[i][j-1] - pref01[i-1][j-1] + int64(j)*val
			}
		}

		for qi := 0; qi < q; qi++ {
			var x1, y1, x2, y2 int
			fmt.Fscan(reader, &x1, &y1, &x2, &y2)
			sum00 := rectSum(pref00, x1, y1, x2, y2)
			sum10 := rectSum(pref10, x1, y1, x2, y2)
			sum01 := rectSum(pref01, x1, y1, x2, y2)
			width := int64(y2 - y1 + 1)
			res := width*(sum10-int64(x1)*sum00) + (sum01 - int64(y1)*sum00) + sum00
			if qi > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, res)
		}
		fmt.Fprintln(writer)
	}
}
