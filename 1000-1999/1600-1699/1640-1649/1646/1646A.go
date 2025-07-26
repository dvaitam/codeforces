package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemA from contest 1646.
// We are given integers n and s. There exists a sequence of n+1 numbers
// where each element is either in [0, n-1] or equal to n^2 and the total
// sum of the sequence is s. The number of elements equal to n^2 is unique.
// Observe that any number equal to n^2 contributes a multiple of n^2 to the
// sum while the remaining numbers are less than n each. Since at most n+1
// numbers are less than n, their total contribution is strictly less than
// n^2. Therefore, the count of n^2 in the sequence is simply floor(s / n^2).
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, s int64
		fmt.Fscan(reader, &n, &s)
		fmt.Fprintln(writer, s/(n*n))
	}
}
