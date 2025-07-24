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
		var n int
		fmt.Fscan(reader, &n)
		if n%2 == 1 {
			// odd n: start with 1 then output pairs (i+1,i)
			fmt.Fprint(writer, 1)
			for i := 2; i <= n; i += 2 {
				fmt.Fprintf(writer, " %d %d", i+1, i)
			}
			fmt.Fprintln(writer)
		} else {
			// even n: output pairs (i+1,i)
			for i := 1; i <= n; i += 2 {
				if i > 1 {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprintf(writer, "%d %d", i+1, i)
			}
			fmt.Fprintln(writer)
		}
	}
}
