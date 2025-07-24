package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, d int
		fmt.Fscan(in, &n, &m, &d)
		p := make([]int, n)
		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			pos[p[i]] = i + 1
		}
		a := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &a[i])
		}
		ans := n + 5 // large enough
		for i := 0; i < m-1; i++ {
			x := pos[a[i]]
			y := pos[a[i+1]]
			if x > y || y-x > d {
				ans = 0
				break
			}
			diff := y - x
			cur := diff
			delta := d + 1 - diff
			if y+delta <= n {
				if delta < cur {
					cur = delta
				}
			}
			if cur < ans {
				ans = cur
			}
		}
		fmt.Println(ans)
	}
}
