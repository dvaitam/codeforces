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
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		m := n / 2
		if m == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		leftVal := func(idx int, state int) int {
			if state == 0 {
				return a[idx]
			}
			return a[n-idx+1]
		}
		rightVal := func(idx int, state int) int {
			if state == 0 {
				return a[n-idx+1]
			}
			return a[idx]
		}
		dpPrev := [2]int{0, 0}
		dpCurr := [2]int{0, 0}
		const INF = int(1e9)
		for i := 2; i <= m; i++ {
			for state := 0; state < 2; state++ {
				dpCurr[state] = INF
			}
			for prev := 0; prev < 2; prev++ {
				for cur := 0; cur < 2; cur++ {
					cost := 0
					if leftVal(i-1, prev) == leftVal(i, cur) {
						cost++
					}
					if rightVal(i-1, prev) == rightVal(i, cur) {
						cost++
					}
					if tmp := dpPrev[prev] + cost; tmp < dpCurr[cur] {
						dpCurr[cur] = tmp
					}
				}
			}
			dpPrev = dpCurr
		}
		if m == 1 {
			dpPrev[0] = 0
			dpPrev[1] = 0
		}
		ans := INF
		for state := 0; state < 2; state++ {
			add := 0
			if n%2 == 0 {
				if leftVal(m, state) == rightVal(m, state) {
					add++
				}
			} else {
				center := a[m+1]
				if leftVal(m, state) == center {
					add++
				}
				if center == rightVal(m, state) {
					add++
				}
			}
			if tmp := dpPrev[state] + add; tmp < ans {
				ans = tmp
			}
		}
		fmt.Fprintln(out, ans)
	}
}
