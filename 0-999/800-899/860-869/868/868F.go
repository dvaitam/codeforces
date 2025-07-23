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

	var n, k int
	fmt.Fscan(in, &n, &k)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	const INF int64 = 1 << 60

	dpPrev := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dpPrev[i] = INF
	}
	cnt := make([]int, n+1)
	curL, curR := 1, 0
	var curCost int64

	add := func(idx int) {
		v := arr[idx]
		curCost += int64(cnt[v])
		cnt[v]++
	}
	remove := func(idx int) {
		v := arr[idx]
		cnt[v]--
		curCost -= int64(cnt[v])
	}
	setSeg := func(l, r int) {
		for curL > l {
			curL--
			add(curL)
		}
		for curR < r {
			curR++
			add(curR)
		}
		for curL < l {
			remove(curL)
			curL++
		}
		for curR > r {
			remove(curR)
			curR--
		}
	}

	dpCur := make([]int64, n+1)
	var solve func(L, R, optL, optR int)
	solve = func(L, R, optL, optR int) {
		if L > R {
			return
		}
		mid := (L + R) / 2
		bestPos := optL
		bestVal := int64(INF)
		start := optL
		end := optR
		if start < 0 {
			start = 0
		}
		if end > mid-1 {
			end = mid - 1
		}
		for i := start; i <= end; i++ {
			setSeg(i+1, mid)
			val := dpPrev[i] + curCost
			if val < bestVal {
				bestVal = val
				bestPos = i
			}
		}
		dpCur[mid] = bestVal
		solve(L, mid-1, optL, bestPos)
		solve(mid+1, R, bestPos, optR)
	}

	dpPrev[0] = 0
	for seg := 1; seg <= k; seg++ {
		for i := 0; i <= n; i++ {
			dpCur[i] = INF
		}
		curL, curR, curCost = 1, 0, 0
		for i := range cnt {
			cnt[i] = 0
		}
		solve(1, n, 0, n-1)
		dpPrev, dpCur = dpCur, dpPrev
	}

	fmt.Fprintln(out, dpPrev[n])
}
