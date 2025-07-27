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
		var u, v int
		fmt.Fscan(in, &n, &u, &v)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		diffGreater := false
		diffOne := false
		for i := 0; i+1 < n; i++ {
			d := a[i] - a[i+1]
			if d < 0 {
				d = -d
			}
			if d > 1 {
				diffGreater = true
			}
			if d == 1 {
				diffOne = true
			}
		}
		if diffGreater {
			fmt.Fprintln(out, 0)
		} else if diffOne {
			if u < v {
				fmt.Fprintln(out, u)
			} else {
				fmt.Fprintln(out, v)
			}
		} else {
			cost1 := u + v
			cost2 := 2 * v
			if cost1 < cost2 {
				fmt.Fprintln(out, cost1)
			} else {
				fmt.Fprintln(out, cost2)
			}
		}
	}
}
