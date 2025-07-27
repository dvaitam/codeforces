package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Curse struct {
	tl, tr int64
	l, r   int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var x int64
	fmt.Fscan(in, &x)
	curses := make([]Curse, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &curses[i].tl, &curses[i].tr, &curses[i].l, &curses[i].r)
	}
	sort.Slice(curses, func(i, j int) bool {
		if curses[i].tl == curses[j].tl {
			return curses[i].tr < curses[j].tr
		}
		return curses[i].tl < curses[j].tl
	})
	var pos int64 = x
	var t int64 = 0
	energy := 0.0
	for _, c := range curses {
		// advance time to c.tl
		pos += c.tl - t
		t = c.tl
		// if inside curse, move to nearest boundary
		if pos >= c.l && pos <= c.r {
			distLeft := float64(pos - c.l)
			distRight := float64(c.r - pos)
			if distLeft < distRight {
				energy += distLeft
				pos = c.l
			} else {
				energy += distRight
				pos = c.r
			}
		}
		// move till end of curse
		pos += c.tr - t
		t = c.tr
	}
	// output ceil of energy
	ans := int64(energy)
	if float64(ans) < energy {
		ans++
	}
	fmt.Println(ans)
}
