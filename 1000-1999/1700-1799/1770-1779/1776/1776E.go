package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Interval struct{ l, r int64 }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	var s, v int64
	if _, err := fmt.Fscan(reader, &n, &m, &s, &v); err != nil {
		return
	}

	rails := make([][]Interval, m)
	for i := 0; i < n; i++ {
		var a, b int64
		var r int
		fmt.Fscan(reader, &a, &b, &r)
		r--
		rails[r] = append(rails[r], Interval{a, b})
	}
	for i := 0; i < m; i++ {
		sort.Slice(rails[i], func(x, y int) bool { return rails[i][x].l < rails[i][y].l })
	}

	// candidate times per rail (0..m inclusive, rail 0 is start)
	times := make([][]int64, m+1)
	times[0] = []int64{0}
	for i := 0; i < m; i++ {
		mp := make(map[int64]struct{})
		mp[0] = struct{}{}
		mp[s] = struct{}{}
		for _, iv := range rails[i] {
			if iv.l <= s {
				mp[iv.l] = struct{}{}
			}
			if iv.r <= s {
				mp[iv.r] = struct{}{}
			}
		}
		arr := make([]int64, 0, len(mp))
		for k := range mp {
			arr = append(arr, k)
		}
		sort.Slice(arr, func(x, y int) bool { return arr[x] < arr[y] })
		// filter allowed times
		allowed := make([]int64, 0, len(arr))
		for _, t := range arr {
			if allowedTime(t, rails[i]) {
				allowed = append(allowed, t)
			}
		}
		times[i+1] = allowed
	}

	const INF = int(1e9)

	if m == 1 {
		ans := INF
		for _, t1 := range times[1] {
			if t1 < v {
				continue
			}
			if t1+t1 <= s {
				if 0 < ans {
					ans = 0
				}
			}
			if t1+v <= s {
				if 1 < ans {
					ans = 1
				}
			}
		}
		if ans == INF {
			ans = -1
		}
		fmt.Fprintln(writer, ans)
		return
	}

	// initialize dp for first two rails
	if len(times[1]) == 0 || len(times[2]) == 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	dp := make([][]int, len(times[1]))
	for i := range dp {
		dp[i] = make([]int, len(times[2]))
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	for i1, t1 := range times[1] {
		if t1 < v {
			continue
		}
		for i2, t2 := range times[2] {
			if t2-t1 < v {
				continue
			}
			cost := 1
			if t1 >= t2-t1 {
				cost = 0
			}
			if cost < dp[i1][i2] {
				dp[i1][i2] = cost
			}
		}
	}

	// iterate for rails 3..m
	for r := 3; r <= m; r++ {
		if len(times[r]) == 0 {
			fmt.Fprintln(writer, -1)
			return
		}
		newdp := make([][]int, len(times[r-1]))
		for i := range newdp {
			newdp[i] = make([]int, len(times[r]))
			for j := range newdp[i] {
				newdp[i][j] = INF
			}
		}
		for iPrev := 0; iPrev < len(times[r-1]); iPrev++ {
			for iPrevPrev := 0; iPrevPrev < len(times[r-2]); iPrevPrev++ {
				cost := dp[iPrevPrev][iPrev]
				if cost == INF {
					continue
				}
				tPrevPrev := times[r-2][iPrevPrev]
				tPrev := times[r-1][iPrev]
				deltaPrev := tPrev - tPrevPrev
				for iCurr, tCurr := range times[r] {
					if tCurr-tPrev < v {
						continue
					}
					deltaCurr := tCurr - tPrev
					newCost := cost
					if deltaCurr != deltaPrev {
						newCost++
					}
					if newCost < newdp[iPrev][iCurr] {
						newdp[iPrev][iCurr] = newCost
					}
				}
			}
		}
		dp = newdp
	}

	ans := INF
	for iPrev := 0; iPrev < len(times[m-1]); iPrev++ {
		for iCurr := 0; iCurr < len(times[m]); iCurr++ {
			cost := dp[iPrev][iCurr]
			if cost == INF {
				continue
			}
			tPrev := times[m-1][iPrev]
			tCurr := times[m][iCurr]
			delta := tCurr - tPrev
			if delta >= v && tCurr+delta <= s {
				if cost < ans {
					ans = cost
				}
			}
			if tCurr+v <= s {
				if cost+1 < ans {
					ans = cost + 1
				}
			}
		}
	}
	if ans == INF {
		ans = -1
	}
	fmt.Fprintln(writer, ans)
}

func allowedTime(t int64, trains []Interval) bool {
	for _, iv := range trains {
		if iv.l < t && t < iv.r {
			return false
		}
	}
	return true
}
