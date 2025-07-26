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
		var s string
		fmt.Fscan(reader, &n, &s)
		maxRun := 1
		cur := 1
		for i := 1; i < n; i++ {
			if s[i] == s[i-1] {
				cur++
			} else {
				if cur > maxRun {
					maxRun = cur
				}
				cur = 1
			}
		}
		if cur > maxRun {
			maxRun = cur
		}
		fmt.Fprintln(writer, maxRun+1)
	}
}
