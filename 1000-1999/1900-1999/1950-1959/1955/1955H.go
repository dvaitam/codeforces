package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Point struct{ x, y int }

func ceilSqrt(x int) int {
	if x <= 0 {
		return 0
	}
	r := int(math.Sqrt(float64(x)))
	for r*r < x {
		r++
	}
	return r
}

func towerCoverage(path []Point, tx, ty, R int) []int {
	freq := make([]int, R+2)
	for _, p := range path {
		dx := p.x - tx
		dy := p.y - ty
		d2 := dx*dx + dy*dy
		r := ceilSqrt(d2)
		if r <= R {
			freq[r]++
		}
	}
	cov := make([]int, R+1)
	sum := 0
	for r := 1; r <= R; r++ {
		sum += freq[r]
		cov[r] = sum
	}
	return cov
}

func hungarianMax(profit [][]int) int {
	n := len(profit)
	if n == 0 {
		return 0
	}
	m := len(profit[0])
	const INF = int(1e15)
	u := make([]int, n+1)
	v := make([]int, m+1)
	p := make([]int, m+1)
	way := make([]int, m+1)
	for i := 1; i <= n; i++ {
		p[0] = i
		j0 := 0
		minv := make([]int, m+1)
		used := make([]bool, m+1)
		for j := 0; j <= m; j++ {
			minv[j] = INF
		}
		for {
			used[j0] = true
			i0 := p[j0]
			delta := INF
			j1 := 0
			for j := 1; j <= m; j++ {
				if !used[j] {
					cur := -profit[i0-1][j-1] - u[i0] - v[j]
					if cur < minv[j] {
						minv[j] = cur
						way[j] = j0
					}
					if minv[j] < delta {
						delta = minv[j]
						j1 = j
					}
				}
			}
			for j := 0; j <= m; j++ {
				if used[j] {
					u[p[j]] += delta
					v[j] -= delta
				} else {
					minv[j] -= delta
				}
			}
			j0 = j1
			if p[j0] == 0 {
				break
			}
		}
		for {
			j1 := way[j0]
			p[j0] = p[j1]
			j0 = j1
			if j0 == 0 {
				break
			}
		}
	}
	assign := make([]int, n)
	for j := 1; j <= m; j++ {
		if p[j] > 0 && p[j]-1 < n {
			assign[p[j]-1] = j - 1
		}
	}
	res := 0
	for i := 0; i < n; i++ {
		res += profit[i][assign[i]]
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	const R = 19
	pow3 := make([]int, R+1)
	pow3[0] = 1
	for i := 1; i <= R; i++ {
		pow3[i] = pow3[i-1] * 3
	}

	for ; T > 0; T-- {
		var n, m, k int
		if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
			return
		}
		path := make([]Point, 0, n*m)
		for i := 1; i <= n; i++ {
			var row string
			fmt.Fscan(in, &row)
			for j := 1; j <= m; j++ {
				if row[j-1] == '#' {
					path = append(path, Point{i, j})
				}
			}
		}
		towers := make([]struct{ x, y, p int }, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &towers[i].x, &towers[i].y, &towers[i].p)
		}
		if len(path) == 0 || k == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		profits := make([][]int, R)
		for r := 0; r < R; r++ {
			profits[r] = make([]int, k+R)
		}
		for j, t := range towers {
			cov := towerCoverage(path, t.x, t.y, R)
			for r := 1; r <= R; r++ {
				profits[r-1][j] = cov[r]*t.p - pow3[r]
			}
		}
		maxProf := hungarianMax(profits)
		if maxProf > 0 {
			fmt.Fprintln(out, maxProf)
		} else {
			fmt.Fprintln(out, 0)
		}
	}
}
