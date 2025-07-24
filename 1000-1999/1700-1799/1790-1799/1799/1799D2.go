package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		solve(reader, writer)
	}
}

func solve(r *bufio.Reader, w *bufio.Writer) {
	var n, k int
	fmt.Fscan(r, &n, &k)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	cold := make([]int64, k+1)
	hot := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		fmt.Fscan(r, &cold[i])
	}
	for i := 1; i <= k; i++ {
		fmt.Fscan(r, &hot[i])
	}
	diff := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		diff[i] = cold[i] - hot[i]
	}
	var base int64
	for _, v := range a {
		base += cold[v]
	}
	// dp stores value minus baseAdd
	dp := make(map[int]int64)
	dp[0] = 0
	var bestStored int64
	var baseAdd int64
	prev := a[0]
	for i := 1; i < n; i++ {
		x := a[i]
		add := int64(0)
		if x == prev {
			add = diff[x]
		}
		candidatePrev := bestStored + baseAdd
		if val, ok := dp[x]; ok {
			cand2 := val + baseAdd + diff[x]
			if cand2 > candidatePrev {
				candidatePrev = cand2
			}
		}
		curPrev := int64(math.MinInt64 / 4)
		if v, ok := dp[prev]; ok {
			curPrev = v
		}
		curPrevActual := curPrev + baseAdd + add
		cand := candidatePrev
		if curPrevActual > cand {
			cand = curPrevActual
		}
		baseAdd += add
		dp[prev] = cand - baseAdd
		if dp[prev] > bestStored {
			bestStored = dp[prev]
		}
		prev = x
	}
	ans := base - (bestStored + baseAdd)
	fmt.Fprintln(w, ans)
}
