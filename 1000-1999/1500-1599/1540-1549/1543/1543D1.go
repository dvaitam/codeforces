package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem 1543D1 - "RPD and Rap Sheet" (easy version).
// This problem is interactive. We repeatedly output the Gray code difference
// i^(i-1) so that after i steps the hidden value becomes x xor i. When our guess
// matches the current hidden value the judge responds with 1.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			return
		}
		for i := 0; i < n; i++ {
			var q int
			if i == 0 {
				q = 0
			} else {
				q = i ^ (i - 1)
			}
			fmt.Fprintln(writer, q)
			writer.Flush()
			var resp int
			if _, err := fmt.Fscan(reader, &resp); err != nil {
				return
			}
			if resp == 1 {
				break
			}
		}
	}
}
