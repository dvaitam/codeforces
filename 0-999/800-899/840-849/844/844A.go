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
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}

	if len(s) < k {
		fmt.Fprintln(writer, "impossible")
		return
	}

	seen := make(map[rune]struct{})
	for _, ch := range s {
		seen[ch] = struct{}{}
	}
	distinct := len(seen)
	if distinct >= k {
		fmt.Fprintln(writer, 0)
	} else {
		fmt.Fprintln(writer, k-distinct)
	}
}
