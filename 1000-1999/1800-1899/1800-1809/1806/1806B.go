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
		c0, c1 := 0, 0
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if arr[i] == 0 {
				c0++
			} else if arr[i] == 1 {
				c1++
			}
		}
		if c0 <= (n+1)/2 {
			fmt.Fprintln(out, 0)
		} else {
			others := n - c0 - c1
			if c1 == 0 || others > 0 {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 2)
			}
		}
	}
}
