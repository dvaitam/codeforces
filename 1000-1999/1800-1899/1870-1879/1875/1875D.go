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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x <= n+1 {
				freq[x]++
			}
		}

		z := freq[0]
		if z == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		// compute mex
		mex := 0
		for {
			if freq[mex] == 0 {
				break
			}
			mex++
		}
		if mex == 1 {
			fmt.Fprintln(out, z-1)
			continue
		}
		costA := int64(mex) * int64(z-1)
		costB := int64(freq[1]-1)*int64(mex) + int64(z)
		if costB < costA {
			costA = costB
		}
		fmt.Fprintln(out, costA)
	}
}
