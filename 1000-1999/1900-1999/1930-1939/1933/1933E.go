package main

import (
	"bufio"
	"fmt"
	"os"
)

func calc(s, u int64) int64 {
	return s*u - s*(s-1)/2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		pref := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + a[i]
		}
		var q int
		fmt.Fscan(reader, &q)
		for qi := 0; qi < q; qi++ {
			var l int
			var u int64
			fmt.Fscan(reader, &l, &u)
			base := pref[l-1]
			low, high := l, n+1
			for low < high {
				mid := (low + high) / 2
				if pref[mid]-base <= u {
					low = mid + 1
				} else {
					high = mid
				}
			}
			r1 := low - 1
			if r1 < l {
				r1 = l - 1
			}
			bestR := r1
			bestVal := int64(-1 << 60)
			if r1 >= l {
				s1 := pref[r1] - base
				bestVal = calc(s1, u)
			}
			r2 := r1 + 1
			if r2 <= n {
				s2 := pref[r2] - base
				val2 := calc(s2, u)
				if val2 > bestVal || r1 < l {
					bestVal = val2
					bestR = r2
				}
			}
			fmt.Fprint(writer, bestR)
			if qi+1 < q {
				fmt.Fprint(writer, " ")
			}
		}
		fmt.Fprintln(writer)
	}
}
