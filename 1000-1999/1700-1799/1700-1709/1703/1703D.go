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
		fmt.Fscan(reader, &n)
		arr := make([]string, n)
		set := make(map[string]bool, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			set[arr[i]] = true
		}
		for i := 0; i < n; i++ {
			s := arr[i]
			ok := false
			for j := 1; j < len(s); j++ {
				if set[s[:j]] && set[s[j:]] {
					ok = true
					break
				}
			}
			if ok {
				writer.WriteByte('1')
			} else {
				writer.WriteByte('0')
			}
		}
		writer.WriteByte('\n')
	}
}
