package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem 1848A - "Vika and Her Friends".
// If any friend starts on a cell with the same parity as Vika's starting cell,
// they can always move so that eventually they catch her. Otherwise, their
// parities differ forever and they never share a cell.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		var x, y int
		fmt.Fscan(in, &x, &y)
		startParity := (x + y) % 2
		caught := false
		for i := 0; i < k; i++ {
			var xi, yi int
			fmt.Fscan(in, &xi, &yi)
			if (xi+yi)%2 == startParity {
				caught = true
			}
		}
		if caught {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
