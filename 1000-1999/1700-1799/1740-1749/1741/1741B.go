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
	var t, n int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		fmt.Fscan(reader, &n)
		if n == 3 {
			fmt.Fprintln(writer, -1)
		} else {
			for i := 3; i <= n; i++ {
				fmt.Fprint(writer, i, " ")
			}
			fmt.Fprintln(writer, "2 1")
		}
	}
}
