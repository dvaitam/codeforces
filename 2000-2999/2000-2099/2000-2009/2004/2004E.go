package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxA = 10_000_000

var grundy = precomputeGrundy()

func precomputeGrundy() []uint32 {
	g := make([]uint32, maxA+1)
	g[1] = 1
	rank := uint32(1) // corresponds to prime 3
	for p := 3; p <= maxA; p += 2 {
		if g[p] != 0 {
			continue
		}
		val := rank + 1
		step := p << 1
		for m := p; m <= maxA; m += step {
			if g[m] == 0 {
				g[m] = val
			}
		}
		rank++
	}
	return g
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var nim uint32
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x&1 == 1 {
				nim ^= grundy[x]
			}
		}
		if nim != 0 {
			fmt.Fprintln(out, "Alice")
		} else {
			fmt.Fprintln(out, "Bob")
		}
	}
}
