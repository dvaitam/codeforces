package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for problemB.txt from contest 1849.
// We repeatedly deal k damage to the monster with the highest current health
// (breaking ties by smaller index). This process kills monsters in the order of
// decreasing a[i] % k (treat 0 as k). We sort monsters by this value and
// output their 1-based indices.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		type pair struct{ r, idx int }
		b := make([]pair, n)
		for i := 0; i < n; i++ {
			r := a[i] % k
			if r == 0 {
				r = k
			}
			b[i] = pair{r, i + 1}
		}
		sort.Slice(b, func(i, j int) bool {
			if b[i].r == b[j].r {
				return b[i].idx < b[j].idx
			}
			return b[i].r > b[j].r
		})
		for i, p := range b {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, p.idx)
		}
		fmt.Fprintln(out)
	}
}
