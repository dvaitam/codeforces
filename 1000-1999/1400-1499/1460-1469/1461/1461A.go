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
	fmt.Fscan(reader, &t)
	for i := 0; i < t; i++ {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		res := make([]byte, n)
		// fill first k with 'a'
		for j := 0; j < k && j < n; j++ {
			res[j] = 'a'
		}
		// cycle through 'b','c','a'
		pattern := []byte{'b', 'c', 'a'}
		for j := k; j < n; j++ {
			res[j] = pattern[(j-k)%3]
		}
		writer.Write(res)
		writer.WriteByte('\n')
	}
}
