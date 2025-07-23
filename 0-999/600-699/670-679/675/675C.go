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
	freq := make(map[int64]int)
	var prefix int64
	maxFreq := 0
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		prefix += x
		freq[prefix]++
		if freq[prefix] > maxFreq {
			maxFreq = freq[prefix]
		}
	}
	fmt.Fprintln(out, n-maxFreq)
}
