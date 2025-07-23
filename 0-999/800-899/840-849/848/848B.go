package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type dancer struct {
	idx int
	g   int
	p   int
	t   int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, w, h int
	if _, err := fmt.Fscan(in, &n, &w, &h); err != nil {
		return
	}

	groups := make(map[int][]dancer)
	for i := 0; i < n; i++ {
		var g, p, t int
		fmt.Fscan(in, &g, &p, &t)
		key := p - t
		groups[key] = append(groups[key], dancer{idx: i, g: g, p: p, t: t})
	}

	ans := make([][2]int, n)
	for _, arr := range groups {
		horiz := make([]dancer, 0)
		vert := make([]dancer, 0)
		for _, d := range arr {
			if d.g == 2 {
				horiz = append(horiz, d)
			} else {
				vert = append(vert, d)
			}
		}
		sort.Slice(horiz, func(i, j int) bool { return horiz[i].p < horiz[j].p })
		sort.Slice(vert, func(i, j int) bool { return vert[i].p < vert[j].p })

		startOrder := append(append([]dancer{}, horiz...), vert...)

		fin := make([][2]int, 0, len(arr))
		tempV := append([]dancer{}, vert...)
		sort.Slice(tempV, func(i, j int) bool { return tempV[i].p < tempV[j].p })
		for _, d := range tempV {
			fin = append(fin, [2]int{d.p, h})
		}
		tempH := append([]dancer{}, horiz...)
		sort.Slice(tempH, func(i, j int) bool { return tempH[i].p < tempH[j].p })
		for _, d := range tempH {
			fin = append(fin, [2]int{w, d.p})
		}

		for i, d := range startOrder {
			ans[d.idx] = fin[i]
		}
	}

	for i := 0; i < n; i++ {
		fmt.Fprintf(out, "%d %d\n", ans[i][0], ans[i][1])
	}
}
