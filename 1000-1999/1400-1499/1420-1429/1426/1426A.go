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
	for i := 0; i < t; i++ {
		var n, x int
		fmt.Fscan(reader, &n, &x)
		if n <= 2 {
			fmt.Fprintln(writer, 1)
		} else {
			rem := n - 2
			floors := (rem + x - 1) / x
			fmt.Fprintln(writer, floors+1)
		}
	}
}
