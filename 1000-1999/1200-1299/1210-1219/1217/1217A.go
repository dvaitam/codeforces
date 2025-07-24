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
		var str, intellect, exp int
		fmt.Fscan(reader, &str, &intellect, &exp)
		tot := intellect + exp - str
		if tot < 0 {
			fmt.Fprintln(writer, exp+1)
			continue
		}
		minX := tot/2 + 1
		if minX > exp {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, exp-minX+1)
		}
	}
}
