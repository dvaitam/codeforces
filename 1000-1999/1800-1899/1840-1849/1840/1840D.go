package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func canCover(a []int64, d int64) bool {
	n := len(a)
	i := 0
	for k := 0; k < 3 && i < n; k++ {
		limit := a[i] + 2*d
		i++
		for i < n && a[i] <= limit {
			i++
		}
	}
	return i == n
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		low, high := int64(0), int64(1_000_000_000)
		for low < high {
			mid := (low + high) / 2
			if canCover(arr, mid) {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(writer, low)
	}
}
