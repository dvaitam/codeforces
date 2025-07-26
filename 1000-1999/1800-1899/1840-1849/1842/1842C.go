package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxBallsRemoved(a []int) int {
	n := len(a)
	best := make(map[int]int)
	dp := 0
	for i := 1; i <= n; i++ {
		val := a[i-1]
		dpCurr := dp
		if v, ok := best[val]; ok {
			cand := i + v
			if cand > dpCurr {
				dpCurr = cand
			}
		}
		prevBest := dp - (i - 1)
		if v, ok := best[val]; !ok || prevBest > v {
			best[val] = prevBest
		}
		dp = dpCurr
	}
	return dp
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		fmt.Fprintln(writer, maxBallsRemoved(arr))
	}
}
