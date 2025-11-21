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

	var a, b int64
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	fmt.Fprintln(writer, a*b)
}

