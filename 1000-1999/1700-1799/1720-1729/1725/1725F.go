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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	L := make([]int64, n)
	R := make([]int64, n)
	lenArr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &L[i], &R[i])
		lenArr[i] = R[i] - L[i] + 1
	}

	const maxT = 30
	ans := make([]int, maxT)
	for t := 0; t < maxT; t++ {
		g := int64(1) << uint(t)
		if g > 1e9 {
			break
		}
		base := 0
		type ev struct {
			pos   int64
			delta int
		}
		events := make([]ev, 0, 4*n)
		for i := 0; i < n; i++ {
			length := lenArr[i]
			if length >= g {
				base++
				continue
			}
			start := L[i] % g
			end := (L[i] + length - 1) % g
			if start+length-1 < g {
				events = append(events, ev{start, 1})
				events = append(events, ev{end + 1, -1})
			} else {
				events = append(events, ev{0, 1})
				events = append(events, ev{end + 1, -1})
				events = append(events, ev{start, 1})
				events = append(events, ev{g, -1})
			}
		}
		sort.Slice(events, func(i, j int) bool { return events[i].pos < events[j].pos })
		curr, maxp := 0, 0
		i := 0
		for i < len(events) {
			pos := events[i].pos
			for i < len(events) && events[i].pos == pos {
				curr += events[i].delta
				i++
			}
			if curr > maxp {
				maxp = curr
			}
		}
		ans[t] = base + maxp
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var W int64
		fmt.Fscan(in, &W)
		g := W & -W
		t := 0
		for (int64(1) << uint(t)) != g {
			t++
		}
		fmt.Fprintln(out, ans[t])
	}
}
