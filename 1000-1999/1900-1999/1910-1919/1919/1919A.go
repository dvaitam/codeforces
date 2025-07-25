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
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		if (a+b)%2 == 1 {
			fmt.Fprintln(writer, "Alice")
		} else {
			fmt.Fprintln(writer, "Bob")
		}
	}
}
