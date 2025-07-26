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
		var n int64
		fmt.Fscan(in, &n)
		// factor out powers of two
		tmp := n
		a := 0
		for tmp%2 == 0 {
			tmp /= 2
			a++
		}
		odd := tmp
		// try using the odd part as k if possible
		if odd > 1 {
			k := odd
			if k*(k+1)/2 <= n {
				fmt.Fprintln(out, k)
				continue
			}
		}
		// try k = 2^(a+1)
		k := int64(1) << (a + 1)
		if k*(k+1)/2 <= n {
			fmt.Fprintln(out, k)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
