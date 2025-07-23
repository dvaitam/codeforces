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
	var t string
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	b := []byte(t)
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			reversePrefix(b, i)
		}
	}
	fmt.Fprintln(writer, string(b))
}

func reversePrefix(b []byte, k int) {
	for l, r := 0, k-1; l < r; l, r = l+1, r-1 {
		b[l], b[r] = b[r], b[l]
	}
}
