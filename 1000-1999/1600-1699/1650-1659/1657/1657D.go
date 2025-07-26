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

	var n, C int
	if _, err := fmt.Fscan(in, &n, &C); err != nil {
		return
	}
	bestUnit := make([]int64, C+1)
	for i := 0; i < n; i++ {
		var c int
		var d, h int64
		fmt.Fscan(in, &c, &d, &h)
		val := d * h
		if val > bestUnit[c] {
			bestUnit[c] = val
		}
	}

	maxVal := make([]int64, C+1)
	for c := 1; c <= C; c++ {
		if bestUnit[c] == 0 {
			continue
		}
		val := bestUnit[c]
		for cost := c; cost <= C; cost += c {
			k := int64(cost / c)
			candidate := val * k
			if candidate > maxVal[cost] {
				maxVal[cost] = candidate
			}
		}
	}

	for i := 1; i <= C; i++ {
		if maxVal[i] < maxVal[i-1] {
			maxVal[i] = maxVal[i-1]
		}
	}

	var m int
	fmt.Fscan(in, &m)
	for j := 0; j < m; j++ {
		var D, H int64
		fmt.Fscan(in, &D, &H)
		target := D * H
		l, r := 1, C
		ans := -1
		for l <= r {
			mid := (l + r) / 2
			if maxVal[mid] > target {
				ans = mid
				r = mid - 1
			} else {
				l = mid + 1
			}
		}
		if j > 0 {
			fmt.Fprint(out, " ")
		}
		if ans == -1 {
			fmt.Fprint(out, -1)
		} else {
			fmt.Fprint(out, ans)
		}
	}
	fmt.Fprintln(out)
}
