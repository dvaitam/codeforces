package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// extract segments of positive numbers and zero block lengths
	zeros := []int{}
	seg := []bool{}
	i := 0
	for i < n && a[i] == 0 {
		i++
	}
	zeros = append(zeros, i)
	for i < n {
		j := i
		has2 := false
		for j < n && a[j] > 0 {
			if a[j] == 2 {
				has2 = true
			}
			j++
		}
		seg = append(seg, has2)
		k := j
		for k < n && a[k] == 0 {
			k++
		}
		zeros = append(zeros, k-j)
		i = k
	}

	if len(seg) == 0 {
		fmt.Fprintln(out, n)
		return
	}
	m := len(seg)
	const INF int = int(1e9)
	dp := make([][2]int, m)
	for i := range dp {
		dp[i][0], dp[i][1] = INF, INF
	}
	for s := 0; s < 2; s++ {
		left := 0
		if seg[0] || s == 0 {
			left = 1
		}
		cost := zeros[0] - left
		if cost < 0 {
			cost = 0
		}
		dp[0][s] = cost
	}
	for idx := 1; idx < m; idx++ {
		ndp0, ndp1 := INF, INF
		for ps := 0; ps < 2; ps++ {
			rightPrev := 0
			if seg[idx-1] || ps == 1 {
				rightPrev = 1
			}
			for cs := 0; cs < 2; cs++ {
				leftCur := 0
				if seg[idx] || cs == 0 {
					leftCur = 1
				}
				c := zeros[idx] - rightPrev - leftCur
				if c < 0 {
					c = 0
				}
				val := dp[idx-1][ps] + c
				if cs == 0 {
					if val < ndp0 {
						ndp0 = val
					}
				} else {
					if val < ndp1 {
						ndp1 = val
					}
				}
			}
		}
		dp[idx][0] = ndp0
		dp[idx][1] = ndp1
	}
	res := INF
	for s := 0; s < 2; s++ {
		right := 0
		if seg[m-1] || s == 1 {
			right = 1
		}
		c := zeros[m] - right
		if c < 0 {
			c = 0
		}
		val := dp[m-1][s] + c
		if val < res {
			res = val
		}
	}
	ans := m + res
	fmt.Fprintln(out, ans)
}
