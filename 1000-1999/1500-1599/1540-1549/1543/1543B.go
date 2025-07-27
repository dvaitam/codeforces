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
		var sum int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			sum += x
		}
		r := sum % int64(n)
		res := r * (int64(n) - r)
		fmt.Fprintln(writer, res)
	}
}
