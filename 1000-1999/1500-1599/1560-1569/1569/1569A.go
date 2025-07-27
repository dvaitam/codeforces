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
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		found := false
		l, r := -1, -1
		for i := 0; i < n-1; i++ {
			if s[i] != s[i+1] {
				l = i + 1
				r = i + 2
				found = true
				break
			}
		}
		if found {
			fmt.Fprintln(writer, l, r)
		} else {
			fmt.Fprintln(writer, -1, -1)
		}
	}
}
