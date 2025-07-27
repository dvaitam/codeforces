package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ x, y int64 }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		if n%2 == 1 && a[n/2] != b[n/2] {
			fmt.Fprintln(writer, "No")
			continue
		}
		m := n / 2
		pa := make([]pair, m)
		pb := make([]pair, m)
		for i := 0; i < m; i++ {
			x1, x2 := a[i], a[n-1-i]
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			pa[i] = pair{x1, x2}
			y1, y2 := b[i], b[n-1-i]
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			pb[i] = pair{y1, y2}
		}
		sort.Slice(pa, func(i, j int) bool {
			if pa[i].x == pa[j].x {
				return pa[i].y < pa[j].y
			}
			return pa[i].x < pa[j].x
		})
		sort.Slice(pb, func(i, j int) bool {
			if pb[i].x == pb[j].x {
				return pb[i].y < pb[j].y
			}
			return pb[i].x < pb[j].x
		})
		ok := true
		for i := 0; i < m; i++ {
			if pa[i] != pb[i] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
