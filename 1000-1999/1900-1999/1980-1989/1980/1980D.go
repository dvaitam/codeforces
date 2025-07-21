package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

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
		if n <= 3 {
			fmt.Fprintln(out, "YES")
			continue
		}
		g := make([]int, n-1)
		for i := 0; i < n-1; i++ {
			g[i] = gcd(a[i], a[i+1])
		}
		pre := make([]bool, n-1)
		pre[0] = true
		for i := 1; i < n-1; i++ {
			pre[i] = pre[i-1] && g[i-1] <= g[i]
		}
		suf := make([]bool, n-1)
		suf[n-2] = true
		for i := n - 3; i >= 0; i-- {
			suf[i] = suf[i+1] && g[i] <= g[i+1]
		}
		possible := false
		for i := 0; i < n; i++ {
			var ok bool
			if i == 0 {
				if n-2 <= 0 {
					ok = true
				} else {
					ok = suf[1]
				}
			} else if i == n-1 {
				if n-3 < 0 {
					ok = true
				} else {
					ok = pre[n-3]
				}
			} else {
				newG := gcd(a[i-1], a[i+1])
				ok = true
				if i-2 >= 0 {
					ok = ok && pre[i-2] && g[i-2] <= newG
				}
				if i+1 <= n-2 {
					ok = ok && suf[i+1] && newG <= g[i+1]
				}
			}
			if ok {
				possible = true
				break
			}
		}
		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
