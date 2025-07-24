package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPossible(n int, s string) bool {
	l, r := -1, -1
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			if l == -1 {
				l = i
			}
			r = i
		}
	}
	if l == -1 {
		return true
	}
	for i := l; i <= r; i++ {
		if s[i] == s[n-1-i] {
			return false
		}
	}
	return true
}

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
		fmt.Fscan(reader, &n, &s)
		if isPossible(n, s) {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
