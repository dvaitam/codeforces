package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problem described in problemD.txt of folder 1476.
// Computes the maximum number of distinct cities a traveler can visit
// starting from each city when road directions flip after every move.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		s = " " + s

		left := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if s[i] == 'L' {
				left[i] = 1
				if i >= 2 && s[i-1] == 'R' {
					left[i] = 2 + left[i-2]
				}
			}
		}

		right := make([]int, n+2)
		for i := n; i >= 1; i-- {
			if s[i] == 'R' {
				right[i] = 1
				if i+1 <= n && s[i+1] == 'L' {
					right[i] = 2 + right[i+2]
				}
			}
		}

		ans := make([]int, n+1)
		ans[0] = 1 + right[1]
		for i := 1; i < n; i++ {
			ans[i] = 1 + left[i] + right[i+1]
		}
		ans[n] = 1 + left[n]

		for i := 0; i <= n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, ans[i])
		}
		writer.WriteByte('\n')
	}
}
