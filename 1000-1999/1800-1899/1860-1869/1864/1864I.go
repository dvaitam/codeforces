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
		// precompute snake path indices
		idx := make([][]int, n)
		cnt := 0
		for i := 0; i < n; i++ {
			idx[i] = make([]int, n)
			if i%2 == 0 {
				for j := 0; j < n; j++ {
					cnt++
					idx[i][j] = cnt
				}
			} else {
				for j := n - 1; j >= 0; j-- {
					cnt++
					idx[i][j] = cnt
				}
			}
		}
		used := make([][]bool, n)
		for i := range used {
			used[i] = make([]bool, n)
		}
		prevAns := 0
		for i := 0; i < q; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x ^= prevAns
			y ^= prevAns
			x--
			y--
			used[x][y] = true
			ans := n*n - idx[x][y] + 1
			prevAns = ans
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}
