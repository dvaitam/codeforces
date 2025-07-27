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
		var n, k int
		fmt.Fscan(in, &n, &k)
		ans := 0
		for i := 0; i < n; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if l <= k && k <= r {
				days := r - k + 1
				if days > ans {
					ans = days
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
