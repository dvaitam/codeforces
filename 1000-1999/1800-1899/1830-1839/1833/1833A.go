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
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)

		mp := make(map[string]struct{})
		for i := 0; i+1 < n; i++ {
			pair := s[i : i+2]
			mp[pair] = struct{}{}
		}
		fmt.Fprintln(writer, len(mp))
	}
}
