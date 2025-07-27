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
		fmt.Fscan(reader, &n, &s)
		if solve1509B(n, s) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

func solve1509B(n int, s string) bool {
	countT := 0
	countM := 0
	for i := 0; i < n; i++ {
		if s[i] == 'T' {
			countT++
		} else {
			countM++
		}
	}
	if countT != 2*countM {
		return false
	}
	// Check prefix: at any point, number of M can't exceed T before it
	balance := 0
	for i := 0; i < n; i++ {
		if s[i] == 'T' {
			balance++
		} else {
			balance--
		}
		if balance < 0 {
			return false
		}
	}
	// Check suffix: scanning from right, same condition
	balance = 0
	for i := n - 1; i >= 0; i-- {
		if s[i] == 'T' {
			balance++
		} else {
			balance--
		}
		if balance < 0 {
			return false
		}
	}
	return true
}
