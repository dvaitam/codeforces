package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(reader, &n, &m)
		xs := make([]int64, n)
		ps := make([]int64, n)
		points := make([]int64, 0, 3*n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &xs[i], &ps[i])
			points = append(points, xs[i]-ps[i], xs[i], xs[i]+ps[i])
		}
		sort.Slice(points, func(i, j int) bool { return points[i] < points[j] })
		uniq := make([]int64, 0, len(points))
		for _, v := range points {
			if len(uniq) == 0 || uniq[len(uniq)-1] != v {
				uniq = append(uniq, v)
			}
		}
		idx := make(map[int64]int, len(uniq))
		for i, v := range uniq {
			idx[v] = i
		}
		k := len(uniq)
		ds := make([]int64, k+1)
		dc := make([]int64, k+1)
		for i := 0; i < n; i++ {
			L := idx[xs[i]-ps[i]]
			M := idx[xs[i]]
			R := idx[xs[i]+ps[i]]
			ds[L] += 1
			dc[L] += -(xs[i] - ps[i])
			ds[M] += -2
			dc[M] += (xs[i] - ps[i]) + (xs[i] + ps[i])
			ds[R] += 1
			dc[R] += -(xs[i] + ps[i])
		}
		slope := int64(0)
		con := int64(0)
		val := make([]int64, k)
		for i := 0; i < k; i++ {
			slope += ds[i]
			con += dc[i]
			val[i] = slope*uniq[i] + con
		}
		inf := int64(1) << 60
		pref := make([]int64, k)
		mx := -inf
		for i := 0; i < k; i++ {
			if val[i] > m {
				tmp := val[i] - m - uniq[i]
				if tmp > mx {
					mx = tmp
				}
			}
			pref[i] = mx
		}
		suff := make([]int64, k)
		mx = -inf
		for i := k - 1; i >= 0; i-- {
			if val[i] > m {
				tmp := val[i] - m + uniq[i]
				if tmp > mx {
					mx = tmp
				}
			}
			suff[i] = mx
		}
		ans := make([]byte, n)
		for i := 0; i < n; i++ {
			pos := idx[xs[i]]
			req := int64(-inf)
			if pref[pos] != -inf {
				req = pref[pos] + uniq[pos]
			}
			if suff[pos] != -inf {
				temp := suff[pos] - uniq[pos]
				if temp > req {
					req = temp
				}
			}
			if req <= ps[i] {
				ans[i] = '1'
			} else {
				ans[i] = '0'
			}
		}
		writer.Write(ans)
		writer.WriteByte('\n')
	}
}
