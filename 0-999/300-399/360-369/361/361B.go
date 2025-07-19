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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	// If all elements are good, no solution
	if n == k {
		fmt.Fprint(writer, -1)
		return
	}
	m := n - k
	// Construct permutation: rotate first m elements by 1, leave rest fixed
	for i := 1; i <= n; i++ {
		var v int
		switch {
		case i < m:
			v = i + 1
		case i == m:
			v = 1
		default:
			v = i
		}
		fmt.Fprint(writer, v)
		if i < n {
			fmt.Fprint(writer, " ")
		}
	}
	fmt.Fprint(writer, '\n')
}
