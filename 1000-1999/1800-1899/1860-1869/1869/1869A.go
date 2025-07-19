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
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		// consume array values
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
		}
		if n%2 == 0 {
			fmt.Fprintln(writer, 2)
			fmt.Fprintln(writer, 1, n)
			fmt.Fprintln(writer, 1, n)
		} else {
			fmt.Fprintln(writer, 4)
			fmt.Fprintln(writer, 1, n-1)
			fmt.Fprintln(writer, 1, n-1)
			fmt.Fprintln(writer, 2, n)
			fmt.Fprintln(writer, 2, n)
		}
	}
}
