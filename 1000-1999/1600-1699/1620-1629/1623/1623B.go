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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		type interval struct{ l, r, ans int }
		ivs := make([]interval, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &ivs[i].l, &ivs[i].r)
		}
		// sort by interval length ascending
		sort.Slice(ivs, func(i, j int) bool {
			return (ivs[i].r - ivs[i].l) < (ivs[j].r - ivs[j].l)
		})
		// assign points
		flag := make([]bool, n+2)
		for i := 0; i < n; i++ {
			for k := ivs[i].l; k <= ivs[i].r; k++ {
				if k >= 0 && k < len(flag) && !flag[k] {
					ivs[i].ans = k
					flag[k] = true
					break
				}
			}
		}
		// output
		for i := 0; i < n; i++ {
			fmt.Fprintf(writer, "%d %d %d\n", ivs[i].l, ivs[i].r, ivs[i].ans)
		}
		fmt.Fprintln(writer)
	}
}
