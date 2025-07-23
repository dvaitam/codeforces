package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func solve(x int64) (int64, int64, bool) {
	for k := int64(1); k*k <= x+int64(1e9); k++ {
		n2 := x + k*k
		n := int64(math.Sqrt(float64(n2)))
		if n*n != n2 {
			continue
		}
		lower := n/(k+1) + 1
		upper := n / k
		if lower <= upper && upper <= n {
			return n, upper, true
		}
	}
	return 0, 0, false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x int64
		fmt.Fscan(in, &x)
		n, m, ok := solve(x)
		if !ok {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintf(out, "%d %d\n", n, m)
		}
	}
}
