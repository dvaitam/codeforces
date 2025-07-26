package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func main() {
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		xor := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			xor ^= x
		}
		fmt.Fprintln(out, xor)
	}
}
