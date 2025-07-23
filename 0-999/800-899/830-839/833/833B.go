package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf = -int(1e9)

var (
	n, k          int
	a             []int
	freq          []int
	curL, curR    int
	curCost       int
	dpPrev, dpCur []int
	reader        *bufio.Reader
)

func add(x int) {
	freq[x]++
	if freq[x] == 1 {
		curCost++
	}
}

func removeVal(x int) {
	freq[x]--
	if freq[x] == 0 {
		curCost--
	}
}

func costLR(L, R int) int {
	for curL > L {
		curL--
		add(a[curL])
	}
	for curR < R {
		curR++
		add(a[curR])
	}
	for curL < L {
		removeVal(a[curL])
		curL++
	}
	for curR > R {
		removeVal(a[curR])
		curR--
	}
	return curCost
}

func solve(l, r, optL, optR int) {
	if l > r {
		return
	}
	mid := (l + r) >> 1
	bestK := -1
	dpCur[mid] = negInf
	maxK := optR
	if mid-1 < maxK {
		maxK = mid - 1
	}
	for t := optL; t <= maxK; t++ {
		if dpPrev[t] == negInf {
			continue
		}
		val := dpPrev[t] + costLR(t+1, mid)
		if val > dpCur[mid] {
			dpCur[mid] = val
			bestK = t
		}
	}
	if l == r {
		return
	}
	if bestK == -1 {
		bestK = optL
	}
	solve(l, mid-1, optL, bestK)
	solve(mid+1, r, bestK, optR)
}

func readInt() int {
	sign := 1
	val := 0
	c, _ := reader.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = reader.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = reader.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		c = b
	}
	return sign * val
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	n = readInt()
	k = readInt()
	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = readInt()
	}
	freq = make([]int, n+1)
	dpPrev = make([]int, n+1)
	dpCur = make([]int, n+1)
	for i := 1; i <= n; i++ {
		dpPrev[i] = negInf
	}
	dpPrev[0] = 0
	curL, curR, curCost = 1, 0, 0
	for j := 1; j <= k; j++ {
		for i := 0; i <= n; i++ {
			dpCur[i] = negInf
		}
		solve(j, n, 0, n-1)
		dpPrev, dpCur = dpCur, dpPrev
	}
	fmt.Println(dpPrev[n])
}
