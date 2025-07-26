package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		digits := make([]int, n+1)
		for i := 0; i < n; i++ {
			digits[i+1] = int(s[i] - '0')
		}
		for i := n; i > 0; i-- {
			if digits[i] >= 5 {
				digits[i] = 0
				// propagate carry to the left
				j := i - 1
				digits[j]++
				for j > 0 && digits[j] == 10 {
					digits[j] = 0
					j--
					digits[j]++
				}
				for k := i + 1; k <= n; k++ {
					digits[k] = 0
				}
			}
		}
		start := 0
		if digits[0] == 0 {
			start = 1
		}
		var b strings.Builder
		for i := start; i <= n; i++ {
			b.WriteByte(byte(digits[i]) + '0')
		}
		fmt.Fprintln(writer, b.String())
	}
}
