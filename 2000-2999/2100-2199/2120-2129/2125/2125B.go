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
	if a < 0 {
		return -a
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
		var a, b, k int64
		fmt.Fscan(in, &a, &b, &k)
		g := gcd(a, b)
		if a/g <= k && b/g <= k {
			fmt.Fprintln(out, 1)
		} else {
			// Using (1,0) and (0,1) always works because k >= 1.
			fmt.Fprintln(out, 2)
		}
	}
}
