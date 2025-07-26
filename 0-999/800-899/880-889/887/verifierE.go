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
)

func solveCase(x1, y1, x2, y2 float64, circles [][3]float64) float64 {
	mx := (x1 + x2) / 2
	my := (y1 + y2) / 2
	dx := x1 - x2
	dy := y1 - y2
	w := math.Hypot(dx, dy) / 2
	vx := -dy
	vy := dx
	norm := math.Hypot(vx, vy)
	vx /= norm
	vy /= norm
	get := func(px, py, val float64) float64 {
		d := func(t float64) float64 {
			cx := mx + vx*t
			cy := my + vy*t
			dist := math.Hypot(px-cx, py-cy)
			return dist - math.Sqrt(w*w+t*t)
		}
		l, r := -1e12, 1e12
		fl := d(l) > val
		for i := 0; i < 100; i++ {
			mid := (l + r) / 2
			f := d(mid) > val
			if f == fl {
				l = mid
			} else {
				r = mid
			}
		}
		return r
	}
	type pair struct {
		x float64
		y int
	}
	q := make([]pair, 0, 2*len(circles)+1)
	q = append(q, pair{0, 0})
	for _, c := range circles {
		l := get(c[0], c[1], c[2])
		rr := get(c[0], c[1], -c[2])
		if l > rr {
			l, rr = rr, l
		}
		q = append(q, pair{l, 1})
		q = append(q, pair{rr, -1})
	}
	sort.Slice(q, func(i, j int) bool {
		if q[i].x == q[j].x {
			return q[i].y < q[j].y
		}
		return q[i].x < q[j].x
	})
	s := 0
	mn := 1e12
	for _, p := range q {
		if s == 0 {
			if math.Abs(p.x) < mn {
				mn = math.Abs(p.x)
			}
		}
		s += p.y
		if s == 0 {
			if math.Abs(p.x) < mn {
				mn = math.Abs(p.x)
			}
		}
	}
	return math.Sqrt(mn*mn + w*w)
}

func runCase(bin string, x1, y1, x2, y2 float64, circles [][3]float64) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%f %f %f %f\n", x1, y1, x2, y2)
	fmt.Fprintf(&sb, "%d\n", len(circles))
	for _, c := range circles {
		fmt.Fprintf(&sb, "%f %f %f\n", c[0], c[1], c[2])
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
	fmt.Fscan(strings.NewReader(out.String()), &got)
	exp := solveCase(x1, y1, x2, y2, circles)
	if math.Abs(got-exp) > 1e-4*math.Max(1, math.Abs(exp)) {
		return fmt.Errorf("expected %.5f got %.5f", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	for i := 0; i < 100; i++ {
		x1 := rand.Float64()*200 - 100
		y1 := rand.Float64()*200 - 100
		x2 := rand.Float64()*200 - 100
		y2 := rand.Float64()*200 - 100
		n := rand.Intn(5) + 1
		circles := make([][3]float64, n)
		for j := 0; j < n; j++ {
			circles[j][0] = rand.Float64()*200 - 100
			circles[j][1] = rand.Float64()*200 - 100
			circles[j][2] = rand.Float64()*10 + 1
		}
		if err := runCase(bin, x1, y1, x2, y2, circles); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
