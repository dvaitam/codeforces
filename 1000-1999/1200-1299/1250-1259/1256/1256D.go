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

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)

		zeroIdx := make([]int, 0)
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				zeroIdx = append(zeroIdx, i)
			}
		}
		res := make([]byte, n)
		for i := 0; i < n; i++ {
			res[i] = '1'
		}
		pos := 0
		for _, idx := range zeroIdx {
			dist := idx - pos
			shift := int64(dist)
			if shift > k {
				shift = k
			}
			final := idx - int(shift)
			res[final] = '0'
			k -= shift
			pos++
		}
		fmt.Fprintln(writer, string(res))
	}
}
