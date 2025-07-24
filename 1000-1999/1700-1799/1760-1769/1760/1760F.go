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
		var c, d int64
		fmt.Fscan(reader, &n, &c, &d)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })

		maxSingle := a[0]
		if maxSingle*int64(d) < c {
			fmt.Fprintln(writer, "Impossible")
			continue
		}

		// prefix sums
		pre := make([]int64, n+1)
		for i := 0; i < n; i++ {
			pre[i+1] = pre[i] + a[i]
		}
		minLen := d
		if int64(n) < d {
			minLen = int64(n)
		}
		if pre[minLen] >= c {
			fmt.Fprintln(writer, "Infinity")
			continue
		}

		// binary search k in [0, d]
		low, high := int64(0), d
		for low < high {
			mid := (low + high + 1) / 2
			if canAchieve(mid, d, c, a, pre) {
				low = mid
			} else {
				high = mid - 1
			}
		}
		fmt.Fprintln(writer, low)
	}
}

func canAchieve(k, d, c int64, a []int64, pre []int64) bool {
	cycle := k + 1
	idx := cycle
	if idx > int64(len(a)) {
		idx = int64(len(a))
	}
	sumCycle := pre[idx]
	full := d / cycle
	rem := d % cycle
	if rem > int64(len(a)) {
		rem = int64(len(a))
	}
	total := sumCycle*full + pre[rem]
	return total >= c
}
