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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	freq := make([]int, 1001)
	maxFreq := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		freq[x]++
		if freq[x] > maxFreq {
			maxFreq = freq[x]
		}
	}
	fmt.Fprintln(out, n-maxFreq)
}
