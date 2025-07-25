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
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	ones := 0
	zeros := 0
	for _, ch := range s {
		if ch == 'n' {
			ones++
		} else if ch == 'z' {
			zeros++
		}
	}

	first := true
	for i := 0; i < ones; i++ {
		if !first {
			writer.WriteByte(' ')
		}
		writer.WriteByte('1')
		first = false
	}
	for i := 0; i < zeros; i++ {
		if !first {
			writer.WriteByte(' ')
		}
		writer.WriteByte('0')
		first = false
	}
	writer.WriteByte('\n')
}
