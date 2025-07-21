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

type Point struct{ x, y float64 }

func (p Point) add(q Point) Point   { return Point{p.x + q.x, p.y + q.y} }
func (p Point) mul(s float64) Point { return Point{p.x * s, p.y * s} }

func dist(p, q Point) float64 { return math.Hypot(p.x-q.x, p.y-q.y) }

func solve(input string) string {
	rdr := strings.NewReader(input)
	var t1, t2 float64
	fmt.Fscan(rdr, &t1, &t2)
	var a, b, c Point
	fmt.Fscan(rdr, &a.x, &a.y, &c.x, &c.y, &b.x, &b.y)
	ab := dist(a, b)
	bc := dist(b, c)
	ac := dist(a, c)
	t1 += ab + bc + 1e-12
	t2 += ac + 1e-12
	if ab+bc < t2 {
		return fmt.Sprintf("%.10f", math.Min(t1, t2))
	}
	cal := func(k float64) float64 {
		p := b.mul(1 - k).add(c.mul(k))
		ap := dist(a, p)
		if ap+(k+1)*bc < t1 && ap+(1-k)*bc < t2 {
			if t1-(k+1)*bc < t2-(1-k)*bc {
				return t1 - (k+1)*bc
			}
			return t2 - (1-k)*bc
		}
		l, r := 0.0, 1.0
		for r-l > 1e-15 {
			m := (l + r) / 2
			p1 := a.mul(1 - m).add(p.mul(m))
			if ap*m+dist(p1, b)+bc < t1 && ap*m+dist(p1, c) < t2 {
				l = m
			} else {
				r = m
			}
		}
		return ((l + r) / 2) * ap
	}
	l, r := 0.0, 1.0
	for r-l > 1e-15 {
		m1 := (2*l + r) / 3
		m2 := (l + 2*r) / 3
		if cal(m1)-cal(m2) < 1e-12 {
			l = m1
		} else {
			r = m2
		}
	}
	ans := cal((l + r) / 2)
	return fmt.Sprintf("%.10f", ans)
}

type test struct{ input, expected string }

func generateTests() []test {
	rand.Seed(123)
	var tests []test
	for len(tests) < 100 {
		t1 := rand.Float64() * 50
		t2 := rand.Float64() * 50
		ax := rand.Float64()*20 - 10
		ay := rand.Float64()*20 - 10
		bx := rand.Float64()*20 - 10
		by := rand.Float64()*20 - 10
		cx := rand.Float64()*20 - 10
		cy := rand.Float64()*20 - 10
		inp := fmt.Sprintf("%.2f %.2f\n%.2f %.2f %.2f %.2f %.2f %.2f\n", t1, t2, ax, ay, cx, cy, bx, by)
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.expected {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
