package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Point struct{ x, y int64 }

func cross(a, b, c Point) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func convexHull(a []Point) []Point {
	n := len(a)
	if n <= 1 {
		return append([]Point{}, a...)
	}
	sort.Slice(a, func(i, j int) bool {
		if a[i].x != a[j].x {
			return a[i].x < a[j].x
		}
		return a[i].y < a[j].y
	})
	lo := []Point{}
	up := []Point{}
	for _, p := range a {
		for len(lo) >= 2 && cross(lo[len(lo)-2], lo[len(lo)-1], p) <= 0 {
			lo = lo[:len(lo)-1]
		}
		lo = append(lo, p)
	}
	for i := n - 1; i >= 0; i-- {
		p := a[i]
		for len(up) >= 2 && cross(up[len(up)-2], up[len(up)-1], p) <= 0 {
			up = up[:len(up)-1]
		}
		up = append(up, p)
	}
	lo = lo[:len(lo)-1]
	up = up[:len(up)-1]
	return append(lo, up...)
}

func pointInConvex(P []Point, q Point) bool {
	n := len(P)
	if n == 0 {
		return false
	}
	if cross(P[0], P[1], q) < 0 || cross(P[0], P[n-1], q) > 0 {
		return false
	}
	low, high := 1, n-1
	for high-low > 1 {
		mid := (low + high) / 2
		if cross(P[0], P[mid], q) >= 0 {
			low = mid
		} else {
			high = mid
		}
	}
	return cross(P[low], P[(low+1)%n], q) >= 0
}

func solve(input string) string {
	lines := strings.Fields(input)
	idx := 0
	q := toInt(lines[idx])
	idx++
	pts := make([]Point, 0)
	var outputs []string
	for i := 0; i < 3; i++ {
		t := toInt(lines[idx])
		idx++
		x := toInt64(lines[idx])
		idx++
		y := toInt64(lines[idx])
		idx++
		_ = t
		pts = append(pts, Point{x, y})
	}
	hull := convexHull(pts)
	for i := 3; i < q; i++ {
		t := toInt(lines[idx])
		idx++
		x := toInt64(lines[idx])
		idx++
		y := toInt64(lines[idx])
		idx++
		p := Point{x, y}
		if t == 1 {
			if pointInConvex(hull, p) {
				continue
			}
			n := len(hull)
			var L, R int
			if cross(hull[0], hull[1], p) < 0 {
				L, R = 0, 1
			} else if cross(hull[0], hull[n-1], p) > 0 {
				L, R = n-1, 0
			} else {
				low, high := 1, n-1
				for high-low > 1 {
					mid := (low + high) / 2
					if cross(hull[0], hull[mid], p) > 0 {
						low = mid
					} else {
						high = mid
					}
				}
				L, R = low, low+1
			}
			for cross(hull[R%len(hull)], hull[(R+1)%len(hull)], p) < 0 {
				R++
			}
			for cross(hull[(L-1+len(hull))%len(hull)], hull[L%len(hull)], p) < 0 {
				L--
			}
			newHull := []Point{p}
			n0 := len(hull)
			idx2 := R % n0
			end := L % n0
			for {
				newHull = append(newHull, hull[idx2])
				if idx2 == end {
					break
				}
				idx2 = (idx2 + 1) % n0
			}
			m := len(newHull)
			minI := 0
			for i := 1; i < m; i++ {
				if newHull[i].x < newHull[minI].x || (newHull[i].x == newHull[minI].x && newHull[i].y < newHull[minI].y) {
					minI = i
				}
			}
			hull = append(newHull[minI:], newHull[:minI]...)
		} else {
			if pointInConvex(hull, p) {
				outputs = append(outputs, "YES")
			} else {
				outputs = append(outputs, "NO")
			}
		}
	}
	return strings.Join(outputs, "\n")
}

func toInt(s string) int     { var x int; fmt.Sscan(s, &x); return x }
func toInt64(s string) int64 { var x int64; fmt.Sscan(s, &x); return x }

func genCase(rng *rand.Rand) string {
	q := rng.Intn(20) + 4
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	pts := make([]Point, 0)
	for len(pts) < 3 {
		x := int64(rng.Intn(21) - 10)
		y := int64(rng.Intn(21) - 10)
		p := Point{x, y}
		unique := true
		for _, t := range pts {
			if t.x == x && t.y == y {
				unique = false
				break
			}
		}
		if unique {
			pts = append(pts, p)
			sb.WriteString(fmt.Sprintf("1 %d %d\n", x, y))
		}
	}
	// ensure non-degenerate
	for cross(pts[0], pts[1], pts[2]) == 0 {
		pts[2].x = int64(rng.Intn(21) - 10)
		pts[2].y = int64(rng.Intn(21) - 10)
	}
	for i := 3; i < q; i++ {
		if rng.Intn(2) == 0 {
			// add unique point
			for {
				x := int64(rng.Intn(21) - 10)
				y := int64(rng.Intn(21) - 10)
				ok := true
				for _, t := range pts {
					if t.x == x && t.y == y {
						ok = false
						break
					}
				}
				if ok {
					pts = append(pts, Point{x, y})
					sb.WriteString(fmt.Sprintf("1 %d %d\n", x, y))
					break
				}
			}
		} else {
			x := int64(rng.Intn(21) - 10)
			y := int64(rng.Intn(21) - 10)
			sb.WriteString(fmt.Sprintf("2 %d %d\n", x, y))
		}
	}
	return sb.String()
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(expected) != got {
		return fmt.Errorf("expected\n%s\nbut got\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expected := solve(strings.TrimRight(input, "\n"))
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
