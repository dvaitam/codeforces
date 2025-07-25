package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const eps = 1e-9

type Point struct{ x, y float64 }

func sub(a, b Point) Point         { return Point{a.x - b.x, a.y - b.y} }
func add(a, b Point) Point         { return Point{a.x + b.x, a.y + b.y} }
func mul(a Point, t float64) Point { return Point{a.x * t, a.y * t} }
func dot(a, b Point) float64       { return a.x*b.x + a.y*b.y }
func cross(a, b Point) float64     { return a.x*b.y - a.y*b.x }

func onSegment(p, a, b Point) bool {
	if math.Abs(cross(sub(b, a), sub(p, a))) > eps {
		return false
	}
	return dot(sub(p, a), sub(p, b)) <= eps
}

func pointInPoly(p Point, poly []Point) bool {
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		if onSegment(p, a, b) {
			return true
		}
	}
	wn := 0
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		if a.y <= p.y {
			if b.y > p.y && cross(sub(b, a), sub(p, a)) > eps {
				wn++
			}
		} else {
			if b.y <= p.y && cross(sub(b, a), sub(p, a)) < -eps {
				wn--
			}
		}
	}
	return wn != 0
}

func intersectionLength(poly []Point, p0, p1 Point) float64 {
	d := sub(p1, p0)
	l := math.Hypot(d.x, d.y)
	dir := Point{d.x / l, d.y / l}
	var ts []float64
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		da := sub(b, a)
		denom := cross(d, da)
		if math.Abs(denom) < eps {
			if math.Abs(cross(d, sub(a, p0))) < eps {
				t1 := dot(sub(a, p0), dir)
				t2 := dot(sub(b, p0), dir)
				if t2 < t1 {
					t1, t2 = t2, t1
				}
				ts = append(ts, t1, t2)
			}
			continue
		}
		t := cross(sub(a, p0), da) / denom
		s := cross(sub(a, p0), d) / denom
		if s >= -eps && s <= 1+eps {
			ts = append(ts, t)
		}
	}
	if len(ts) == 0 {
		return 0
	}
	sort.Float64s(ts)
	uniq := []float64{ts[0]}
	for i := 1; i < len(ts); i++ {
		if math.Abs(ts[i]-ts[i-1]) > eps {
			uniq = append(uniq, ts[i])
		}
	}
	res := 0.0
	for i := 0; i+1 < len(uniq); i++ {
		tmid := (uniq[i] + uniq[i+1]) / 2
		pt := add(p0, mul(dir, tmid))
		if pointInPoly(pt, poly) {
			res += uniq[i+1] - uniq[i]
		}
	}
	return res * l
}

func expectedF(poly []Point, lines [][2]Point) []float64 {
	ans := make([]float64, len(lines))
	for i, ln := range lines {
		ans[i] = intersectionLength(poly, ln[0], ln[1])
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 3
	m := rng.Intn(5) + 1
	poly := make([]Point, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		x := rng.Float64()*20 - 10
		y := rng.Float64()*20 - 10
		poly[i] = Point{x, y}
		sb.WriteString(fmt.Sprintf("%.2f %.2f\n", x, y))
	}
	lines := make([][2]Point, m)
	for i := 0; i < m; i++ {
		x1 := rng.Float64()*20 - 10
		y1 := rng.Float64()*20 - 10
		x2 := rng.Float64()*20 - 10
		y2 := rng.Float64()*20 - 10
		lines[i] = [2]Point{{x1, y1}, {x2, y2}}
		sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f\n", x1, y1, x2, y2))
	}
	ans := expectedF(poly, lines)
	var out strings.Builder
	for i, v := range ans {
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(fmt.Sprintf("%.10f", v))
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Fields(strings.TrimSpace(out.String()))
	expLines := strings.Fields(strings.TrimSpace(expected))
	if len(gotLines) != len(expLines) {
		return fmt.Errorf("wrong number of lines")
	}
	for i := range gotLines {
		g, err1 := strconv.ParseFloat(gotLines[i], 64)
		e, err2 := strconv.ParseFloat(expLines[i], 64)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("parse error")
		}
		if math.Abs(g-e) > 1e-6*math.Max(1, math.Abs(e)) {
			return fmt.Errorf("expected %.10f got %.10f", e, g)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
