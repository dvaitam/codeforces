package main

import (
	"bufio"
	"fmt"
	"os"
)

const B = 30

func canSet(x, y int64, b uint) bool {
	start := x
	if ((start >> b) & 1) == 0 {
		start = (start>>(b+1))<<(b+1) + (1 << b)
	}
	return start <= y
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
		xs := make([]int64, n+1)
		ys := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &xs[i], &ys[i])
		}
		pref := make([][]int, B)
		for b := 0; b < B; b++ {
			pref[b] = make([]int, n+1)
		}
		for i := 1; i <= n; i++ {
			for b := 0; b < B; b++ {
				pref[b][i] = pref[b][i-1]
				if canSet(xs[i], ys[i], uint(b)) {
					pref[b][i]++
				}
			}
		}
		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			var res int
			for b := 0; b < B; b++ {
				if pref[b][r]-pref[b][l-1] > 0 {
					res |= 1 << b
				}
			}
			fmt.Fprintln(out, res)
		}
	}
}
