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

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		seg := append([]int(nil), a[l:r]...)
		sort.Ints(seg)
		ans := seg[1] - seg[0]
		for i := 2; i < len(seg); i++ {
			d := seg[i] - seg[i-1]
			if d < ans {
				ans = d
			}
		}
		fmt.Fprintln(out, ans)
	}
}
