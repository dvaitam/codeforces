package main

import (
	"bufio"
	"fmt"
	"os"
)

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
		var x int64
		fmt.Fscan(reader, &n, &x)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		l, r := a[0]-x, a[0]+x
		changes := 0
		for i := 1; i < n; i++ {
			nl, nr := a[i]-x, a[i]+x
			if nl > r || nr < l {
				changes++
				l, r = nl, nr
			} else {
				if nl > l {
					l = nl
				}
				if nr < r {
					r = nr
				}
			}
		}
		fmt.Fprintln(writer, changes)
	}
}
