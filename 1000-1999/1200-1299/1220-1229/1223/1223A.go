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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int64
		fmt.Fscan(reader, &n)
		var result int64
		if n == 2 {
			result = 2
		} else if n%2 == 1 {
			result = 1
		} else {
			result = 0
		}
		fmt.Fprintln(writer, result)
	}
}
