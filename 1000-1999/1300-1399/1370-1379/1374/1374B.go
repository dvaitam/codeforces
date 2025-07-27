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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		c2, c3 := 0, 0
		for n%2 == 0 {
			n /= 2
			c2++
		}
		for n%3 == 0 {
			n /= 3
			c3++
		}
		if n != 1 || c2 > c3 {
			fmt.Fprintln(writer, -1)
			continue
		}
		ops := (c3 - c2) + c3
		fmt.Fprintln(writer, ops)
	}
}
