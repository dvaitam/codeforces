package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt for contest 1696A.
// We can apply the operation on any element of the array, replacing a[i] with
// a[i] OR z and updating z to a[i] AND z. Bits removed from z never reappear,
// so to maximize any element we should apply the operation to that element
// immediately, yielding value a[i] | z. The best possible maximum value is
// therefore the maximum of a[i] | z over all i.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, z int
		fmt.Fscan(reader, &n, &z)
		maxVal := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if v := x | z; v > maxVal {
				maxVal = v
			}
		}
		fmt.Fprintln(writer, maxVal)
	}
}
