package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pixel struct {
	x, y int
	t    int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k, q int
	if _, err := fmt.Fscan(in, &n, &m, &k, &q); err != nil {
		return
	}
	px := make([]pixel, q)
	times := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &px[i].x, &px[i].y, &px[i].t)
		px[i].x--
		px[i].y--
		times[i] = px[i].t
	}

	if k == 0 {
		fmt.Fprintln(out, 0)
		return
	}
	if q == 0 {
		fmt.Fprintln(out, -1)
		return
	}

	sort.Ints(times)
	uniq := make([]int, 0, len(times))
	for _, v := range times {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}

	need := k * k
	check := func(T int) bool {
		grid := make([][]int, n)
		for i := range grid {
			grid[i] = make([]int, m)
		}
		for _, p := range px {
			if p.t <= T {
				grid[p.x][p.y] = 1
			}
		}
		pref := make([][]int, n+1)
		for i := range pref {
			pref[i] = make([]int, m+1)
		}
		for i := 1; i <= n; i++ {
			row := grid[i-1]
			prow := pref[i-1]
			srow := pref[i]
			for j := 1; j <= m; j++ {
				srow[j] = srow[j-1] + prow[j] - prow[j-1] + row[j-1]
			}
		}
		for i := k; i <= n; i++ {
			for j := k; j <= m; j++ {
				cnt := pref[i][j] - pref[i-k][j] - pref[i][j-k] + pref[i-k][j-k]
				if cnt == need {
					return true
				}
			}
		}
		return false
	}

	ans := -1
	l, r := 0, len(uniq)-1
	for l <= r {
		mid := (l + r) / 2
		if check(uniq[mid]) {
			ans = uniq[mid]
			r = mid - 1
		} else {
			l = mid + 1
		}
	}

	fmt.Fprintln(out, ans)
}
