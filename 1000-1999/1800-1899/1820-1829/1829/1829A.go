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
	const target = "codeforces"
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		diff := 0
		for i := 0; i < len(target) && i < len(s); i++ {
			if s[i] != target[i] {
				diff++
			}
		}
		fmt.Fprintln(writer, diff)
	}
}
