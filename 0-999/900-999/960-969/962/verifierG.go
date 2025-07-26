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
type Rect struct{ x1, y1, x2, y2 int }

func insideRect(p Point, r Rect) bool {
	return p.x > r.x1 && p.x < r.x2 && p.y > r.y2 && p.y < r.y1
}

func pointInPolygon(pt Point, poly []Point) bool {
	inside := false
	n := len(poly)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := poly[i].x, poly[i].y
		xj, yj := poly[j].x, poly[j].y
		if (yi > pt.y) != (yj > pt.y) {
			if int64(pt.x-xi) < int64(xj-xi)*int64(pt.y-yi)/int64(yj-yi) {
				inside = !inside
			}
		}
	}
	return inside
}

func solve(r Rect, poly []Point) int {
	cross := 0
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		ia := insideRect(a, r)
		ib := insideRect(b, r)
		if ia != ib {
			cross++
			continue
		}
		if ia || ib {
			continue
		}
		if a.x == b.x {
			x := a.x
			if x > r.x1 && x < r.x2 {
				low, high := a.y, b.y
				if low > high {
					low, high = high, low
				}
				if low < r.y2 && high > r.y2 {
					cross++
				}
				if low < r.y1 && high > r.y1 {
					cross++
				}
			}
		} else if a.y == b.y {
			y := a.y
			if y > r.y2 && y < r.y1 {
				low, high := a.x, b.x
				if low > high {
					low, high = high, low
				}
				if low < r.x1 && high > r.x1 {
					cross++
				}
				if low < r.x2 && high > r.x2 {
					cross++
				}
			}
		}
	}
	if cross > 0 {
		return cross / 2
	}
	center := Point{(r.x1 + r.x2) / 2, (r.y1 + r.y2) / 2}
	if pointInPolygon(center, poly) || insideRect(poly[0], r) {
		return 1
	}
	return 0
}

func genCase(rng *rand.Rand) (string, string) {
	x1 := rng.Intn(20)
	x2 := x1 + rng.Intn(5) + 1
	y2 := rng.Intn(20)
	y1 := y2 + rng.Intn(5) + 1
	r := Rect{x1, y1, x2, y2}
	// polygon as rectangle
	px1 := rng.Intn(20)
	py2 := rng.Intn(20)
	px2 := px1 + rng.Intn(5) + 1
	py1 := py2 + rng.Intn(5) + 1
	poly := []Point{{px1, py1}, {px2, py1}, {px2, py2}, {px1, py2}}
	n := len(poly)
	input := fmt.Sprintf("%d %d %d %d\n%d\n", r.x1, r.y1, r.x2, r.y2, n)
	for _, p := range poly {
		input += fmt.Sprintf("%d %d\n", p.x, p.y)
	}
	out := solve(r, poly)
	expected := fmt.Sprintf("%d\n", out)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
