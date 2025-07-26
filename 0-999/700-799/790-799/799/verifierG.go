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

type Point struct{ x, y float64 }

func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }
func sub(a, b Point) Point     { return Point{a.x - b.x, a.y - b.y} }

func polygonArea(poly []Point) float64 {
	area := 0.0
	n := len(poly)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += cross(poly[i], poly[j])
	}
	if area < 0 {
		area = -area
	}
	return area / 2
}

func cutPolygon(poly []Point, p, dir Point) []Point {
	res := make([]Point, 0, len(poly)+2)
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		da := cross(dir, sub(a, p))
		db := cross(dir, sub(b, p))
		if da >= 0 {
			res = append(res, a)
		}
		if da*db < 0 {
			t := da / (da - db)
			inter := Point{a.x + t*(b.x-a.x), a.y + t*(b.y-a.y)}
			res = append(res, inter)
		}
	}
	return res
}

func findAngle(poly []Point, p Point, total float64) float64 {
	lo, hi := 0.0, math.Pi
	for iter := 0; iter < 60; iter++ {
		mid := (lo + hi) / 2
		dir := Point{math.Cos(mid), math.Sin(mid)}
		half := cutPolygon(poly, p, dir)
		area := polygonArea(half)
		if area > total/2 {
			hi = mid
		} else {
			lo = mid
		}
	}
	return (lo + hi) / 2
}

func genConvex(n int) []Point {
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = rand.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	poly := make([]Point, n)
	for i := 0; i < n; i++ {
		r := rand.Float64()*5 + 1
		poly[i] = Point{r * math.Cos(angles[i]), r * math.Sin(angles[i])}
	}
	return poly
}

func genG() (string, []float64) {
	n := rand.Intn(3) + 3
	q := rand.Intn(3) + 1
	poly := genConvex(n)
	cx, cy := 0.0, 0.0
	for _, p := range poly {
		cx += p.x
		cy += p.y
	}
	cx /= float64(n)
	cy /= float64(n)
	queries := make([]Point, q)
	for i := 0; i < q; i++ {
		queries[i] = Point{cx + (rand.Float64()-0.5)*0.1, cy + (rand.Float64()-0.5)*0.1}
	}
	total := polygonArea(poly)
	answers := make([]float64, q)
	for i := 0; i < q; i++ {
		answers[i] = findAngle(poly, queries[i], total)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%.6f %.6f\n", poly[i].x, poly[i].y))
	}
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%.6f %.6f\n", queries[i].x, queries[i].y))
	}
	return sb.String(), answers
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genG()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s", i+1, err, in, got)
			return
		}
		lines := strings.Fields(got)
		if len(lines) != len(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected %d numbers, got %d\n%s", i+1, in, len(exp), len(lines), got)
			return
		}
		for j, field := range lines {
			var val float64
			fmt.Sscan(field, &val)
			if math.Abs(val-exp[j]) > 1e-4 {
				fmt.Printf("Test %d failed at line %d\nInput:\n%sExpected: %.6f Got: %s", i+1, j+1, in, exp[j], field)
				return
			}
		}
	}
	fmt.Println("All tests passed")
}
