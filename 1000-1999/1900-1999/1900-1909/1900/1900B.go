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
		one := 0
		two := 0
		three := 0
		if b%2 == c%2 {
			one = 1
		}
		if a%2 == c%2 {
			two = 1
		}
		if a%2 == b%2 {
			three = 1
		}
		fmt.Fprintf(writer, "%d %d %d\n", one, two, three)
	}
}
