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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		ops := 0
		len := 1
		for len < n {
			if a[n-len-1] == a[n-1] {
				len++
			} else {
				ops++
				len *= 2
			}
		}
		fmt.Fprintln(out, ops)
	}
}
