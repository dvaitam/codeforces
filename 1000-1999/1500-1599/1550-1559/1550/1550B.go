package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, a, b int
		fmt.Fscan(reader, &n, &a, &b)
		var s string
		fmt.Fscan(reader, &s)

		if b >= 0 {
			fmt.Fprintln(writer, (a+b)*n)
			continue
		}

		zeroGroups := 0
		oneGroups := 0
		prev := byte(0)
		for i := 0; i < n; i++ {
			if i == 0 || s[i] != prev {
				if s[i] == '0' {
					zeroGroups++
				} else {
					oneGroups++
				}
				prev = s[i]
			}
		}
		ops := 1 + min(zeroGroups, oneGroups)
		ans := a*n + b*ops
		fmt.Fprintln(writer, ans)
	}
}
