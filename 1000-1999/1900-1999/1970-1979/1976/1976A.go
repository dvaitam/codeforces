package main

import (
	"bufio"
	"fmt"
	"os"
)

func strongPassword(s string) bool {
	lastDigit := byte('0')
	lastLetter := byte('a')
	seenLetter := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			if seenLetter {
				return false
			}
			if c < lastDigit {
				return false
			}
			lastDigit = c
		} else {
			seenLetter = true
			if c < lastLetter {
				return false
			}
			lastLetter = c
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		if strongPassword(s) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
