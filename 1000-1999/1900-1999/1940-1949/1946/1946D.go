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
	const mask30 int64 = (1 << 30) - 1
	for ; t > 0; t-- {
		var n int
		var x int64
		fmt.Fscan(in, &n, &x)
		mask := (^x) & mask30
		px := int64(0)
		base := int64(0)
		segments := 0
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			px ^= v & mask
			if px == base {
				segments++
				base = px
			}
		}
		if px == base {
			fmt.Fprintln(out, segments)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
