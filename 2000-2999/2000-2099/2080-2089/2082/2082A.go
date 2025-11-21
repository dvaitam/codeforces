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
		var n, m int
		fmt.Fscan(in, &n, &m)
		colParity := make([]int, m)
		rowOdd := 0
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			parity := 0
			for j := 0; j < m; j++ {
				if s[j] == '1' {
					parity ^= 1
					colParity[j] ^= 1
				}
			}
			if parity == 1 {
				rowOdd++
			}
		}
		colOdd := 0
		for j := 0; j < m; j++ {
			if colParity[j] == 1 {
				colOdd++
			}
		}
		if rowOdd < colOdd {
			fmt.Fprintln(out, colOdd)
		} else {
			fmt.Fprintln(out, rowOdd)
		}
	}
}
