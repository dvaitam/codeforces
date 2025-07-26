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
		var s string
		fmt.Fscan(reader, &s)
		if isPossible(s) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

func isPossible(s string) bool {
	if len(s) == 0 || s[len(s)-1] != 'B' {
		return false
	}
	countA := 0
	countB := 0
	for _, ch := range s {
		if ch == 'A' {
			countA++
		} else if ch == 'B' {
			countB++
			if countB > countA {
				return false
			}
		}
	}
	return true
}
