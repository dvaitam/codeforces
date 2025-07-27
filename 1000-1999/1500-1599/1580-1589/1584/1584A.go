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
		var u, v int64
		fmt.Fscan(reader, &u, &v)
		x := -u * u
		y := v * v
		fmt.Fprintln(writer, x, y)
	}
}
