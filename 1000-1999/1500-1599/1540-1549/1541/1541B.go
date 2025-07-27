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
		pairs := make([]struct{ val, idx int }, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pairs[i].val)
			pairs[i].idx = i + 1
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
		ans := 0
		limit := 2 * n
		for i := 0; i < n; i++ {
			for j := i + 1; j < n && pairs[i].val*pairs[j].val <= limit; j++ {
				if pairs[i].val*pairs[j].val == pairs[i].idx+pairs[j].idx {
					ans++
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
