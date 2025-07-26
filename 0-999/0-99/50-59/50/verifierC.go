package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
)

type Point struct{ x, y int64 }

func cross(a, b, c Point) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func convexHull(pts []Point) []Point {
	n := len(pts)
	if n <= 1 {
		return pts
	}
	lower := make([]Point, 0, n)
	for _, p := range pts {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}
	upper := make([]Point, 0, n)
	for i := n - 1; i >= 0; i-- {
		p := pts[i]
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}
	hull := append(lower[:len(lower)-1], upper[:len(upper)-1]...)
	return hull
}

func expected(pts []Point) int64 {
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x != pts[j].x {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	uniq := pts[:0]
	for i, p := range pts {
		if i == 0 || p.x != pts[i-1].x || p.y != pts[i-1].y {
			uniq = append(uniq, p)
		}
	}
	pts = uniq
	m := len(pts)
	if m == 0 {
		return 0
	}
	if m == 1 {
		return 4
	}
	hull := convexHull(pts)
	hsz := len(hull)
	if hsz == 1 {
		return 4
	}
	if hsz == 2 {
		dx := abs64(hull[1].x - hull[0].x)
		dy := abs64(hull[1].y - hull[0].y)
		d := dx
		if dy > d {
			d = dy
		}
		return 2*d + 4
	}
	var perim int64
	for i := 0; i < hsz; i++ {
		j := (i + 1) % hsz
		dx := abs64(hull[j].x - hull[i].x)
		dy := abs64(hull[j].y - hull[i].y)
		if dx > dy {
			perim += dx
		} else {
			perim += dy
		}
	}
	return perim
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(43))
	const cases = 120
	for t := 0; t < cases; t++ {
		n := r.Intn(10) + 1
		pts := make([]Point, n)
		for i := 0; i < n; i++ {
			pts[i] = Point{int64(r.Intn(21) - 10), int64(r.Intn(21) - 10)}
		}
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for _, p := range pts {
			fmt.Fprintf(&input, "%d %d\n", p.x, p.y)
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Fscan(&out, &got); err != nil {
			fmt.Printf("Test %d parse error: %v\n", t+1, err)
			os.Exit(1)
		}
		want := expected(append([]Point(nil), pts...))
		if got != want {
			fmt.Printf("Test %d failed: expected %d got %d\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
