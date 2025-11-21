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
		a := make([]int64, n+1)
		pref := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
			if pref[i-1] > a[i] {
				pref[i] = pref[i-1]
			} else {
				pref[i] = a[i]
			}
		}
		suff := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			suff[i] = suff[i-1] + a[n-i+1]
		}
		for k := 1; k <= n; k++ {
			idx := n - k + 1
			base := suff[k-1]
			val := pref[idx]
			if a[idx] > val {
				val = a[idx]
			}
			if k > 1 {
				temp := base + a[idx]
				temp2 := base + pref[idx]
				if temp2 > temp {
					temp = temp2
				}
				fmt.Fprint(out, temp)
			} else {
				fmt.Fprint(out, val)
			}
			if k != n {
				fmt.Fprint(out, " ")
			}
		}
		fmt.Fprintln(out)
	}
}
