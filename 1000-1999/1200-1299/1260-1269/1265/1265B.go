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
		p := make([]int, n)
		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			pos[p[i]] = i
		}

		ans := make([]byte, n)
		l, r := pos[1], pos[1]
		ans[0] = '1'
		for m := 2; m <= n; m++ {
			if pos[m] < l {
				l = pos[m]
			}
			if pos[m] > r {
				r = pos[m]
			}
			if r-l+1 == m {
				ans[m-1] = '1'
			} else {
				ans[m-1] = '0'
			}
		}
		fmt.Fprintln(out, string(ans))
	}
}
