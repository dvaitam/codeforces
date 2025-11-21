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
		var s string
		fmt.Fscan(in, &s)

		lo, hi := 1, n
		ans := make([]int, n)
		for i := 0; i < n; i++ {
			ans[i] = -1
		}

		if len(s) > 0 && s[0] == '>' {
			ans[0] = hi
			hi--
		} else {
			ans[0] = lo
			lo++
		}

		for i := 1; i < n; i++ {
			if s[i-1] == '<' {
				ans[i] = lo
				lo++
			} else {
				ans[i] = hi
				hi--
			}
		}

		for i := 0; i < n; i++ {
			fmt.Fprintf(out, "%d ", ans[i])
		}
		fmt.Fprintln(out)
	}
}
