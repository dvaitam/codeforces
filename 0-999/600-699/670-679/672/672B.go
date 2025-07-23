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

	if n > 26 {
		fmt.Fprintln(writer, -1)
		return
	}

	seen := make([]bool, 26)
	changes := 0
	for _, ch := range s {
		idx := ch - 'a'
		if seen[idx] {
			changes++
		} else {
			seen[idx] = true
		}
	}
	fmt.Fprintln(writer, changes)
}
