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
		var a, b int
		fmt.Fscan(reader, &a, &b)
		rem := a % b
		if rem == 0 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, b-rem)
		}
	}
}
