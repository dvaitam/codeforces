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
	freq := make(map[int]int)
	pair := make(map[[2]int]int)
	for i := 0; i < n; i++ {
		var w, h int
		fmt.Fscan(in, &w, &h)
		if w == h {
			freq[w]++
		} else {
			freq[w]++
			freq[h]++
			if w < h {
				pair[[2]int{w, h}]++
			} else {
				pair[[2]int{h, w}]++
			}
		}
	}
	var total int64
	for _, c := range freq {
		total += int64(c) * int64(c-1) / 2
	}
	var dup int64
	for _, c := range pair {
		dup += int64(c) * int64(c-1) / 2
	}
	fmt.Fprintln(out, total-dup)
}
