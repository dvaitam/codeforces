package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var a, n, m int
	if _, err := fmt.Fscan(reader, &a, &n, &m); err != nil {
		return
	}

	rain := make([]bool, a)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		for x := l; x < r; x++ {
			rain[x] = true
		}
	}

	umbIndex := make([]int, a+1)
	weights := []int64{0}

	for i := 0; i < m; i++ {
		var x int
		var p int64
		fmt.Fscan(reader, &x, &p)
		if umbIndex[x] == 0 {
			weights = append(weights, p)
			umbIndex[x] = len(weights) - 1
		} else if p < weights[umbIndex[x]] {
			weights[umbIndex[x]] = p
		}
	}

	U := len(weights) - 1
	const INF int64 = 1 << 60
	dp := make([][]int64, a+1)
	for i := 0; i <= a; i++ {
		dp[i] = make([]int64, U+1)
		for j := 0; j <= U; j++ {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0
	if umbIndex[0] != 0 {
		dp[0][umbIndex[0]] = 0
	}

	for pos := 0; pos < a; pos++ {
		for j := 0; j <= U; j++ {
			val := dp[pos][j]
			if val == INF {
				continue
			}
			// possible umbrellas to use this step
			candidates := []int{}
			if j != 0 {
				candidates = append(candidates, j)
			}
			if umbIndex[pos] != 0 && umbIndex[pos] != j {
				candidates = append(candidates, umbIndex[pos])
			}
			if !rain[pos] {
				candidates = append(candidates, 0)
			}
			if len(candidates) == 0 && rain[pos] {
				continue
			}
			used := make(map[int]bool)
			for _, k := range candidates {
				if used[k] {
					continue
				}
				used[k] = true
				if k == 0 && rain[pos] {
					continue
				}
				cost := val + weights[k]
				if cost < dp[pos+1][k] {
					dp[pos+1][k] = cost
				}
				if umbIndex[pos+1] != 0 {
					idx := umbIndex[pos+1]
					if cost < dp[pos+1][idx] {
						dp[pos+1][idx] = cost
					}
				}
				if cost < dp[pos+1][0] {
					dp[pos+1][0] = cost
				}
			}
		}
	}

	ans := INF
	for j := 0; j <= U; j++ {
		if dp[a][j] < ans {
			ans = dp[a][j]
		}
	}
	if ans == INF {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}
