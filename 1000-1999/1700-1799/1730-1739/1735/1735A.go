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
		var n int64
		fmt.Fscan(reader, &n)
		if n < 6 {
			fmt.Fprintln(writer, 0)
			continue
		}
		fmt.Fprintln(writer, (n-6)/3)
	}
}
