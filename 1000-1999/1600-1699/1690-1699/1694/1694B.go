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
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		var res int64 = 1
		for i := 1; i < n; i++ {
			if s[i] != s[i-1] {
				res += int64(i + 1)
			} else {
				res++
			}
		}
		fmt.Fprintln(writer, res)
	}
}
