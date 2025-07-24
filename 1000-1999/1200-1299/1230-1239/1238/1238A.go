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
		var x, y int64
		fmt.Fscan(reader, &x, &y)
		if x-y == 1 {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
