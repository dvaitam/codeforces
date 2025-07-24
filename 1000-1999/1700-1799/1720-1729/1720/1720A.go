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
		var a, b, c, d int64
		fmt.Fscan(in, &a, &b, &c, &d)
		ad := a * d
		bc := b * c
		if ad == bc {
			fmt.Fprintln(out, 0)
		} else if ad == 0 || bc == 0 {
			fmt.Fprintln(out, 1)
		} else if ad%bc == 0 || bc%ad == 0 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 2)
		}
	}
}
