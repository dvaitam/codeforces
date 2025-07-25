package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const (
	magic = 49
	D     = 105
)

type Point struct{ x, y int }

type testCase struct{ pts []Point }

var half [D]float64
var two [D]float64

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCase(p []Point) float64 {
	n := len(p)
	half[0], two[0] = 1.0, 1.0
	for i := 1; i < D; i++ {
		half[i] = half[i-1] / 2.0
		two[i] = two[i-1] * 2.0
	}
	area, border := 0.0, 0.0
	var z float64
	if n <= 50 {
		z = two[n] - 1.0 - float64(n) - float64(n*(n-1)/2)
	}
	for i := 0; i < n; i++ {
		for j0 := i + 1; j0 <= min(i+n-2, i+magic); j0++ {
			var now float64
			if n > 50 {
				now = half[j0-i+1]
			} else {
				now = two[n-(j0-i+1)] - 1.0
				now = now / z
			}
			jj := j0 % n
			dx := p[i].x - p[jj].x
			if dx < 0 {
				dx = -dx
			}
			dy := p[i].y - p[jj].y
			if dy < 0 {
				dy = -dy
			}
			line := gcd(dx, dy)
			border += float64(line) * now
			cross := float64(p[i].x)*float64(p[jj].y) - float64(p[i].y)*float64(p[jj].x)
			area += 0.5 * now * cross
		}
	}
	return area - border/2.0 + 1.0
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.pts))
	for _, pt := range tc.pts {
		fmt.Fprintf(&sb, "%d %d\n", pt.x, pt.y)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	exp := solveCase(tc.pts)
	if math.Abs(got-exp) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", exp, got)
	}
	return nil
}

func genPolygon(n, r, ox, oy int) []Point {
	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		angle := 2 * math.Pi * float64(i) / float64(n)
		x := ox + int(math.Round(float64(r)*math.Cos(angle)))
		y := oy + int(math.Round(float64(r)*math.Sin(angle)))
		pts[i] = Point{x, y}
	}
	return pts
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	cases := make([]testCase, 100)
	for i := range cases {
		n := rand.Intn(5) + 3
		r := rand.Intn(20) + 5
		ox := rand.Intn(20) - 10
		oy := rand.Intn(20) - 10
		cases[i] = testCase{genPolygon(n, r, ox, oy)}
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
