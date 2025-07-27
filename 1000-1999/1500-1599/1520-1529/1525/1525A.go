package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemA from contest 1525.
// To obtain a potion with exactly k percent magic essence,
// the ratio of essence to water must be k:(100-k).
// The minimal total volume is 100 divided by gcd(k,100).
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var k int
		fmt.Fscan(in, &k)
		g := gcd(k, 100)
		fmt.Fprintln(out, 100/g)
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
