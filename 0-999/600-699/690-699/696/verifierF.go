package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		poly := make([]com, 2*n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			scan.Scan()
			x, _ := strconv.ParseFloat(scan.Text(), 64)
			scan.Scan()
			y, _ := strconv.ParseFloat(scan.Text(), 64)
			poly[i] = com{x, y}
			sb.WriteString(fmt.Sprintf("%g %g\n", x, y))
		}
		for i := n; i < 2*n; i++ {
			poly[i] = poly[i-n]
		}
		input := sb.String()
		r, b, l := solveF(poly, n)
		expected := formatAns(r, b, l)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
