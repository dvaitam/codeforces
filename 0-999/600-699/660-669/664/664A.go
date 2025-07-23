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

	var a, b string
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	if a == b {
		fmt.Fprintln(writer, a)
	} else {
		fmt.Fprintln(writer, 1)
	}
}
