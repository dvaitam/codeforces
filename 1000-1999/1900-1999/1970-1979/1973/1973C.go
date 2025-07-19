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
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		// separate by index parity
		ev := make([][2]int, 0, n)
		od := make([][2]int, 0, n)
		pos1 := 0
		for i, v := range a {
			if v == 1 {
				pos1 = i
			}
			if i&1 == 0 {
				ev = append(ev, [2]int{v, i})
			} else {
				od = append(od, [2]int{v, i})
			}
		}
		sort.Slice(ev, func(i, j int) bool { return ev[i][0] < ev[j][0] })
		sort.Slice(od, func(i, j int) bool { return od[i][0] < od[j][0] })
		b := make([]int, n)
		cur := n
		if pos1&1 == 1 {
			// starting with ev
			for _, x := range ev {
				b[x[1]] = cur
				cur--
			}
			for _, x := range od {
				b[x[1]] = cur
				cur--
			}
		} else {
			// starting with od
			for _, x := range od {
				b[x[1]] = cur
				cur--
			}
			for _, x := range ev {
				b[x[1]] = cur
				cur--
			}
		}
		// output
		for i, v := range b {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, v)
		}
		out.WriteByte('\n')
	}
}
