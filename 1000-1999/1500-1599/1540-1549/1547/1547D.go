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
	const mask int = (1 << 30) - 1
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		x := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}
		y := make([]int, n)
		prev := x[0] // x[0] ^ y[0] where y[0] = 0
		for i := 1; i < n; i++ {
			y[i] = (^x[i] & mask) & prev
			prev = x[i] ^ y[i]
		}
		for i, v := range y {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(out, "%d", v)
		}
		out.WriteByte('\n')
	}
}
