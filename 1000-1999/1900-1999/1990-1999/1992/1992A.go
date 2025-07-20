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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		best := 0
		for i := 0; i <= 5; i++ {
			for j := 0; j <= 5-i; j++ {
				k := 5 - i - j
				prod := (a + i) * (b + j) * (c + k)
				if prod > best {
					best = prod
				}
			}
		}
		fmt.Fprintln(writer, best)
	}
}
