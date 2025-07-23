package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isBeautiful(s string) bool {
	count1 := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			count1++
			if count1 > 1 {
				return false
			}
		} else if s[i] != '0' {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	special := "1"
	zeros := 0
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		if s == "0" {
			fmt.Fprintln(writer, "0")
			return
		}
		if isBeautiful(s) {
			if s != "1" {
				zeros += len(s) - 1
			}
		} else {
			special = s
		}
	}
	fmt.Fprint(writer, special)
	if zeros > 0 {
		fmt.Fprint(writer, strings.Repeat("0", zeros))
	}
}
