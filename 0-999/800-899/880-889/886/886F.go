package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct {
	x, y int64
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isGood(dx, dy int64, pts []Point, sumX, sumY int64, n int) bool {
	m := make(map[int64]int, n)
	target := sumX*dx + sumY*dy
	for _, p := range pts {
		val := (p.x*dx+p.y*dy)*int64(n) - target
		m[val]++
	}
	for k, v := range m {
		if m[-k] != v {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pts := make([]Point, n)
	var sumX, sumY int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
		sumX += pts[i].x
		sumY += pts[i].y
	}

	// check global central symmetry
	sorted := make([]Point, n)
	copy(sorted, pts)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].x == sorted[j].x {
			return sorted[i].y < sorted[j].y
		}
		return sorted[i].x < sorted[j].x
	})
	candX := sorted[0].x + sorted[n-1].x
	candY := sorted[0].y + sorted[n-1].y
	symmetric := true
	for i := 0; i < n; i++ {
		if sorted[i].x+sorted[n-1-i].x != candX || sorted[i].y+sorted[n-1-i].y != candY {
			symmetric = false
			break
		}
	}
	if symmetric {
		fmt.Fprintln(out, -1)
		return
	}

	type pair struct{ x, y int64 }
	dirs := make(map[pair]struct{})
	idx := []int{0, 1}
	if n == 1 {
		// already handled symmetric case above, but just in case
		fmt.Fprintln(out, 0)
		return
	}
	for _, i := range idx {
		if i >= n {
			continue
		}
		for j := i + 1; j < n; j++ {
			ux := (pts[i].x+pts[j].x)*int64(n) - 2*sumX
			uy := (pts[i].y+pts[j].y)*int64(n) - 2*sumY
			if ux == 0 && uy == 0 {
				continue
			}
			dx := -uy
			dy := ux
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			dirs[pair{dx, dy}] = struct{}{}
		}
	}

	count := 0
	for d := range dirs {
		if isGood(d.x, d.y, pts, sumX, sumY, n) {
			count++
		}
	}
	fmt.Fprintln(out, count)
}
