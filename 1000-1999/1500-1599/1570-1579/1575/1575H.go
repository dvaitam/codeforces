package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var aStr, bStr string
	fmt.Fscan(reader, &aStr)
	fmt.Fscan(reader, &bStr)

	maxK := n - m + 1

	// build prefix function for b
	b := []byte(bStr)
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && b[i] != b[j] {
			j = pi[j-1]
		}
		if b[i] == b[j] {
			j++
		}
		pi[i] = j
	}

	next := make([][2]int, m+1)
	for j := 0; j <= m; j++ {
		for c := 0; c < 2; c++ {
			ch := byte('0' + c)
			if j < m && b[j] == ch {
				next[j][c] = j + 1
			} else {
				k := j
				for k > 0 && (k == m || b[k] != ch) {
					k = pi[k-1]
				}
				if k < m && b[k] == ch {
					k++
				}
				next[j][c] = k
			}
		}
	}

	const INF int = 1 << 30
	dp := make([][]int, m+1)
	ndp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]int, maxK+1)
		ndp[i] = make([]int, maxK+1)
		for k := 0; k <= maxK; k++ {
			dp[i][k] = INF
		}
	}
	dp[0][0] = 0

	for i := 0; i < n; i++ {
		for j := 0; j <= m; j++ {
			for k := 0; k <= maxK; k++ {
				ndp[j][k] = INF
			}
		}
		for j := 0; j <= m; j++ {
			for k := 0; k <= maxK; k++ {
				val := dp[j][k]
				if val == INF {
					continue
				}
				for c := 0; c < 2; c++ {
					nj := next[j][c]
					nk := k
					if nj == m {
						nk++
					}
					if nk > maxK {
						continue
					}
					cost := val
					if aStr[i] != byte('0'+c) {
						cost++
					}
					if cost < ndp[nj][nk] {
						ndp[nj][nk] = cost
					}
				}
			}
		}
		dp, ndp = ndp, dp
	}

	ans := make([]int, maxK+1)
	for k := 0; k <= maxK; k++ {
		best := INF
		for j := 0; j <= m; j++ {
			if dp[j][k] < best {
				best = dp[j][k]
			}
		}
		if best == INF {
			best = -1
		}
		ans[k] = best
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
