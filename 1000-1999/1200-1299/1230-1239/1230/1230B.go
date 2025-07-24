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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	if n == 1 {
		if k > 0 {
			fmt.Fprintln(writer, "0")
		} else {
			fmt.Fprintln(writer, s)
		}
		return
	}

	b := []byte(s)
	if k > 0 && b[0] != '1' {
		b[0] = '1'
		k--
	}
	for i := 1; i < n && k > 0; i++ {
		if b[i] != '0' {
			b[i] = '0'
			k--
		}
	}

	fmt.Fprintln(writer, string(b))
}
