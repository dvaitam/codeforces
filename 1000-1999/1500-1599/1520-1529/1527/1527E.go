package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

var (
	n, k       int
	a          []int
	prevPos    []int
	nextPos    []int
	dpPrev     []int64
	dpCur      []int64
	curL, curR int
	curCost    int64
)

func costLR(L, R int) int64 {
	for curL > L {
		curL--
		v := nextPos[curL]
		if v <= curR {
			curCost += int64(v - curL)
		}
	}
	for curR < R {
		curR++
		p := prevPos[curR]
		if p >= curL {
			curCost += int64(curR - p)
		}
	}
	for curL < L {
		v := nextPos[curL]
		if v <= curR {
			curCost -= int64(v - curL)
		}
		curL++
	}
	for curR > R {
		p := prevPos[curR]
		if p >= curL {
			curCost -= int64(curR - p)
		}
		curR--
	}
	return curCost
}

func solve(l, r, optL, optR int) {
	if l > r {
		return
	}
	mid := (l + r) >> 1
	best := -1
	dpCur[mid] = inf
	mx := optR
	if mid-1 < mx {
		mx = mid - 1
	}
	for j := optL; j <= mx; j++ {
		val := dpPrev[j] + costLR(j+1, mid)
		if val < dpCur[mid] {
			dpCur[mid] = val
			best = j
		}
	}
	if l == r {
		return
	}
	if best == -1 {
		best = optL
	}
	solve(l, mid-1, optL, best)
	solve(mid+1, r, best, optR)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	prevPos = make([]int, n+1)
	last := make([]int, n+1)
	for i := 1; i <= n; i++ {
		v := a[i]
		prevPos[i] = last[v]
		last[v] = i
	}
	for i := 0; i <= n; i++ {
		last[i] = n + 1
	}
	nextPos = make([]int, n+1)
	for i := n; i >= 1; i-- {
		v := a[i]
		nextPos[i] = last[v]
		last[v] = i
	}
	dpPrev = make([]int64, n+1)
	dpCur = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dpPrev[i] = inf
	}
	dpPrev[0] = 0
	curL, curR, curCost = 1, 0, 0

	for seg := 1; seg <= k; seg++ {
		for i := 0; i <= n; i++ {
			dpCur[i] = inf
		}
		solve(seg, n, seg-1, n-1)
		dpPrev, dpCur = dpCur, dpPrev
	}

	fmt.Fprintln(out, dpPrev[n])
}
