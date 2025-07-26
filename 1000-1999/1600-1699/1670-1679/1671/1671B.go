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
		x := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}
		L := int(-1 << 60)
		R := int(1 << 60)
		for i, v := range x {
			a := v - (i + 1)
			if a > L {
				L = a
			}
			if a+2 < R {
				R = a + 2
			}
		}
		if L <= R {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
