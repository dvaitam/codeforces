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
		var l1, r1, l2, r2 int
		fmt.Fscan(reader, &l1, &r1, &l2, &r2)
		if r1 < l2 || r2 < l1 {
			fmt.Fprintln(writer, l1+l2)
		} else {
			if l1 > l2 {
				fmt.Fprintln(writer, l1)
			} else {
				fmt.Fprintln(writer, l2)
			}
		}
	}
}
