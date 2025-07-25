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
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(reader, &n)
		oddLen := 0
		totalOnes := 0
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			if len(s)%2 == 1 {
				oddLen++
			}
			for j := 0; j < len(s); j++ {
				if s[j] == '1' {
					totalOnes++
				}
			}
		}
		if oddLen > 0 || totalOnes%2 == 0 {
			fmt.Fprintln(writer, n)
		} else {
			fmt.Fprintln(writer, n-1)
		}
	}
}
