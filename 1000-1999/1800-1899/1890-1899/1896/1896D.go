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
		a := make([]int, n+1)
		total := 0
		firstOne := n + 1
		lastOne := 0
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
			total += a[i]
			if a[i] == 1 {
				if firstOne > n {
					firstOne = i
				}
				lastOne = i
			}
		}
		for ; q > 0; q-- {
			var op int
			fmt.Fscan(in, &op)
			if op == 1 {
				var s int
				fmt.Fscan(in, &s)
				if s > total {
					fmt.Fprintln(out, "NO")
					continue
				}
				L := n
				if firstOne <= n {
					L = firstOne - 1
				}
				R := n
				if lastOne >= 1 {
					R = n - lastOne
				}
				if L < R {
					R = L
				}
				d := total - s
				if d%2 == 1 && (d+1)/2 <= R {
					fmt.Fprintln(out, "NO")
				} else {
					fmt.Fprintln(out, "YES")
				}
			} else {
				var i, v int
				fmt.Fscan(in, &i, &v)
				if a[i] == v {
					continue
				}
				old := a[i]
				a[i] = v
				total += v - old
				if old == 1 {
					if i == firstOne {
						for firstOne <= n && a[firstOne] == 2 {
							firstOne++
						}
						if firstOne > n {
							firstOne = n + 1
						}
					}
					if i == lastOne {
						for lastOne >= 1 && a[lastOne] == 2 {
							lastOne--
						}
						if lastOne < 1 {
							lastOne = 0
						}
					}
				} else {
					if v == 1 {
						if i < firstOne {
							firstOne = i
						}
						if i > lastOne {
							lastOne = i
						}
					}
				}
			}
		}
	}
}
