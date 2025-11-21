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
		if n == 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		res := make([]int, n)
		l, r := 0, n-1
		cur := 1
		for l <= r {
			res[l] = cur
			cur++
			l++
			if l <= r {
				res[r] = cur
				cur++
				r--
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
