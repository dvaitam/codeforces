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

	var t, n int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		fmt.Fscan(reader, &n)
		var p []int
		if n%2 == 1 {
			// odd n: pairs from n-2 down to 5, then 1,2,3
			for i := n - 2; i >= 5; i -= 2 {
				p = append(p, i, i-1)
			}
			// append 1,2,3
			p = append(p, 1, 2, 3)
		} else {
			// even n: pairs from n-2 down to 2
			for i := n - 2; i >= 2; i -= 2 {
				p = append(p, i, i-1)
			}
		}
		// append the last two elements n-1, n
		p = append(p, n-1, n)

		// output permutation
		for i, v := range p {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprintf(writer, "%d", v)
		}
		writer.WriteByte('\n')
	}
}
