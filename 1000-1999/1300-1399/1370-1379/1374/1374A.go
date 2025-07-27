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
		var x, y, n int64
		fmt.Fscan(reader, &x, &y, &n)
		k := n - (n-y)%x
		fmt.Fprintln(writer, k)
	}
}
