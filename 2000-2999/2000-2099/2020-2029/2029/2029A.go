package main

import (
	"bufio"
	"fmt"
	"os"
)

func countOps(l, r, k int64) int64 {
	if k == 1 {
		return r - l + 1
	}
	ans := int64(0)
	cur := l
	for cur <= r {
		q := r / cur
		blockEnd := r / q
		if blockEnd > r {
			blockEnd = r
		}
		if q >= k {
			t := q - k
			threshold := ((l - 1) / (t + 1)) + 1
			if threshold < l {
				threshold = l
			}
			start := cur
			if threshold > start {
				start = threshold
			}
			if start <= blockEnd {
				ans += blockEnd - start + 1
			}
		}
		cur = blockEnd + 1
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r, k int64
		fmt.Fscan(in, &l, &r, &k)
		fmt.Fprintln(out, countOps(l, r, k))
	}
}
