package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Point struct{ x, y int }

func pointInPoly(poly []Point, p Point) bool {
	inside := false
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		if onSegment(a, b, p) {
			return true
		}
		yi, yj := a.y, b.y
		xi, xj := a.x, b.x
		if (yi > p.y) != (yj > p.y) {
			x := xi + (p.y-yi)*(xj-xi)/(yj-yi)
			if x >= p.x {
				inside = !inside
			}
		}
	}
	return inside
}

func onSegment(a, b, p Point) bool {
	dx1 := b.x - a.x
	dy1 := b.y - a.y
	dx2 := p.x - a.x
	dy2 := p.y - a.y
	if dx1*dy2 != dy1*dx2 {
		return false
	}
	if dx2 < 0 && dx1 > 0 || dx2 > 0 && dx1 < 0 {
		return false
	}
	if dy2 < 0 && dy1 > 0 || dy2 > 0 && dy1 < 0 {
		return false
	}
	if abs(dx2) > abs(dx1) || abs(dy2) > abs(dy1) {
		return false
	}
	return true
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func solve(points []Point, polys [][]int) []int {
	res := make([]int, len(polys))
	for i, polyIdx := range polys {
		poly := make([]Point, len(polyIdx))
		for j, id := range polyIdx {
			poly[j] = points[id]
		}
		// bounding box
		minx, maxx := poly[0].x, poly[0].x
		miny, maxy := poly[0].y, poly[0].y
		for _, p := range poly {
			if p.x < minx {
				minx = p.x
			}
			if p.x > maxx {
				maxx = p.x
			}
			if p.y < miny {
				miny = p.y
			}
			if p.y > maxy {
				maxy = p.y
			}
		}
		cnt := 0
		for _, pt := range points {
			if pt.x < minx || pt.x > maxx || pt.y < miny || pt.y > maxy {
				continue
			}
			if pointInPoly(poly, pt) {
				cnt++
			}
		}
		res[i] = cnt
	}
	return res
}

func genCase(rng *rand.Rand) ([]Point, [][]int) {
	n := rng.Intn(5) + 3
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		points[i] = Point{rng.Intn(11) - 5, rng.Intn(11) - 5}
	}
	q := 1 + rng.Intn(3)
	polys := make([][]int, q)
	for i := 0; i < q; i++ {
		k := rng.Intn(n-2) + 3
		polys[i] = rng.Perm(n)[:k]
	}
	return points, polys
}

func runCase(bin string, pts []Point, polys [][]int, expected []int) error {
	var sb strings.Builder
	n := len(pts)
	m := n // dummy edges count
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		sb.WriteString("1 1\n")
	}
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(polys)))
	for _, poly := range polys {
		sb.WriteString(fmt.Sprintf("%d", len(poly)))
		for _, id := range poly {
			sb.WriteString(fmt.Sprintf(" %d", id+1))
		}
		sb.WriteString("\n")
	}
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(expected) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		var x int
		fmt.Sscan(f, &x)
		if x != expected[i] {
			return fmt.Errorf("expected %d got %d", expected[i], x)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		pts, polys := genCase(rng)
		exp := solve(pts, polys)
		if err := runCase(bin, pts, polys, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
