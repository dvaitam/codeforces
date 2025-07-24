package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var r, b, k int64
		fmt.Fscan(in, &r, &b, &k)
		if r > b {
			r, b = b, r
		}
		g := gcd(r, b)
		r /= g
		b /= g
		if r > b {
			r, b = b, r
		}
		// maximum possible run length of the larger color
		maxRun := (b + r - 2) / r
		if maxRun >= k {
			fmt.Fprintln(out, "REBEL")
		} else {
			fmt.Fprintln(out, "OBEY")
		}
	}
}
