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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	maxDistinct := 0
	seen := make(map[rune]bool)
	for _, ch := range s {
		if 'a' <= ch && ch <= 'z' {
			seen[ch] = true
		} else {
			if len(seen) > maxDistinct {
				maxDistinct = len(seen)
			}
			seen = make(map[rune]bool)
		}
	}
	if len(seen) > maxDistinct {
		maxDistinct = len(seen)
	}
	fmt.Fprintln(writer, maxDistinct)
}
