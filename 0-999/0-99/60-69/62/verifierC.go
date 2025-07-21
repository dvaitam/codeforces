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

type Point struct{ x, y float64 }

func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }

func segmentIntersect(p, r, q, s Point) (float64, bool) {
	denom := cross(r, s)
	if math.Abs(denom) < 1e-9 {
		return 0, false
	}
	qp := Point{q.x - p.x, q.y - p.y}
	t := cross(qp, s) / denom
	u := cross(qp, r) / denom
	if t > -1e-9 && t < 1+1e-9 && u > -1e-9 && u < 1+1e-9 {
		return t, true
	}
	return 0, false
}

func pointInTriangle(p, a, b, c Point) bool {
	ab := Point{b.x - a.x, b.y - a.y}
	bc := Point{c.x - b.x, c.y - b.y}
	ca := Point{a.x - c.x, a.y - c.y}
	ap := Point{p.x - a.x, p.y - a.y}
	bp := Point{p.x - b.x, p.y - b.y}
	cp := Point{p.x - c.x, p.y - c.y}
	c1 := cross(ab, ap)
	c2 := cross(bc, bp)
	c3 := cross(ca, cp)
	if (c1 >= -1e-9 && c2 >= -1e-9 && c3 >= -1e-9) || (c1 <= 1e-9 && c2 <= 1e-9 && c3 <= 1e-9) {
		return true
	}
	return false
}

func solveCase(tris [][3]Point) float64 {
	n := len(tris)
	total := 0.0
	for i := 0; i < n; i++ {
		for e := 0; e < 3; e++ {
			a := tris[i][e]
			b := tris[i][(e+1)%3]
			r := Point{b.x - a.x, b.y - a.y}
			ts := []float64{0.0, 1.0}
			for j := 0; j < n; j++ {
				if j == i {
					continue
				}
				for f := 0; f < 3; f++ {
					c := tris[j][f]
					d := tris[j][(f+1)%3]
					s := Point{d.x - c.x, d.y - c.y}
					if t, ok := segmentIntersect(a, r, c, s); ok {
						if t > 1e-9 && t < 1-1e-9 {
							ts = append(ts, t)
						}
					}
				}
			}
			sort.Float64s(ts)
			unique := ts[:0]
			for _, t := range ts {
				if len(unique) == 0 || math.Abs(t-unique[len(unique)-1]) > 1e-8 {
					unique = append(unique, t)
				}
			}
			for k := 0; k+1 < len(unique); k++ {
				t0 := unique[k]
				t1 := unique[k+1]
				tm := (t0 + t1) * 0.5
				pm := Point{a.x + r.x*tm, a.y + r.y*tm}
				covered := false
				for j := 0; j < n && !covered; j++ {
					if j == i {
						continue
					}
					if pointInTriangle(pm, tris[j][0], tris[j][1], tris[j][2]) {
						covered = true
					}
				}
				if !covered {
					total += (t1 - t0) * math.Hypot(r.x, r.y)
				}
			}
		}
	}
	return total
}

func runCase(bin string, tris [][3]Point) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tris)))
	for _, tr := range tris {
		sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f\n", tr[0].x, tr[0].y, tr[1].x, tr[1].y, tr[2].x, tr[2].y))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(3) + 1
		tris := make([][3]Point, n)
		for i := 0; i < n; i++ {
			for {
				a := Point{rng.Float64() * 10, rng.Float64() * 10}
				b := Point{rng.Float64() * 10, rng.Float64() * 10}
				c := Point{rng.Float64() * 10, rng.Float64() * 10}
				if math.Abs(cross(Point{b.x - a.x, b.y - a.y}, Point{c.x - a.x, c.y - a.y})) > 1e-3 {
					tris[i] = [3]Point{a, b, c}
					break
				}
			}
		}
		expected := solveCase(tris)
		gotStr, err := runCase(bin, tris)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(gotStr), 64)
		if err != nil || math.Abs(got-expected) > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %s\n", t+1, expected, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
