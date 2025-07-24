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
		var h, m int
		fmt.Fscan(reader, &h, &m)
		remaining := 24*60 - (h*60 + m)
		fmt.Fprintln(writer, remaining)
	}
}
