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
		var n int64
		fmt.Fscan(reader, &n)
		for n%2 == 0 {
			n /= 2
		}
		if n > 1 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
