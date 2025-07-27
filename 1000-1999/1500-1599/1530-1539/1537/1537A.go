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
		var n int
		fmt.Fscan(reader, &n)
		sum := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += x
		}
		if sum == n {
			fmt.Fprintln(writer, 0)
		} else if sum < n {
			fmt.Fprintln(writer, 1)
		} else {
			fmt.Fprintln(writer, sum-n)
		}
	}
}
