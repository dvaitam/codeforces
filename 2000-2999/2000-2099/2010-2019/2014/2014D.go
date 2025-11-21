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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, d, k int
		fmt.Fscan(in, &n, &d, &k)
		N := n - d + 1
		if N < 1 {
			N = 1
		}
		diff := make([]int, N+3)
		for i := 0; i < k; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			start := l - d + 1
			if start < 1 {
				start = 1
			}
			end := r
			if end > N {
				end = N
			}
			if start <= end {
				diff[start]++
				diff[end+1]--
			}
		}
		count := make([]int, N+2)
		cur := 0
		for i := 1; i <= N; i++ {
			cur += diff[i]
			count[i] = cur
		}
		maxVal, minVal := -1, int(1e9)
		maxPos, minPos := 1, 1
		for i := 1; i <= N; i++ {
			if count[i] > maxVal {
				maxVal = count[i]
				maxPos = i
			}
			if count[i] < minVal {
				minVal = count[i]
				minPos = i
			}
		}
		fmt.Fprintln(out, maxPos, minPos)
	}
}
