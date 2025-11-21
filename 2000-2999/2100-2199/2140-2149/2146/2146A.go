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

		freq := make([]int, n+1)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}

		maxFreq := 0
		for _, v := range freq {
			if v > maxFreq {
				maxFreq = v
			}
		}

		best := 0
		for c := 1; c <= maxFreq; c++ {
			cnt := 0
			for _, v := range freq {
				if v >= c {
					cnt++
				}
			}
			if cand := c * cnt; cand > best {
				best = cand
			}
		}

		fmt.Fprintln(writer, best)
	}
}
