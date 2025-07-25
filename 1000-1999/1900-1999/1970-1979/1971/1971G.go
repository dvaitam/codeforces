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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		groups := make(map[int][]int)
		vals := make(map[int][]int)
		for i, v := range a {
			g := v >> 2
			groups[g] = append(groups[g], i)
			vals[g] = append(vals[g], v)
		}
		for g := range groups {
			idxs := groups[g]
			vs := vals[g]
			sort.Ints(vs)
			for j, idx := range idxs {
				a[idx] = vs[j]
			}
		}
		for i, v := range a {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
