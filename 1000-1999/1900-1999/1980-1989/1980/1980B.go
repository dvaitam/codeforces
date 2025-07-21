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
		var n, f, k int
		fmt.Fscan(in, &n, &f, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fav := a[f-1]
		gt, eq := 0, 0
		for _, v := range a {
			if v > fav {
				gt++
			} else if v == fav {
				eq++
			}
		}
		earliest := gt + 1
		latest := gt + eq
		if k < earliest {
			fmt.Fprintln(out, "NO")
		} else if k >= latest {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "MAYBE")
		}
	}
}
