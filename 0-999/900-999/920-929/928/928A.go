package main

import (
	"bufio"
	"fmt"
	"os"
)

func canonical(s string) string {
	b := []byte(s)
	for i, c := range b {
		switch c {
		case 'O', 'o', '0':
			b[i] = '0'
		case 'I', 'i', 'l', 'L', '1':
			b[i] = '1'
		default:
			if c >= 'A' && c <= 'Z' {
				b[i] = c - 'A' + 'a'
			}
		}
	}
	return string(b)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	normS := canonical(s)
	ans := "Yes"
	for i := 0; i < n; i++ {
		var t string
		fmt.Fscan(reader, &t)
		if len(t) != len(s) {
			continue
		}
		if canonical(t) == normS {
			ans = "No"
		}
	}
	fmt.Fprintln(writer, ans)
}
