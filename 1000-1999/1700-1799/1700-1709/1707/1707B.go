package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		freq := make(map[int]int, n)
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			freq[v]++
		}
		total := n
		for total > 1 {
			keys := make([]int, 0, len(freq))
			for k := range freq {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			newFreq := make(map[int]int, len(keys))
			zeros := 0
			prev := 0
			havePrev := false
			for _, k := range keys {
				c := freq[k]
				zeros += c - 1
				if havePrev {
					newFreq[k-prev]++
				}
				prev = k
				havePrev = true
			}
			if zeros > 0 {
				newFreq[0] += zeros
			}
			freq = newFreq
			total--
		}
		for ans := range freq {
			fmt.Fprintln(out, ans)
			break
		}
	}
}
