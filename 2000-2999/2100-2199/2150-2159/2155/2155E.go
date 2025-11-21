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
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)

		if n == 1 {
			parity := 0
			for i := 0; i < k; i++ {
				var x, y int
				fmt.Fscan(in, &x, &y)
				if y == 2 {
					parity ^= 1
				}
			}
			if parity == 1 {
				fmt.Fprintln(out, "Mimo")
			} else {
				fmt.Fprintln(out, "Yuyu")
			}
			continue
		}

		odd := make(map[int]bool)
		for i := 0; i < k; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if y <= 1 {
				continue
			}
			if odd[y] {
				delete(odd, y)
			} else {
				odd[y] = true
			}
		}
		if len(odd) > 0 {
			fmt.Fprintln(out, "Mimo")
		} else {
			fmt.Fprintln(out, "Yuyu")
		}
	}
}
