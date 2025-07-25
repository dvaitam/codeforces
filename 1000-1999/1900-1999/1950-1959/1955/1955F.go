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
		var c1, c2, c3, c4 int
		if _, err := fmt.Fscan(reader, &c1, &c2, &c3, &c4); err != nil {
			return
		}
		pairs := c1/2 + c2/2 + c3/2 + c4/2
		if c1%2 == 1 && c2%2 == 1 && c3%2 == 1 {
			pairs++
		}
		fmt.Fprintln(writer, pairs)
	}
}
