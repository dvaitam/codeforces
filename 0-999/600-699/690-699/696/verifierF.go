package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"100",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"4",
	"10 0",
	"0 10",
	"-10 0",
	"0 -10",
	"3",
	"10 0",
	"-5 9",
	"-5 -9",
	"5",
	"10 0",
	"3 10",
	"-8 6",
	"-8 -6",
	"3 -10",
}

type R = float64

type com struct{ x, y R }

func (c com) add(d com) com { return com{c.x + d.x, c.y + d.y} }
func (c com) sub(d com) com { return com{c.x - d.x, c.y - d.y} }
func (c com) mulf(f R) com  { return com{c.x * f, c.y * f} }
func (c com) divf(f R) com  { return com{c.x / f, c.y / f} }
func (c com) abs() R        { return math.Hypot(c.x, c.y) }
func (c com) rot90() com    { return com{-c.y, c.x} }

func dot(a, b com) R { return a.x*b.x + a.y*b.y }
func det(a, b com) R { return a.x*b.y - a.y*b.x }

const eps = 1e-9

type line struct {
	n com
	c R
}

func newLinePoints(p1, p2 com) line {
	d := p2.sub(p1)
	nd := d.divf(d.abs())
	n := nd.rot90()
	return line{n: n, c: dot(p1, n)}
}

func (l line) dir() com { return l.n.rot90() }

func dist(l line, p com) R { return math.Abs(dot(p, l.n) - l.c) }

func check_r(radius R, polygon []com, n int, candB, candL *com) bool {
	i1, m := 2, 0
	best := make([]int, n)
	cand := make([]com, n)
	for i := 0; i < n; i++ {
		cand[i] = polygon[i+1]
	}
	for i := 0; i < n; i++ {
		if i1 < i+2 {
			i1 = i + 2
		}
		if m < i {
			m = i
		}
		best[i] = i1 - i
		l1 := newLinePoints(polygon[i], polygon[i+1])
		for i1 != i+n {
			l2 := newLinePoints(polygon[i1], polygon[i1+1])
			if det(polygon[i+1].sub(polygon[i]), polygon[i1+1].sub(polygon[i1])) < eps {
				break
			}
			for dist(l1, polygon[m+1]) < dist(l2, polygon[m+1]) {
				m++
			}
			low, high := polygon[m], polygon[m+1]
			for k := 0; k < 42; k++ {
				mid := low.add(high).mulf(0.5)
				if dist(l1, mid) < dist(l2, mid) {
					low = mid
				} else {
					high = mid
				}
			}
			if dist(l1, low) < radius {
				i1++
				best[i] = i1 - i
				cand[i] = low
			} else {
				break
			}
		}
	}
	for i := 0; i < n; i++ {
		next := (i + best[i]) % n
		if best[i]+best[next] >= n {
			*candB = cand[i]
			*candL = cand[next]
			return true
		}
	}
	return false
}

func solveF(poly []com, n int) (R, com, com) {
	start, fin := R(0), R(100000)
	var Bar, Lya com
	for it := 0; it < 84; it++ {
		m := (start + fin) / 2
		var candB, candL com
		if check_r(m, poly, n, &candB, &candL) {
			fin = m
			Bar = candB
			Lya = candL
		} else {
			start = m
		}
	}
	return fin, Bar, Lya
}

func formatAns(r R, b, l com) string {
	return fmt.Sprintf("%.10f\n%.10f %.10f\n%.10f %.10f\n", r, b.x, b.y, l.x, l.y)
}

type testCase struct {
	n     int
	points []com
	input string
}

func parseCases() []testCase {
	lines := rawTestcases
	pos := 0
	nextLine := func() string {
		s := strings.TrimSpace(lines[pos])
		pos++
		return s
	}
	t, _ := strconv.Atoi(nextLine())
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, _ := strconv.Atoi(nextLine())
		points := make([]com, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			parts := strings.Fields(nextLine())
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			points[j] = com{x, y}
			fmt.Fprintf(&sb, "%g %g\n", x, y)
		}
		cases = append(cases, testCase{n: n, points: points, input: sb.String()})
	}
	return cases
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases := parseCases()
	for caseIdx, tc := range cases {
		poly := make([]com, 2*tc.n)
		copy(poly, tc.points)
		copy(poly[tc.n:], tc.points)
		r, b, l := solveF(poly, tc.n)
		expected := formatAns(r, b, l)
		if err := runCase(exe, tc.input, expected); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
