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
		var s1, s2 string
		fmt.Fscan(reader, &s1)
		fmt.Fscan(reader, &s2)
		seen := make(map[byte]struct{})
		for i := 0; i < 2; i++ {
			seen[s1[i]] = struct{}{}
			seen[s2[i]] = struct{}{}
		}
		fmt.Fprintln(writer, len(seen)-1)
	}
}
