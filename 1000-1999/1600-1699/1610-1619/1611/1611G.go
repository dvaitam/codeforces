package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type point struct{ u, v int }

func width(pts []point) int {
	if len(pts) == 0 {
		return 0
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].u == pts[j].u {
			return pts[i].v < pts[j].v
		}
		return pts[i].u < pts[j].u
	})
	tails := make([]int, 0, len(pts))
	for _, p := range pts {
		x := -p.v
		idx := sort.Search(len(tails), func(i int) bool { return tails[i] >= x })
		if idx == len(tails) {
			tails = append(tails, x)
		} else {
			tails[idx] = x
		}
	}
	return len(tails)
}

func solve(n, m int, grid []string) int {
	var lists [2][]point
	for i := 0; i < n; i++ {
		row := grid[i]
		for j := 0; j < m; j++ {
			if row[j] == '1' {
				u := i + j + 2
				v := i - j
				par := u & 1
				lists[par] = append(lists[par], point{u, v})
			}
		}
	}
	ans := 0
	for _, pts := range lists {
		ans += width(pts)
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &grid[i])
		}
		fmt.Fprintln(writer, solve(n, m, grid))
	}
}
