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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	if n == 1 {
		fmt.Fprintln(writer, "9 8")
	} else {
		fmt.Fprintf(writer, "%d %d\n", 3*n, 2*n)
	}
}
