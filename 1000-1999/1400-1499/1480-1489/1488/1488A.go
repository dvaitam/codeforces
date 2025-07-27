package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// We want to reach y from 0 using two types of operations:
//  1. add 1;
//  2. add x * 10^p for any p >= 0.
//
// The available increments form a canonical coin system
// {1, x, 10x, 100x, ...}. Hence the minimal number of operations
// is obtained greedily by always using the largest possible
// multiple of x with trailing zeros that does not exceed the
// remaining value. After processing all such denominations,
// the remainder equals the number of +1 operations.
// Complexity per test case is O(log_10 y).
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y int64
		fmt.Fscan(in, &x, &y)

		denom := x
		for denom <= y {
			denom *= 10
		}
		denom /= 10

		ops := int64(0)
		for denom >= x {
			ops += y / denom
			y %= denom
			denom /= 10
		}
		ops += y // remaining +1 operations

		fmt.Fprintln(out, ops)
	}
}
