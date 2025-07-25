package main

import (
	"bufio"
	"fmt"
	"os"
)

func solveQuery(a []int64, k int64) int64 {
	n := len(a)
	cur := make([]int64, n)
	copy(cur, a)
	var ans int64
	for bit := 60; bit >= 0; bit-- {
		var cost int64
		limit := int64(1) << (bit + 1)
		half := int64(1) << bit
		for i := 0; i < n; i++ {
			if (cur[i]>>bit)&1 == 0 {
				rem := cur[i] % limit
				cost += half - rem
			}
		}
		if cost <= k {
			k -= cost
			ans |= half
			for i := 0; i < n; i++ {
				if (cur[i]>>bit)&1 == 0 {
					rem := cur[i] % limit
					cur[i] += half - rem
				}
			}
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	for ; q > 0; q-- {
		var k int64
		fmt.Fscan(in, &k)
		res := solveQuery(arr, k)
		fmt.Fprintln(out, res)
	}
}
