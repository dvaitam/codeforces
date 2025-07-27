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
		start := 2*k - n
		// print ascending part
		for i := 1; i < start; i++ {
			fmt.Fprintf(writer, "%d ", i)
		}
		// print descending part
		for i := k; i >= start; i-- {
			fmt.Fprintf(writer, "%d", i)
			if i > start {
				fmt.Fprint(writer, " ")
			}
		}
		if t > 1 {
			fmt.Fprintln(writer)
		}
	}
}
