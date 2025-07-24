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
		var A, B int64
		fmt.Fscan(in, &A, &B)
		d := 0
		x := B + 1
		for x > 0 {
			d++
			x /= 10
		}
		fmt.Fprintln(out, A*int64(d-1))
	}
}
