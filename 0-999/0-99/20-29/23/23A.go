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
	n := len(s)
	// search longest substring length that appears at least twice
	for L := n - 1; L >= 1; L-- {
		seen := make(map[string]bool)
		for i := 0; i+L <= n; i++ {
			sub := s[i : i+L]
			if seen[sub] {
				fmt.Fprintln(writer, L)
				return
			}
			seen[sub] = true
		}
	}
	fmt.Fprintln(writer, 0)
}
