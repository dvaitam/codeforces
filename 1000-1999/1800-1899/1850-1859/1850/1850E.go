package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var c int64
		fmt.Fscan(in, &n, &c)
		s := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s[i])
		}

		// helper function to calculate total cardboard used for given w
		calc := func(w int64) int64 {
			var sum int64
			for _, v := range s {
				d := v + 2*w
				sum += d * d
				if sum > c {
					return sum
				}
			}
			return sum
		}

		low, high := int64(0), int64(1)
		for calc(high) < c {
			high <<= 1
		}
		for low+1 < high {
			mid := (low + high) / 2
			if calc(mid) >= c {
				high = mid
			} else {
				low = mid
			}
		}
		fmt.Fprintln(out, high)
	}
}
