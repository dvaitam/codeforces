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
		var n, k int
		fmt.Fscan(in, &n, &k)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			freq[x]++
		}
		maxFreq := 0
		for _, c := range freq {
			if c > maxFreq {
				maxFreq = c
			}
		}
		if maxFreq >= k {
			fmt.Fprintln(out, k-1)
		} else {
			fmt.Fprintln(out, n)
		}
	}
}
