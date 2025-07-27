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
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		zeros := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				zeros++
			}
		}
		if zeros == 1 {
			fmt.Fprintln(writer, "BOB")
		} else if zeros%2 == 1 {
			fmt.Fprintln(writer, "ALICE")
		} else {
			fmt.Fprintln(writer, "BOB")
		}
	}
}
