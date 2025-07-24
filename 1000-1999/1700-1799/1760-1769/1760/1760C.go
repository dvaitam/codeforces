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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		max1, max2 := -1, -1
		for _, v := range a {
			if v > max1 {
				max2 = max1
				max1 = v
			} else if v > max2 {
				max2 = v
			}
		}
		for i, v := range a {
			diff := 0
			if v == max1 {
				diff = v - max2
			} else {
				diff = v - max1
			}
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, diff)
		}
		out.WriteByte('\n')
	}
}
