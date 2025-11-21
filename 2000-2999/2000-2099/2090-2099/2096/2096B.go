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
		l := make([]int64, n)
		r := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &r[i])
		}
		mins := make([]int64, n)
		var base int64
		for i := 0; i < n; i++ {
			if l[i] > r[i] {
				base += l[i]
				mins[i] = r[i]
			} else {
				base += r[i]
				mins[i] = l[i]
			}
		}
		sort.Slice(mins, func(i, j int) bool { return mins[i] > mins[j] })
		var extra int64
		for i := 0; i < k-1 && i < n; i++ {
			extra += mins[i]
		}
		ans := base + extra + 1
		fmt.Fprintln(out, ans)
	}
}
