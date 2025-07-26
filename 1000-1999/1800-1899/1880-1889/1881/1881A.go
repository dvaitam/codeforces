package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// This program solves problem A from contest 1881.
// We repeatedly double the string x until the pattern s becomes
// a substring or we have checked enough repetitions. The length of
// strings is tiny (n*m <= 25), so it is sufficient to stop once
// the doubled string exceeds len(s)+len(x). If s is not found by
// then, it will never appear in further doublings.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var x, s string
		fmt.Fscan(reader, &x)
		fmt.Fscan(reader, &s)

		cur := x
		ops := 0
		limit := len(s) + len(x)
		ans := -1
		for {
			if strings.Contains(cur, s) {
				ans = ops
				break
			}
			if len(cur) > limit {
				break
			}
			cur += cur
			ops++
		}
		fmt.Fprintln(writer, ans)
	}
}
