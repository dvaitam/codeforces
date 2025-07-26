package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the actual algorithm for problem F1.
// This placeholder reads the input and outputs the initial sum
// minus k as a very naive heuristic.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
			return
		}
		sum := int64(0)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			sum += x
		}
		// Extremely naive approximation: assume each operation
		// decreases the sum by 1. This is not correct, but serves as
		// a placeholder until a proper algorithm is implemented.
		ans := sum - int64(k)
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(out, ans)
	}
}
