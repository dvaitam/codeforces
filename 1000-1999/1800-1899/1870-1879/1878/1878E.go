package main

import (
	"bufio"
	"fmt"
	"os"
)

const B = 30

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		pref := make([][]int, B)
		for b := 0; b < B; b++ {
			pref[b] = make([]int, n+1)
		}
		for i := 1; i <= n; i++ {
			v := a[i]
			for b := 0; b < B; b++ {
				pref[b][i] = pref[b][i-1]
				if (v>>b)&1 == 0 {
					pref[b][i]++
				}
			}
		}
		var q int
		fmt.Fscan(reader, &q)
		for ; q > 0; q-- {
			var l, k int
			fmt.Fscan(reader, &l, &k)
			lo, hi := l, n
			ans := -1
			for lo <= hi {
				mid := (lo + hi) / 2
				val := (1 << B) - 1
				for b := 0; b < B; b++ {
					if pref[b][mid]-pref[b][l-1] > 0 {
						val &^= 1 << b
					}
				}
				if val >= k {
					ans = mid
					lo = mid + 1
				} else {
					hi = mid - 1
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
