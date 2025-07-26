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
		var q int
		fmt.Fscan(in, &q)
		res := make([]byte, q)
		var first, last int
		flag := false
		have := false
		for i := 0; i < q; i++ {
			var x int
			fmt.Fscan(in, &x)
			if !have {
				// first element is always accepted
				first = x
				last = x
				have = true
				res[i] = '1'
				continue
			}
			if !flag {
				if x >= last {
					last = x
					res[i] = '1'
				} else if x <= first {
					last = x
					flag = true
					res[i] = '1'
				} else {
					res[i] = '0'
				}
			} else {
				if x >= last && x <= first {
					last = x
					res[i] = '1'
				} else {
					res[i] = '0'
				}
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
