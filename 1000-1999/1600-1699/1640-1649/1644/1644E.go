package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; t > 0; t-- {
		var n int64
		var s string
		fmt.Fscan(in, &n, &s)
		hasR := false
		hasD := false
		for _, ch := range s {
			if ch == 'R' {
				hasR = true
			}
			if ch == 'D' {
				hasD = true
			}
		}
		if !hasR || !hasD {
			fmt.Fprintln(writer, n)
			continue
		}
		posR := int64(-1)
		posD := int64(-1)
		for i, ch := range s {
			if ch == 'R' && posR == -1 {
				posR = int64(i + 1)
			}
			if ch == 'D' && posD == -1 {
				posD = int64(i + 1)
			}
		}
		var c int64
		if posR < posD {
			c = posD - 1
		} else {
			c = posR - 1
		}
		ans := n*n - int64(n-1)*c
		fmt.Fprintln(writer, ans)
	}
}
