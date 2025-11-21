package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		zeros, negatives := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x == 0 {
				zeros++
			} else if x < 0 {
				negatives++
			}
		}

		// Convert every zero to one (one operation each) and if the number of negatives is odd,
		// flip one of them to positive (two more operations) to make the total product positive.
		operations := zeros
		if negatives%2 == 1 {
			operations += 2
		}

		fmt.Fprintln(writer, operations)
	}
}
