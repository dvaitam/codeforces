package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type employee struct {
	l int64
	r int64
}

// feasible checks if we can achieve median >= x with total salary <= s.
func feasible(x int64, s int64, arr []employee) bool {
	n := len(arr)
	need := (n + 1) / 2

	sum := int64(0)
	countHigh := 0
	var candidates []int64

	for _, e := range arr {
		if e.r < x {
			sum += e.l
		} else if e.l >= x {
			sum += e.l
			countHigh++
		} else {
			sum += e.l
			candidates = append(candidates, e.l)
		}
	}

	remaining := need - countHigh
	if remaining <= 0 {
		return sum <= s
	}
	if len(candidates) < remaining {
		return false
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i] < candidates[j] })
	for i := len(candidates) - remaining; i < len(candidates); i++ {
		sum += x - candidates[i]
		if sum > s {
			return false
		}
	}
	return sum <= s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var s int64
		fmt.Fscan(reader, &n, &s)
		arr := make([]employee, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i].l, &arr[i].r)
		}

		low := int64(1)
		high := int64(1_000_000_000)
		for low < high {
			mid := (low + high + 1) / 2
			if feasible(mid, s, arr) {
				low = mid
			} else {
				high = mid - 1
			}
		}
		fmt.Fprintln(writer, low)
	}
}
