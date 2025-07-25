package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(s string) bool {
	if len(s) < 2 || s[0] != '1' || s[len(s)-1] == '9' {
		return false
	}
	for i := 1; i < len(s)-1; i++ {
		if s[i] == '0' {
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
		var x string
		fmt.Fscan(reader, &x)
		if check(x) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
