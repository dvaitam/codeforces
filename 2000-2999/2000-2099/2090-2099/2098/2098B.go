package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func countHouses(arr []int, k int) int64 {
	n := len(arr)
	sort.Ints(arr)
	window := n - k
	if window <= 0 {
		return 0
	}
	var total int64
	start := -1
	end := -1
	for i := 0; i+window <= n; i++ {
		var low, high int
		if window%2 == 1 {
			low = arr[i+window/2]
			high = low
		} else {
			low = arr[i+window/2-1]
			high = arr[i+window/2]
		}
		if start == -1 {
			start = low
			end = high
		} else if low <= end+1 {
			if high > end {
				end = high
			}
		} else {
			total += int64(end - start + 1)
			start = low
			end = high
		}
	}
	if start != -1 {
		total += int64(end - start + 1)
	}
	return total
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		result := countHouses(arr, k)
		fmt.Fprintln(out, result)
	}
}
