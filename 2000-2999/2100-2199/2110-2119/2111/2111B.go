package main

import (
	"bufio"
	"fmt"
	"os"
)

// fib[i] = i-th Fibonacci number with f1=1, f2=2, f3=3, ...
var fib = []int{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}

// canFit returns true if all n Fibonacci cubes fit in a box with the given dimensions.
// The optimal arrangement tiles a 2D floor using the Fibonacci rectangle property:
// cubes f_1..f_n tile a f_n × f_{n+1} rectangle as a single layer.
// Necessary and sufficient condition (sorted a<=b<=c): a>=fib[n], b>=fib[n], c>=fib[n+1].
func canFit(n int, x, y, z int) bool {
	a, b, c := x, y, z
	if a > b {
		a, b = b, a
	}
	if b > c {
		b, c = c, b
	}
	if a > b {
		a, b = b, a
	}
	return a >= fib[n] && b >= fib[n] && c >= fib[n+1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		for i := 0; i < m; i++ {
			var x, y, z int
			fmt.Fscan(in, &x, &y, &z)
			if canFit(n, x, y, z) {
				fmt.Fprint(out, "1")
			} else {
				fmt.Fprint(out, "0")
			}
		}
		fmt.Fprintln(out)
	}
}
