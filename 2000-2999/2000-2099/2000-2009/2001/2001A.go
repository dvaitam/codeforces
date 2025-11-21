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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		freq := make(map[int]int)
		maxFreq := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			freq[x]++
			if freq[x] > maxFreq {
				maxFreq = freq[x]
			}
		}

		// We can keep at most one value per remaining element, so the best scenario
		// is to keep all occurrences of the most frequent number and delete the rest.
		// Thus the minimum number of operations equals the number of deletions needed.
		fmt.Fprintln(out, n-maxFreq)
	}
}
