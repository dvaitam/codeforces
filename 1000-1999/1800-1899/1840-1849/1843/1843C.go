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
		var sum int64
		for n > 0 {
			sum += n
			n /= 2
		}
		fmt.Fprintln(writer, sum)
	}
}
