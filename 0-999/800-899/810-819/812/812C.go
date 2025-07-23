package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var S int64
	if _, err := fmt.Fscan(in, &n, &S); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	costs := make([]int64, n)
	feasible := func(k int) (bool, int64) {
		if k == 0 {
			return true, 0
		}
		for i := 0; i < n; i++ {
			costs[i] = a[i] + int64(i+1)*int64(k)
		}
		sort.Slice(costs, func(i, j int) bool { return costs[i] < costs[j] })
		var sum int64
		for i := 0; i < k; i++ {
			sum += costs[i]
			if sum > S {
				return false, sum
			}
		}
		return sum <= S, sum
	}

	l, r := 0, n
	bestK := 0
	bestSum := int64(0)
	for l <= r {
		m := (l + r) / 2
		ok, sum := feasible(m)
		if ok {
			bestK = m
			bestSum = sum
			l = m + 1
		} else {
			r = m - 1
		}
	}
	fmt.Printf("%d %d\n", bestK, bestSum)
}
