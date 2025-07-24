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
		var q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		res := make([]byte, n)
		iq := q
		for i := n - 1; i >= 0; i-- {
			if a[i] <= iq {
				res[i] = '1'
			} else if iq > 0 {
				res[i] = '1'
				iq--
			} else {
				res[i] = '0'
			}
		}
		out.Write(res)
		if t > 0 {
			out.WriteByte('\n')
		}
	}
}
