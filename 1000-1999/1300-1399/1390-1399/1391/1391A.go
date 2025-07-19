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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		fmt.Fscan(reader, &n)
		for i := n; i >= 1; i-- {
			fmt.Fprint(writer, i)
			if i > 1 {
				fmt.Fprint(writer, " ")
			}
		}
		fmt.Fprint(writer, "\n")
	}
}
