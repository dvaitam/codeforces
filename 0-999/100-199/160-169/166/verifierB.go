package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Point struct{ x, y int64 }

func cross(a, b, c Point) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func pointInConvex(poly []Point, q Point) bool {
	n := len(poly)
	if cross(poly[0], poly[1], q) <= 0 {
		return false
	}
	if cross(poly[0], poly[n-1], q) >= 0 {
		return false
	}
	l, r := 1, n-1
	for l+1 < r {
		m := (l + r) >> 1
		if cross(poly[0], poly[m], q) > 0 {
			l = m
		} else {
			r = m
		}
	}
	if cross(poly[l], poly[l+1], q) <= 0 {
		return false
	}
	return true
}

func expected(polyA, polyB []Point) string {
	poly := make([]Point, len(polyA))
	copy(poly, polyA)
	for i := 0; i < len(poly)/2; i++ {
		poly[i], poly[len(poly)-1-i] = poly[len(poly)-1-i], poly[i]
	}
	for _, q := range polyB {
		if !pointInConvex(poly, q) {
			return "NO"
		}
	}
	return "YES"
}

func randomPointInside(poly []Point, rng *rand.Rand) Point {
	for {
		var px, py float64
		remain := 1.0
		for i := 0; i < len(poly)-1; i++ {
			w := rng.Float64() * remain
			px += w * float64(poly[i].x)
			py += w * float64(poly[i].y)
			remain -= w
		}
		px += remain * float64(poly[len(poly)-1].x)
		py += remain * float64(poly[len(poly)-1].y)
		p := Point{int64(math.Round(px)), int64(math.Round(py))}
		if pointInConvex(poly, p) {
			return p
		}
	}
}

func genPolyA(rng *rand.Rand) []Point {
	n := rng.Intn(6) + 3
	base := rng.Float64() * 2 * math.Pi
	r := float64(100 + rng.Intn(50))
	poly := make([]Point, n)
	for i := 0; i < n; i++ {
		ang := base - float64(i)*2*math.Pi/float64(n)
		x := int64(math.Round(r*math.Cos(ang))) + 1000
		y := int64(math.Round(r*math.Sin(ang))) + 1000
		poly[i] = Point{x, y}
	}
	return poly
}

func genPolyBInside(A []Point, rng *rand.Rand) []Point {
	m := rng.Intn(6) + 3
	center := randomPointInside(A, rng)
	rad := float64(rng.Intn(30)+1) / 2
	base := rng.Float64() * 2 * math.Pi
	poly := make([]Point, m)
	for i := 0; i < m; i++ {
		ang := base - float64(i)*2*math.Pi/float64(m)
		x := center.x + int64(math.Round(rad*math.Cos(ang)))
		y := center.y + int64(math.Round(rad*math.Sin(ang)))
		poly[i] = Point{x, y}
	}
	return poly
}

func genPolyBOutside(A []Point, rng *rand.Rand) []Point {
	m := rng.Intn(6) + 3
	cx := A[0].x + 300
	cy := A[0].y + 300
	rad := float64(rng.Intn(30)+1) / 2
	base := rng.Float64() * 2 * math.Pi
	poly := make([]Point, m)
	for i := 0; i < m; i++ {
		ang := base - float64(i)*2*math.Pi/float64(m)
		x := cx + int64(math.Round(rad*math.Cos(ang)))
		y := cy + int64(math.Round(rad*math.Sin(ang)))
		poly[i] = Point{x, y}
	}
	return poly
}

func generateCase(rng *rand.Rand) (string, string) {
	A := genPolyA(rng)
	inside := rng.Intn(2) == 0
	var B []Point
	if inside {
		B = genPolyBInside(A, rng)
	} else {
		B = genPolyBOutside(A, rng)
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, len(A))
	for _, p := range A {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	fmt.Fprintln(&sb, len(B))
	for _, p := range B {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	return sb.String(), expected(A, B)
}

func runCase(exe, input, expectedAns string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	res := strings.TrimSpace(out.String())
	if res != expectedAns {
		return fmt.Errorf("expected %s got %s", expectedAns, res)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
