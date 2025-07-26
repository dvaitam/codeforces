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
		var s string
		fmt.Fscan(reader, &s)

		prefixA := make([]bool, n)
		hasA := false
		for i := 0; i < n; i++ {
			if s[i] == 'A' {
				hasA = true
			}
			prefixA[i] = hasA
		}

		suffixB := make([]bool, n+1)
		hasB := false
		for i := n - 1; i >= 0; i-- {
			if s[i] == 'B' {
				hasB = true
			}
			suffixB[i] = hasB
		}

		ans := 0
		for i := 0; i < n-1; i++ {
			if prefixA[i] && suffixB[i+1] {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
