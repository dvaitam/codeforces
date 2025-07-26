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
		var x, y, k int
		fmt.Fscan(reader, &x, &y, &k)
		if y <= x {
			fmt.Fprintln(writer, x)
		} else if y-x <= k {
			fmt.Fprintln(writer, y)
		} else {
			fmt.Fprintln(writer, 2*y-x-k)
		}
	}
}
