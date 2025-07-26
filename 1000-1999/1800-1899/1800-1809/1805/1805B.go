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
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)

		minChar := byte('z' + 1)
		for i := 0; i < len(s); i++ {
			if s[i] < minChar {
				minChar = s[i]
			}
		}
		idx := -1
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] == minChar {
				idx = i
				break
			}
		}
		result := make([]byte, 0, len(s))
		result = append(result, minChar)
		result = append(result, s[:idx]...)
		if idx+1 < len(s) {
			result = append(result, s[idx+1:]...)
		}
		fmt.Fprintln(writer, string(result))
	}
}
