package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)

		var w int
		fmt.Fscan(in, &w)
		heights := make([]int64, w)
		for i := 0; i < w; i++ {
			fmt.Fscan(in, &heights[i])
		}

		rowWeights := dimensionWeights(n, k)
		colWeights := dimensionWeights(m, k)

		totalCells := n * m
		weights := make([]int64, 0, totalCells)
		for _, rw := range rowWeights {
			for _, cw := range colWeights {
				weights = append(weights, rw*cw)
			}
		}

		sort.Slice(heights, func(i, j int) bool { return heights[i] > heights[j] })
		sort.Slice(weights, func(i, j int) bool { return weights[i] > weights[j] })

		limit := w
		if len(weights) < limit {
			limit = len(weights)
		}

		var result int64
		for i := 0; i < limit; i++ {
			result += heights[i] * weights[i]
		}

		fmt.Fprintln(out, result)
	}
}

func dimensionWeights(length, k int) []int64 {
	res := make([]int64, length)
	limit := length - k + 1
	for i := 1; i <= length; i++ {
		l := i - k + 1
		if l < 1 {
			l = 1
		}
		r := i
		if r > limit {
			r = limit
		}
		if r >= l {
			res[i-1] = int64(r - l + 1)
		} else {
			res[i-1] = 0
		}
	}
	return res
}
