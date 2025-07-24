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

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	seen := make([]bool, 26)
	next := byte('a')
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'a' || c > 'z' {
			fmt.Fprintln(writer, "NO")
			return
		}
		if !seen[c-'a'] {
			if c != next {
				fmt.Fprintln(writer, "NO")
				return
			}
			seen[c-'a'] = true
			next++
		}
	}
	fmt.Fprintln(writer, "YES")
}
