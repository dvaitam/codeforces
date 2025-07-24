package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func absorb(order []int64, a []int64, h int64) int {
	idx := 0
	n := len(a)
	for _, m := range order {
		for idx < n && a[idx] < h {
			h += a[idx] / 2
			idx++
		}
		h *= m
	}
	for idx < n && a[idx] < h {
		h += a[idx] / 2
		idx++
	}
	return idx
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var h int64
		fmt.Fscan(reader, &n, &h)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

		orders := [][]int64{{2, 2, 3}, {2, 3, 2}, {3, 2, 2}}
		best := 0
		for _, ord := range orders {
			res := absorb(ord, arr, h)
			if res > best {
				best = res
			}
		}
		fmt.Fprintln(writer, best)
	}
}
