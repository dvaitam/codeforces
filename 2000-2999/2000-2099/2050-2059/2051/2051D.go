package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func countPairsLE(arr []int64, limit int64) int64 {
	if limit < 0 {
		return 0
	}
	n := len(arr)
	var res int64
	j := n - 1
	for i := 0; i < n; i++ {
		if i >= j {
			break
		}
		for i < j && arr[i]+arr[j] > limit {
			j--
		}
		if i >= j {
			break
		}
		res += int64(j - i)
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var x, y int64
		fmt.Fscan(in, &n, &x, &y)
		arr := make([]int64, n)
		var total int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			total += arr[i]
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		L := total - y
		R := total - x
		upper := countPairsLE(arr, R)
		lower := countPairsLE(arr, L-1)
		fmt.Fprintln(out, upper-lower)
	}
}
