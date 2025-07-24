package main

import (
	"bufio"
	"fmt"
	"os"
)

// Operation represents a copy operation with source range [l,r]
// and the target range [start,end] in the growing string.
type Operation struct {
	l, r  int64
	start int64
	end   int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, c, q int
		fmt.Fscan(reader, &n, &c, &q)
		var s string
		fmt.Fscan(reader, &s)

		ops := make([]Operation, c)
		currLen := int64(n)
		for i := 0; i < c; i++ {
			var l, r int64
			fmt.Fscan(reader, &l, &r)
			ops[i] = Operation{l: l, r: r, start: currLen + 1, end: currLen + (r - l + 1)}
			currLen += r - l + 1
		}
		origLen := int64(n)
		for ; q > 0; q-- {
			var k int64
			fmt.Fscan(reader, &k)
			// Walk backwards through operations until k refers to the original string
			for k > origLen {
				for i := c - 1; i >= 0; i-- {
					op := ops[i]
					if k >= op.start && k <= op.end {
						k = op.l + (k - op.start)
						break
					}
				}
			}
			writer.WriteByte(s[k-1])
			writer.WriteByte('\n')
		}
	}
}
