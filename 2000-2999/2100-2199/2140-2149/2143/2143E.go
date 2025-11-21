package main

import (
	"bufio"
	"fmt"
	"os"
)

func canTransform(s string) bool {
	open := 0
	balance := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			open++
			balance++
		} else {
			balance--
		}
		if balance < 0 {
			return false
		}
	}
	return open == len(s)/2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		if n%2 != 0 {
			fmt.Fprintln(out, -1)
			continue
		}

		if canTransform(s) {
			fmt.Fprintln(out, s)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
