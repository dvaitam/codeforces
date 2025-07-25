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
		var k, m int64
		fmt.Fscan(reader, &k, &m)
		mod := m % (3 * k)
		if mod >= 2*k {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, 2*k-mod)
		}
	}
}
