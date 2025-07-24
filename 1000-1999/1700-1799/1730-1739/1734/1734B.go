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
		for i := 1; i <= n; i++ {
			for j := 1; j <= i; j++ {
				val := 0
				if j == 1 || j == i {
					val = 1
				}
				if j > 1 {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprint(writer, val)
			}
			fmt.Fprintln(writer)
		}
	}
}
