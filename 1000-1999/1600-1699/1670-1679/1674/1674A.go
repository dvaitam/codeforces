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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		if y%x != 0 {
			fmt.Fprintln(writer, "0 0")
		} else {
			fmt.Fprintf(writer, "1 %d\n", y/x)
		}
	}
}
