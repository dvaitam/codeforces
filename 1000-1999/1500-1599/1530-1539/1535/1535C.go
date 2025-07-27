package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		dp0, dp1 := 0, 0
		var total int64
		for i := 0; i < len(s); i++ {
			ch := s[i]
			expected0, expected1 := byte('0'), byte('1')
			if i%2 == 1 {
				expected0, expected1 = '1', '0'
			}
			if ch == '?' || ch == expected0 {
				dp0++
			} else {
				dp0 = 0
			}
			if ch == '?' || ch == expected1 {
				dp1++
			} else {
				dp1 = 0
			}
			if dp0 > dp1 {
				total += int64(dp0)
			} else {
				total += int64(dp1)
			}
		}
		fmt.Fprintln(out, total)
	}
}
