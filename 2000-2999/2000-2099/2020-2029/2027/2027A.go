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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		maxW, maxH := 0, 0
		for i := 0; i < n; i++ {
			var w, h int
			fmt.Fscan(in, &w, &h)
			if w > maxW {
				maxW = w
			}
			if h > maxH {
				maxH = h
			}
		}
		fmt.Fprintln(out, 2*(maxW+maxH))
	}
}
