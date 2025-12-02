package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for i := 0; i < t; i++ {
		var n, m uint64
		fmt.Fscan(in, &n, &m)

		if (n ^ m) < n {
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, n, m)
			continue
		}

		k := bits.Len64(n) - 1
		xBit := uint64(1) << k
		r := n ^ xBit

		if r == 0 {
			fmt.Fprintln(out, -1)
			continue
		}

		if m < r {
			fmt.Fprintln(out, 2)
			fmt.Fprintln(out, n, xBit|m, m)
			continue
		}

		fmt.Fprintln(out, -1)
	}
}