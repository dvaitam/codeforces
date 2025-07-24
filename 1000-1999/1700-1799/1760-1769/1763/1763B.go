package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Monster struct {
	h int64
	p int64
}

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
		var k int64
		fmt.Fscan(in, &n, &k)
		h := make([]int64, n)
		p := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}

		monsters := make([]Monster, n)
		for i := 0; i < n; i++ {
			monsters[i] = Monster{h: h[i], p: p[i]}
		}
		sort.Slice(monsters, func(i, j int) bool {
			return monsters[i].h < monsters[j].h
		})

		suf := make([]int64, n)
		suf[n-1] = monsters[n-1].p
		for i := n - 2; i >= 0; i-- {
			if monsters[i].p < suf[i+1] {
				suf[i] = monsters[i].p
			} else {
				suf[i] = suf[i+1]
			}
		}

		damage := int64(0)
		idx := 0
		cur := k
		for cur > 0 {
			damage += cur
			for idx < n && monsters[idx].h <= damage {
				idx++
			}
			if idx == n {
				break
			}
			cur -= suf[idx]
		}
		if idx == n {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
