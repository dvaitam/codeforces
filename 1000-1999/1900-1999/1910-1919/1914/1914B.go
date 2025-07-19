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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		k++
		// print from n+1-k to n
		for j := n + 1 - k; j <= n; j++ {
			fmt.Fprint(writer, j, " ")
		}
		// print from n-k down to 1
		for j := n - k; j >= 1; j-- {
			fmt.Fprint(writer, j, " ")
		}
		fmt.Fprint(writer, "\n")
	}
}
