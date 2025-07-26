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
		var d int
		fmt.Fscan(reader, &n, &d)
		var s string
		fmt.Fscan(reader, &s)
		inserted := false
		for i := 0; i < n; i++ {
			if !inserted && int(s[i]-'0') < d {
				writer.WriteByte(byte('0' + d))
				inserted = true
			}
			writer.WriteByte(s[i])
		}
		if !inserted {
			writer.WriteByte(byte('0' + d))
		}
		writer.WriteByte('\n')
	}
}
