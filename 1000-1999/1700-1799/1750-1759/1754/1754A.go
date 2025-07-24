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
		var s string
		fmt.Fscan(reader, &s)

		pending := 0
		valid := true
		for i := 0; i < n && valid; i++ {
			if s[i] == 'Q' {
				pending++
			} else {
				if pending > 0 {
					pending--
				} else {
					valid = false
				}
			}
		}
		if valid && pending == 0 {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
