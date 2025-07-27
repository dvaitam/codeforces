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
		var s string
		fmt.Fscan(reader, &s)
		hasLess := false
		hasGreater := false
		for i := 0; i < len(s); i++ {
			switch s[i] {
			case '<':
				hasLess = true
			case '>':
				hasGreater = true
			}
		}
		if hasLess && hasGreater {
			writer.WriteByte('?')
		} else if hasLess {
			writer.WriteByte('<')
		} else if hasGreater {
			writer.WriteByte('>')
		} else {
			writer.WriteByte('=')
		}
		writer.WriteByte('\n')
	}
}
