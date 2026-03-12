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
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		var total int64
		for l := 0; l < n; l++ {
			p := a[l]
			var cnt int64
			if b[l] == p {
				cnt++
			}
			total += cnt
			for r := l + 1; r < n; r++ {
				if a[r] > p {
					p = a[r]
					if b[r] == p {
						cnt++
					}
				} else {
					if b[r] <= p {
						cnt++
					}
				}
				total += cnt
			}
		}
		fmt.Fprintln(out, total)
	}
}
