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
	for t > 0 {
		t--
		var n, k int
		fmt.Fscan(in, &n, &k)
		l, r := 1, n
		d := 1
		for l <= r {
			d ^= 1
			if d == 1 {
				fmt.Fprint(out, l, " ")
				l++
			} else {
				fmt.Fprint(out, r, " ")
				r--
			}
		}
		fmt.Fprintln(out)
	}
}
