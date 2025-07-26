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
		var a, b, c, d int
		fmt.Fscan(reader, &a, &b, &c, &d)
		ans := 0
		if b > a {
			ans++
		}
		if c > a {
			ans++
		}
		if d > a {
			ans++
		}
		fmt.Fprintln(writer, ans)
	}
}
