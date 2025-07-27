package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var n uint
		fmt.Fscan(in, &n)
		// Find the highest power of two <= n
		if n == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		pow := uint(1) << (bits.Len(n) - 1)
		fmt.Fprintln(out, pow-1)
	}
}
