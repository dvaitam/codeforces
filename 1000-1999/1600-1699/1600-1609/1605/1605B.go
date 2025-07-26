package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Solution to problemB.txt for 1605B (Reverse Sort).
// We compare the string with its sorted version. If already sorted,
// no operations are needed. Otherwise we output one operation that
// reverses all mismatched positions, which form a non-increasing subsequence.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n, &s)
		zeros := strings.Count(s, "0")
		if zeros == 0 || zeros == n {
			fmt.Fprintln(out, 0)
			continue
		}
		target := strings.Repeat("0", zeros) + strings.Repeat("1", n-zeros)
		if s == target {
			fmt.Fprintln(out, 0)
			continue
		}
		indices := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if s[i] != target[i] {
				indices = append(indices, i+1)
			}
		}
		fmt.Fprintln(out, 1)
		fmt.Fprint(out, len(indices))
		for _, idx := range indices {
			fmt.Fprint(out, " ", idx)
		}
		fmt.Fprintln(out)
	}
}
