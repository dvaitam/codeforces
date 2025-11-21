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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		ok := true
		for i := 0; i < n-1 && ok; i++ {
			if a[i] > a[i+1] {
				diff := a[i] - a[i+1]
				a[i] -= diff
				a[i+1] -= diff
			}
			if a[i] < 0 || a[i+1] < 0 {
				ok = false
			}
		}

		if ok {
			for i := 0; i < n-1; i++ {
				if a[i] > a[i+1] {
					ok = false
					break
				}
			}
		}

		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
