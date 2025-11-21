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
		var n int
		var x int64
		fmt.Fscan(in, &n, &x)

		buckets := make(map[int64][]int64)
		for i := 0; i < n; i++ {
			var val int64
			fmt.Fscan(in, &val)
			r := val % x
			buckets[r] = append(buckets[r], val)
		}

		for r := range buckets {
			sort.Slice(buckets[r], func(i, j int) bool { return buckets[r][i] < buckets[r][j] })
		}

		pointers := make(map[int64]int)
		var mex int64
		for {
			r := mex % x
			list := buckets[r]
			ptr := pointers[r]
			if ptr < len(list) && list[ptr] <= mex {
				pointers[r] = ptr + 1
				mex++
			} else {
				break
			}
		}
		fmt.Fprintln(out, mex)
	}
}
