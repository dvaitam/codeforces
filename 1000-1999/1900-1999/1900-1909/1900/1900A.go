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

		empty := 0
		hasTriple := false
		for i := 0; i < n; i++ {
			if s[i] == '.' {
				empty++
			}
			if i+2 < n && s[i] == '.' && s[i+1] == '.' && s[i+2] == '.' {
				hasTriple = true
			}
		}

		if empty == 0 {
			fmt.Fprintln(writer, 0)
		} else if hasTriple {
			fmt.Fprintln(writer, 2)
		} else {
			fmt.Fprintln(writer, empty)
		}
	}
}
