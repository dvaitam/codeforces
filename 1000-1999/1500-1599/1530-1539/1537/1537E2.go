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
	fmt.Fscan(reader, &n, &k)
	var s string
	fmt.Fscan(reader, &s)

	best := 1
	for i := 1; i < n; i++ {
		if s[i] > s[i%best] {
			break
		}
		if s[i] < s[i%best] {
			best = i + 1
		}
	}

	prefix := s[:best]
	res := make([]byte, k)
	for i := 0; i < k; i++ {
		res[i] = prefix[i%best]
	}
	fmt.Fprintln(writer, string(res))
}
