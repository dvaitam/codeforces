package main

import (
	"bufio"
	"fmt"
	"os"
)

type target struct {
	x int64
	y int64
	d int64
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func canonical(x, y int64) target {
	ox, oy := x, y
	if y < 0 {
		y = -y
		x = -x
	}
	if y == 0 && x < 0 {
		x = -x
	} else if y == 0 && x == 0 {
		if ox < 0 {
			x = -x
		}
		if oy < 0 {
			y = -y
		}
	}
	return target{x: x, y: y, d: abs64(ox) + abs64(oy)}
}

func pointsOnManhattanCircle(d int64) []target {
	points := make([]target, 0, 4*d)
	for x := int64(0); x <= d; x++ {
		y := d - x
		points = append(points, canonical(x, y))
		points = append(points, canonical(-x, y))
		points = append(points, canonical(x, -y))
		points = append(points, canonical(-x, -y))
	}
	return points
}

func canonicalList(points []target) []target {
	unique := make(map[[3]int64]struct{})
	list := make([]target, 0, len(points))
	for _, p := range points {
		key := [3]int64{p.x, p.y, p.d}
		if _, ok := unique[key]; !ok {
			unique[key] = struct{}{}
			list = append(list, p)
		}
	}
	return list
}

var (
	n   int
	a   int64
	b   int64
	d   []int64
	pos []target
)

func backtrack(idx int, curX, curY int64, used map[[2]int64]struct{}, coords [][2]int64) bool {
	if idx == n {
		if curX == a && curY == b {
			fmt.Println("YES")
			for _, c := range coords {
				fmt.Println(c[0], c[1])
			}
			return true
		}
		return false
	}
	for _, p := range pos[idx] {
		nx := curX + p.x
		ny := curY + p.y
		key := [2]int64{nx, ny}
		if _, ok := used[key]; ok {
			continue
		}
		rem := int64(0)
		for i := idx + 1; i < n; i++ {
			rem += d[i]
		}
		if abs64(nx-a)+abs64(ny-b) > rem {
			continue
		}
		used[key] = struct{}{}
		coords[idx] = key
		if backtrack(idx+1, nx, ny, used, coords) {
			return true
		}
		delete(used, key)
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	fmt.Fscan(reader, &a, &b)
	d = make([]int64, n)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &d[i])
	}
	d[n-1] = 0

	pos = make([][]target, n)
	pos[1] = canonicalList(pointsOnManhattanCircle(d[0]))
	for i := 2; i < n; i++ {
		pos[i] = make([]target, 0)
		for _, p := range pos[i-1] {
			nextTargets := canonicalList(pointsOnManhattanCircle(d[i-1]))
			for _, q := range nextTargets {
				pos[i] = append(pos[i], target{x: p.x + q.x, y: p.y + q.y, d: p.d + q.d})
			}
		}
		pos[i] = canonicalList(pos[i])
	}

	coords := make([][2]int64, n)
	coords[0] = [2]int64{0, 0}
	used := map[[2]int64]struct{}{{0, 0}: {}}
	if !backtrack(1, 0, 0, used, coords) {
		fmt.Println("NO")
	}
}
