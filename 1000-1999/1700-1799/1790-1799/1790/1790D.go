package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		freq := make(map[int]int, n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}
		keys := make([]int, 0, len(freq))
		for k := range freq {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		total := 0
		active := 0
		prev := 0
		first := true
		for _, k := range keys {
			if first {
				first = false
			} else if k != prev+1 {
				active = 0
			}
			need := freq[k] - active
			if need > 0 {
				total += need
			}
			active = freq[k]
			prev = k
		}
		fmt.Fprintln(writer, total)
	}
}
