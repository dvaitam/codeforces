package main

import (
	"bufio"
	"fmt"
	"os"
)

func minOnes(n int) int {
	for ones := 0; ones <= n; ones++ {
		m := n - ones
		if m < 0 {
			break
		}
		for a := 0; 3*a <= m; a++ {
			if (m-3*a)%5 == 0 {
				return ones
			}
		}
	}
	return n
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		fmt.Fprintln(writer, minOnes(n))
	}
}
