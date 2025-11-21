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
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		ans := make([]int, n)
		for r := 0; r < n; r++ {
			best := 0
			for l := r; l >= 0; l-- {
				mex := 0
				present := make(map[int]bool)
				for i := l; i <= r; i++ {
					present[arr[i]] = true
				}
				for present[mex] {
					mex++
				}
				w := 0
				for i := l; i <= r; i++ {
					if arr[i] > mex {
						w++
					}
				}
				if w > best {
					best = w
				}
			}
			ans[r] = best
		}
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
