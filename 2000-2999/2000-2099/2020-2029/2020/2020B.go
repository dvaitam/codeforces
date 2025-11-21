package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isqrt(x int64) int64 {
	if x <= 0 {
		return 0
	}
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func bulbsOnCount(n int64) int64 {
	return n - isqrt(n)
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
		var k int64
		fmt.Fscan(in, &k)
		if k <= 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		low := int64(1)
		high := k
		for bulbsOnCount(high) < k {
			high *= 2
		}
		for low < high {
			mid := (low + high) / 2
			if bulbsOnCount(mid) >= k {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(out, low)
	}
}

