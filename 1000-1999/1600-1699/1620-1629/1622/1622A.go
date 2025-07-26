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
		if a+b == c || a+c == b || b+c == a {
			fmt.Fprintln(writer, "YES")
		} else if a == b && c%2 == 0 {
			fmt.Fprintln(writer, "YES")
		} else if a == c && b%2 == 0 {
			fmt.Fprintln(writer, "YES")
		} else if b == c && a%2 == 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
