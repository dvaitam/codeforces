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
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		desc := make([]int, n)
		for i := 0; i+2 < n; i++ {
			if a[i] > a[i+1] && a[i+1] > a[i+2] {
				desc[i+1] = 1 // use 1-based indexing for prefix
			}
		}
		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + desc[i-1]
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if r-l+1 < 3 {
				fmt.Fprintln(out, "YES")
				continue
			}
			// desc index uses 0-based start positions
			if pref[r-2]-pref[l-1] > 0 {
				fmt.Fprintln(out, "NO")
			} else {
				fmt.Fprintln(out, "YES")
			}
		}
	}
}
