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
		lowerSkip := 0
		upperSkip := 0
		result := make([]byte, 0, len(s))
		for i := len(s) - 1; i >= 0; i-- {
			c := s[i]
			if c == 'b' {
				lowerSkip++
			} else if c == 'B' {
				upperSkip++
			} else if c >= 'a' && c <= 'z' {
				if lowerSkip > 0 {
					lowerSkip--
				} else {
					result = append(result, c)
				}
			} else { // uppercase letter
				if upperSkip > 0 {
					upperSkip--
				} else {
					result = append(result, c)
				}
			}
		}
		for i := len(result) - 1; i >= 0; i-- {
			writer.WriteByte(result[i])
		}
		writer.WriteByte('\n')
	}
}
