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
		fmt.Fscan(in, &n)
		h := make([]int64, n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}
		maxA := int64(0)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxA {
				maxA = a[i]
			}
		}
		res := make([]int64, n)
		for x := int64(1); x <= maxA+1; x++ {
			times := make([]int64, n)
			max1, max2 := int64(-1), int64(-1)
			maxIdx := -1
			for i := 0; i < n; i++ {
				cur := h[i] * ((a[i] + x - 1) / x)
				times[i] = cur
				if cur > max1 {
					max2 = max1
					max1 = cur
					maxIdx = i
				} else if cur > max2 {
					max2 = cur
				}
			}
			cntMax := 0
			for i := 0; i < n; i++ {
				if times[i] == max1 {
					cntMax++
				}
			}
			if cntMax == 1 {
				diff := max1 - max2
				if diff > res[maxIdx] {
					res[maxIdx] = diff
				}
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}
