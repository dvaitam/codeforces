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

		deg := make([]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			deg[u]++
			deg[v]++
		}

		var s string
		fmt.Fscan(in, &s)

		c0, c1, k, d := 0, 0, 0, 0
		for i := 1; i < n; i++ {
			if deg[i] == 1 {
				switch s[i] {
				case '0':
					c0++
				case '1':
					c1++
				default:
					k++
				}
			} else {
				if s[i] == '?' {
					d++
				}
			}
		}

		root := s[0]
		var ans int
		if root == '0' {
			ans = c1 + (k+1)/2
		} else if root == '1' {
			ans = c0 + (k+1)/2
		} else {
			scoreA := c0
			if c1 > scoreA {
				scoreA = c1
			}
			scoreA += k / 2
			scoreB := -1
			if d%2 == 1 {
				if c0 < c1 {
					scoreB = c0
				} else {
					scoreB = c1
				}
				scoreB += (k + 1) / 2
			}
			if scoreB > scoreA {
				ans = scoreB
			} else {
				ans = scoreA
			}
		}
		fmt.Fprintln(out, ans)
	}
}
