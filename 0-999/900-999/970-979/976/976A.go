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
	var s string
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}

	zeros := 0
	hasOne := false
	for _, ch := range s {
		if ch == '0' {
			zeros++
		} else if ch == '1' {
			hasOne = true
		}
	}

	if !hasOne {
		fmt.Fprintln(writer, "0")
		return
	}

	fmt.Fprint(writer, "1")
	for i := 0; i < zeros; i++ {
		fmt.Fprint(writer, "0")
	}
	fmt.Fprintln(writer)
}
