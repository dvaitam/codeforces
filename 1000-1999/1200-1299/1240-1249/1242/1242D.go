package main

import (
	"bufio"
	"fmt"
	"os"
)

func indexInSequence(n int64, k int) int64 {
	used := make(map[int64]bool)
	var p int64 = 1
	var idx int64
	for {
		nums := make([]int64, 0, k)
		for len(nums) < k {
			if !used[p] {
				nums = append(nums, p)
				used[p] = true
			}
			p++
		}
		for _, v := range nums {
			idx++
			if v == n {
				return idx
			}
		}
		var sum int64
		for _, v := range nums {
			sum += v
		}
		idx++
		if sum == n {
			return idx
		}
		used[sum] = true
		if idx > n*2 { // naive bailout
			break
		}
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		var k int
		fmt.Fscan(in, &n, &k)
		fmt.Fprintln(out, indexInSequence(n, k))
	}
}
