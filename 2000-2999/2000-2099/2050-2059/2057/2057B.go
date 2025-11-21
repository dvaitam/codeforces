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
		counts := make([]int, 0, len(freq))
		for _, v := range freq {
			counts = append(counts, v)
		}
		sort.Ints(counts)
		distinct := len(counts)
		for _, c := range counts {
			if k >= c {
				k -= c
				distinct--
			} else {
				break
			}
		}
		if distinct == 0 {
			distinct = 1
		}
		fmt.Fprintln(out, distinct)
	}
}
