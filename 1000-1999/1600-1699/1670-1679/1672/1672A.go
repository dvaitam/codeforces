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
		moves := sum - n
		if moves%2 == 1 {
			fmt.Fprintln(writer, "errorgorn")
		} else {
			fmt.Fprintln(writer, "maomao90")
		}
	}
}
