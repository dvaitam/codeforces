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

type pt struct{ x, y float64 }

func sub(a, b pt) pt        { return pt{a.x - b.x, a.y - b.y} }
func cross(a, b pt) float64 { return a.x*b.y - a.y*b.x }
func dot(a, b pt) float64   { return a.x*b.x + a.y*b.y }

func work(ext []pt, n int, ans *float64) {
	sum := 0.0
	j := 1
	for i := 0; i < n; i++ {
		for j+1 < len(ext) && dot(sub(ext[i+1], ext[i]), sub(ext[j], ext[j+1])) < 0 {
			sum += cross(ext[j+1], ext[j])
			j++
		}
		u := sub(ext[i+1], ext[i])
		l1 := math.Hypot(u.x, u.y)
		w := sub(ext[j], ext[i])
		l2 := math.Hypot(w.x, w.y)
		h := cross(w, u) / l1
		l := math.Sqrt(l2*l2 - h*h)
		t := l / l1
		res := pt{ext[i].x + u.x*t, ext[i].y + u.y*t}
		area := cross(ext[j], res) + cross(res, ext[i+1]) - sum
		if area < *ans {
			*ans = area
		}
		if i+2 < len(ext) {
			sum -= cross(ext[i+2], ext[i+1])
		}
	}
}

func solve(pts []pt) float64 {
	n := len(pts)
	ext := make([]pt, 0, 2*n)
	ext = append(ext, pts...)
	ext = append(ext, pts...)
	if n >= 3 {
		u := sub(ext[1], ext[0])
		v := sub(ext[2], ext[1])
		if cross(u, v) > 0 {
			for i, j := 0, len(ext)-1; i < j; i, j = i+1, j-1 {
				ext[i], ext[j] = ext[j], ext[i]
			}
		}
	}
	ans := math.Inf(1)
	work(ext, n, &ans)
	for i := range ext {
		ext[i].x = -ext[i].x
	}
	for i, j := 0, len(ext)-1; i < j; i, j = i+1, j-1 {
		ext[i], ext[j] = ext[j], ext[i]
	}
	work(ext, n, &ans)
	return ans / 2
}

func convexHull(pts []pt) []pt {
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x == pts[j].x {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	n := len(pts)
	if n <= 1 {
		return pts
	}
	hull := make([]pt, 0, 2*n)
	for _, p := range pts {
		for len(hull) >= 2 && cross(sub(hull[len(hull)-1], hull[len(hull)-2]), sub(p, hull[len(hull)-1])) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}
	t := len(hull)
	for i := n - 2; i >= 0; i-- {
		p := pts[i]
		for len(hull) > t && cross(sub(hull[len(hull)-1], hull[len(hull)-2]), sub(p, hull[len(hull)-1])) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}
	if len(hull) > 1 {
		hull = hull[:len(hull)-1]
	}
	return hull
}

func runCase(exe string, input string, expected float64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6*math.Max(1, math.Abs(expected)) {
		return fmt.Errorf("expected %f got %f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m := rng.Intn(5) + 3
		raw := make([]pt, m)
		for j := 0; j < m; j++ {
			raw[j] = pt{rng.Float64()*20 - 10, rng.Float64()*20 - 10}
		}
		hull := convexHull(raw)
		n := len(hull)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, p := range hull {
			sb.WriteString(fmt.Sprintf("%f %f\n", p.x, p.y))
		}
		input := sb.String()
		exp := solve(hull)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
