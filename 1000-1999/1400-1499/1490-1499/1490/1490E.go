package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	val int64
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		pairs := make([]Pair, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pairs[i].val)
			pairs[i].idx = i + 1
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
		prefix := make([]int64, n)
		prefix[0] = pairs[0].val
		for i := 1; i < n; i++ {
			prefix[i] = prefix[i-1] + pairs[i].val
		}
		pos := n - 1
		for i := n - 2; i >= 0; i-- {
			if prefix[i] < pairs[i+1].val {
				break
			}
			pos = i
		}
		winners := make([]int, 0, n-pos)
		for i := pos; i < n; i++ {
			winners = append(winners, pairs[i].idx)
		}
		sort.Ints(winners)
		fmt.Fprintln(out, len(winners))
		for i, id := range winners {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, id)
		}
		fmt.Fprintln(out)
	}
}
