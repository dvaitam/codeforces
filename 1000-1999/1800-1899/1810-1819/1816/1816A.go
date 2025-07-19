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
	for i := 0; i < t; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		fmt.Fprintln(writer, 2)
		fmt.Fprintln(writer, a-1, 1)
		fmt.Fprintln(writer, a, b)
	}
}
