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
		var s string
		fmt.Fscan(reader, &s)
		// expression is of form a+b where a and b are digits
		if len(s) >= 3 {
			a := int(s[0] - '0')
			b := int(s[2] - '0')
			fmt.Fprintln(writer, a+b)
		}
	}
}
