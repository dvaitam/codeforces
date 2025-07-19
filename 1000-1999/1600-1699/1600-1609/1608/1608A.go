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
		for i := 0; i < n; i++ {
			fmt.Fprint(writer, 3+2*i)
			if i+1 < n {
				fmt.Fprint(writer, " ")
			}
		}
		fmt.Fprint(writer, "\n")
	}
}
