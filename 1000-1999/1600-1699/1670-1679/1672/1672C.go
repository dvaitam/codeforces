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
		l, r := -1, -1
		for i := 0; i < n-1; i++ {
			if a[i] == a[i+1] {
				if l == -1 {
					l = i
				}
				r = i
			}
		}
		if l == -1 || l == r {
			fmt.Fprintln(out, 0)
		} else {
			res := r - l - 1
			if res < 1 {
				res = 1
			}
			fmt.Fprintln(out, res)
		}
	}
}
