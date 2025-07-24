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

	var n int
	var s string
	fmt.Fscan(reader, &n, &s)
	fmt.Fprintln(writer, n+1)
}
