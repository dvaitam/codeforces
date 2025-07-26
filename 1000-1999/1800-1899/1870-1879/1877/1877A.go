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
		var n int
		fmt.Fscan(reader, &n)
		sum := 0
		for i := 0; i < n-1; i++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += x
		}
		fmt.Fprintln(writer, -sum)
	}
}
