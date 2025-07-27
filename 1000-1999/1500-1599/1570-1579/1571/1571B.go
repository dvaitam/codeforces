package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var a, va int
		fmt.Fscan(reader, &a, &va)
		var c, vc int
		fmt.Fscan(reader, &c, &vc)
		var b int
		fmt.Fscan(reader, &b)

		l := va
		if tmp := vc - (c - b); tmp > l {
			l = tmp
		}
		if l < 1 {
			l = 1
		}

		r := vc
		if tmp := va + (b - a); tmp < r {
			r = tmp
		}
		if r > n {
			r = n
		}
		if r > b {
			r = b
		}

		if l < va {
			l = va
		}
		if r < l {
			r = l
		}

		fmt.Fprintln(writer, l)
	}
}
