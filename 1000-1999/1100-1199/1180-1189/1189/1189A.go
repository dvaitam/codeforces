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

	// If length is odd, or zeros and ones differ, output whole string
	if n%2 == 1 {
		fmt.Fprintln(writer, 1)
		fmt.Fprintln(writer, s)
		return
	}
	zeros := 0
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			zeros++
		}
	}
	ones := n - zeros
	if zeros != ones {
		fmt.Fprintln(writer, 1)
		fmt.Fprintln(writer, s)
	} else {
		fmt.Fprintln(writer, 2)
		fmt.Fprintf(writer, "%c %s\n", s[0], s[1:])
	}
}
