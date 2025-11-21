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
		var w, b int64
		fmt.Fscan(in, &w, &b)
		sum := w + b
		lo, hi := int64(0), int64(1)
		for hi*(hi+1)/2 <= sum {
			hi <<= 1
		}
		for lo < hi {
			mid := (lo + hi + 1) >> 1
			if mid*(mid+1)/2 <= sum {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		fmt.Fprintln(out, lo)
	}
}
