package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type point struct{ x, y float64 }

func cross(a, b, c point) float64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func areaPoly(p []point) float64 {
	s := 0.0
	n := len(p)
	for i := 1; i < n-1; i++ {
		s += cross(p[0], p[i], p[i+1])
	}
	return math.Abs(s) * 0.5
}

func areaTri(p1, p2, p3 point) float64 {
	return math.Abs(cross(p1, p2, p3)) * 0.5
}

func convexHull(pts []point) []point {
	n := len(pts)
	if n <= 1 {
		return append([]point(nil), pts...)
	}
	sort.Slice(pts, func(i, j int) bool {
		if math.Abs(pts[i].y-pts[j].y) > 1e-9 {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	lower := make([]point, 0, n)
	for _, p := range pts {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}
	upper := make([]point, 0, n)
	for i := n - 1; i >= 0; i-- {
		p := pts[i]
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}
	hull := make([]point, 0, len(lower)+len(upper)-2)
	hull = append(hull, lower...)
	if len(upper) > 1 {
		hull = append(hull, upper[1:len(upper)-1]...)
	}
	return hull
}

func expectedAnswer(n int, pts []point) float64 {
	hull := convexHull(append([]point(nil), pts...))
	c := len(hull)
	baseArea := areaPoly(hull)
	vis := make([]bool, n)
	for _, hp := range hull {
		for i, p := range pts {
			if math.Abs(p.x-hp.x) < 1e-9 && math.Abs(p.y-hp.y) < 1e-9 {
				vis[i] = true
			}
		}
	}
	ans := -1.0
	if c < 4 {
		for i, p := range pts {
			if vis[i] {
				continue
			}
			minA := math.Inf(1)
			for j := 0; j < c; j++ {
				k := (j + 1) % c
				a := areaTri(p, hull[j], hull[k])
				if a < minA {
					minA = a
				}
			}
			if baseArea-minA > ans {
				ans = baseArea - minA
			}
		}
	} else {
		for i := 0; i < c-3; i++ {
			for j := i + 1; j < c-2; j++ {
				for k := j + 1; k < c-1; k++ {
					for l := k + 1; l < c; l++ {
						quad := []point{hull[i], hull[j], hull[k], hull[l]}
						a := areaPoly(quad)
						if a > ans {
							ans = a
						}
					}
				}
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (int, []point) {
	n := rng.Intn(10) + 4 // small for runtime
	pts := make([]point, 0, n)
	for len(pts) < n {
		x := float64(rng.Intn(2001) - 1000)
		y := float64(rng.Intn(2001) - 1000)
		ok := true
		for _, p := range pts {
			if p.x == x && p.y == y {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// check no three collinear
		collinear := false
		for i := 0; i < len(pts) && !collinear; i++ {
			for j := i + 1; j < len(pts) && !collinear; j++ {
				if math.Abs(cross(pts[i], pts[j], point{x, y})) < 1e-9 {
					collinear = true
				}
			}
		}
		if collinear {
			continue
		}
		pts = append(pts, point{x, y})
	}
	return n, pts
}

func runCase(bin string, n int, pts []point) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%.0f %.0f\n", p.x, p.y))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%.9f", expectedAnswer(n, pts))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, pts := generateCase(rng)
		if err := runCase(bin, n, pts); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
