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
		var n int
		fmt.Fscan(in, &n)
		freq := make(map[int]int)
		maxFreq := 0
		for i := 0; i < n; i++ {
			var h int
			fmt.Fscan(in, &h)
			freq[h]++
			if freq[h] > maxFreq {
				maxFreq = freq[h]
			}
		}
		fmt.Fprintln(out, n-maxFreq)
	}
}
