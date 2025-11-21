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
		var s, m int64
		fmt.Fscan(in, &n, &s, &m)
		prev := int64(0)
		can := false
		for i := 0; i < n; i++ {
			var l, r int64
			fmt.Fscan(in, &l, &r)
			if !can && l-prev >= s {
				can = true
			}
			prev = r
		}
		if !can && m-prev >= s {
			can = true
		}
		if can {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
