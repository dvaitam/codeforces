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
		base := n / 3
		r := n % 3
		c1, c2 := base, base
		if r == 1 {
			c1++
		} else if r == 2 {
			c2++
		}
		fmt.Fprintln(writer, c1, c2)
	}
}
