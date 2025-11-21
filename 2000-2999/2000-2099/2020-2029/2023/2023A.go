package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	a, b int64
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
		pairs := make([]pair, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pairs[i].a, &pairs[i].b)
		}
		// Already read line by line? Input format? Actually first line n, followed by n lines each two ints but spec indicates maybe two arrays, but handle per pairs read.
		sort.Slice(pairs, func(i, j int) bool {
			si := pairs[i].a + pairs[i].b
			sj := pairs[j].a + pairs[j].b
			if si == sj {
				if pairs[i].a == pairs[j].a {
					return pairs[i].b < pairs[j].b
				}
				return pairs[i].a < pairs[j].a
			}
			return si < sj
		})
		res := make([]int64, 0, 2*n)
		for _, p := range pairs {
			res = append(res, p.a, p.b)
		}
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
