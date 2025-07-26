package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isMeow(s string) bool {
	s = strings.ToLower(s)
	letters := []byte{'m', 'e', 'o', 'w'}
	idx := 0
	for _, ch := range letters {
		if idx >= len(s) || s[idx] != ch {
			return false
		}
		for idx < len(s) && s[idx] == ch {
			idx++
		}
	}
	return idx == len(s)
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
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		if isMeow(s) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
