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
		switch {
		case n == 1:
			fmt.Fprintln(out, 9)
		case n == 2:
			fmt.Fprintln(out, "98")
		case n == 3:
			fmt.Fprintln(out, "989")
		default:
			res := make([]byte, n)
			res[0] = '9'
			res[1] = '8'
			res[2] = '9'
			for i := 3; i < n; i++ {
				res[i] = byte('0' + (i-3)%10)
			}
			fmt.Fprintln(out, string(res))
		}
	}
}
