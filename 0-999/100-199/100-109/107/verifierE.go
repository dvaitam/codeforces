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

const EPS = 1e-8

type PT struct{ x, y float64 }

func (p PT) Sub(q PT) PT        { return PT{p.x - q.x, p.y - q.y} }
func (p PT) Cross(q PT) float64 { return p.x*q.y - p.y*q.x }

func SG(x float64) int {
	if x > EPS {
		return 1
	}
	if x < -EPS {
		return -1
	}
	return 0
}

func tri(p1, p2, p3 PT) float64 { return p2.Sub(p1).Cross(p3.Sub(p1)) }

func segP(p, p1, p2 PT) float64 {
	if SG(p1.x-p2.x) == 0 {
		return (p.y - p1.y) / (p2.y - p1.y)
	}
	return (p.x - p1.x) / (p2.x - p1.x)
}

func polyUnion(polys [][]PT) float64 {
	n := len(polys)
	var sum float64
	for i := 0; i < n; i++ {
		p := polys[i]
		m := len(p)
		p2 := make([]PT, m+1)
		copy(p2, p)
		p2[m] = p[0]
		for ii := 0; ii < m; ii++ {
			a, b := p2[ii], p2[ii+1]
			type inter struct {
				t float64
				d int
			}
			ev := []inter{{0, 0}, {1, 0}}
			for j := 0; j < n; j++ {
				if j == i {
					continue
				}
				q := polys[j]
				lq := len(q)
				q2 := make([]PT, lq+1)
				copy(q2, q)
				q2[lq] = q[0]
				for jj := 0; jj < lq; jj++ {
					c, dpt := q2[jj], q2[jj+1]
					ta := SG(tri(a, b, c))
					tb := SG(tri(a, b, dpt))
					if ta == 0 && tb == 0 {
						if q2[jj+1].Sub(q2[jj]).Cross(b.Sub(a)) > 0 && j < i {
							t1 := segP(c, a, b)
							t2 := segP(dpt, a, b)
							ev = append(ev, inter{t1, 1}, inter{t2, -1})
						}
					} else if ta >= 0 && tb < 0 {
						tc := tri(c, dpt, a)
						td := tri(c, dpt, b)
						ev = append(ev, inter{tc / (tc - td), 1})
					} else if ta < 0 && tb >= 0 {
						tc := tri(c, dpt, a)
						td := tri(c, dpt, b)
						ev = append(ev, inter{tc / (tc - td), -1})
					}
				}
			}
			sort.Slice(ev, func(i, j int) bool { return ev[i].t < ev[j].t })
			z := math.Min(math.Max(ev[0].t, 0), 1)
			dcnt := ev[0].d
			var covered float64
			for k := 1; k < len(ev); k++ {
				w := math.Min(math.Max(ev[k].t, 0), 1)
				if dcnt == 0 {
					covered += w - z
				}
				dcnt += ev[k].d
				z = w
			}
			sum += a.Cross(b) * covered
		}
	}
	return sum / 2
}

func solveE(rects [][]PT) float64 {
	var sumArea float64
	for i := 0; i < len(rects); i++ {
		pts := rects[i]
		var s float64
		for j := 0; j < 3; j++ {
			s += pts[j].Cross(pts[j+1])
		}
		s += pts[3].Cross(pts[0])
		area := s / 2
		if area < 0 {
			for l, r := 0, 3; l < r; l, r = l+1, r-1 {
				pts[l], pts[r] = pts[r], pts[l]
			}
			area = -area
		}
		sumArea += area
		rects[i] = pts
	}
	union := polyUnion(rects)
	return sumArea / union
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

func randomRect(rng *rand.Rand) []PT {
	// axis-aligned for simplicity
	x1 := rng.Float64()*20 - 10
	y1 := rng.Float64()*20 - 10
	w := rng.Float64()*5 + 0.1
	h := rng.Float64()*5 + 0.1
	rect := []PT{{x1, y1}, {x1 + w, y1}, {x1 + w, y1 + h}, {x1, y1 + h}}
	// random rotate
	ang := rng.Float64() * math.Pi * 2
	cx := x1 + w/2
	cy := y1 + h/2
	for i := 0; i < 4; i++ {
		rx := rect[i].x - cx
		ry := rect[i].y - cy
		rect[i].x = cx + rx*math.Cos(ang) - ry*math.Sin(ang)
		rect[i].y = cy + rx*math.Sin(ang) + ry*math.Cos(ang)
	}
	return rect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		rects := make([][]PT, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			rects[j] = randomRect(rng)
			for k := 0; k < 4; k++ {
				sb.WriteString(fmt.Sprintf("%f %f", rects[j][k].x, rects[j][k].y))
				if k == 3 {
					sb.WriteByte('\n')
				} else {
					sb.WriteByte(' ')
				}
			}
		}
		input := sb.String()
		expect := solveE(rects)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		val, err2 := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err2 != nil {
			fmt.Printf("case %d: cannot parse output\n", i+1)
			os.Exit(1)
		}
		if math.Abs(val-expect) > 1e-6*math.Max(1, math.Abs(expect)) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:%.6f\ngot:%s\n", i+1, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
