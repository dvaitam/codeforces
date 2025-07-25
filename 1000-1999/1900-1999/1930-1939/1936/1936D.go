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
		var v int64
		fmt.Fscan(in, &n, &v)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var typ int
			fmt.Fscan(in, &typ)
			if typ == 1 {
				var idx int
				var x int64
				fmt.Fscan(in, &idx, &x)
				b[idx-1] = x
			} else if typ == 2 {
				var l, r int
				fmt.Fscan(in, &l, &r)
				l--
				r--
				const INF int64 = 1<<63 - 1
				ans := INF
				for L := l; L <= r; L++ {
					orVal := int64(0)
					maxA := int64(0)
					for R := L; R <= r; R++ {
						orVal |= b[R]
						if a[R] > maxA {
							maxA = a[R]
						}
						if orVal >= v {
							if maxA < ans {
								ans = maxA
							}
							break
						}
					}
				}
				if ans == INF {
					fmt.Fprintln(out, -1)
				} else {
					fmt.Fprintln(out, ans)
				}
			}
		}
	}
}
