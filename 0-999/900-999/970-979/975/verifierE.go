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

// PT represents a point or vector in 2D.
type PT struct{ x, y float64 }

func (p PT) Add(q PT) PT      { return PT{p.x + q.x, p.y + q.y} }
func (p PT) Sub(q PT) PT      { return PT{p.x - q.x, p.y - q.y} }
func (p PT) Mul(t float64) PT { return PT{p.x * t, p.y * t} }
func (p PT) Div(t float64) PT { return PT{p.x / t, p.y / t} }
func (p PT) Rot(t float64) PT {
	c, s := math.Cos(t), math.Sin(t)
	return PT{p.x*c - p.y*s, p.x*s + p.y*c}
}

func Dot(p, q PT) float64   { return p.x*q.x + p.y*q.y }
func Cross(p, q PT) float64 { return p.x*q.y - p.y*q.x }
func Angle(p, q PT) float64 { return math.Atan2(Cross(p, q), Dot(p, q)) }

// query types

type query struct {
	typ  int
	f, t int
	v    int
}

type testCase struct {
	pts []PT
	qs  []query
}

func solveCase(tc testCase) string {
	n := len(tc.pts)
	pts := make([]PT, n)
	copy(pts, tc.pts)
	// compute centroid
	var cen PT
	mass := 0.0
	for i := 2; i < n; i++ {
		p0 := pts[0]
		p1 := pts[i-1]
		p2 := pts[i]
		temp := PT{(p0.x + p1.x + p2.x) / 3.0, (p0.y + p1.y + p2.y) / 3.0}
		area2 := math.Abs(Cross(p1.Sub(p0), p2.Sub(p0)))
		cen = cen.Mul(mass).Add(temp.Mul(area2)).Div(mass + area2)
		mass += area2
	}
	for i := 0; i < n; i++ {
		pts[i] = pts[i].Sub(cen)
	}
	a, b := 0, 1
	ang := 0.0
	const twoPi = 2 * math.Pi
	up := PT{0, 1}
	var sb strings.Builder
	for _, qu := range tc.qs {
		if qu.typ == 1 {
			c1 := qu.f - 1
			if b == c1 {
				a, b = b, a
			}
			r := pts[b].Rot(ang)
			cen = cen.Add(r)
			tang := Angle(r, up)
			ang += tang
			ang = math.Mod(ang, twoPi)
			if ang < 0 {
				ang += twoPi
			}
			cen = cen.Sub(pts[b].Rot(ang))
			a = qu.t - 1
		} else {
			c := qu.v - 1
			r := pts[c].Rot(ang)
			p := r.Add(cen)
			sb.WriteString(fmt.Sprintf("%.8f %.8f\n", p.x, p.y))
		}
	}
	return sb.String()
}

func buildCase(tc testCase) (string, string) {
	var sb strings.Builder
	n := len(tc.pts)
	q := len(tc.qs)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, p := range tc.pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", int(p.x), int(p.y)))
	}
	for _, qu := range tc.qs {
		if qu.typ == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d\n", qu.f, qu.t))
		} else {
			sb.WriteString(fmt.Sprintf("2 %d\n", qu.v))
		}
	}
	return sb.String(), solveCase(tc)
}

func randomPolygon(rng *rand.Rand) []PT {
	n := rng.Intn(3) + 3 // 3..5
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	pts := make([]PT, n)
	for i, a := range angles {
		r := float64(rng.Intn(9) + 1)
		x := math.Round(r * math.Cos(a))
		y := math.Round(r * math.Sin(a))
		pts[i] = PT{x, y}
	}
	return pts
}

func genCase(rng *rand.Rand) (string, string) {
	tc := testCase{}
	tc.pts = randomPolygon(rng)
	q := rng.Intn(5) + 1
	tc.qs = make([]query, q)
	pinned := [2]int{1, 2}
	has2 := false
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			idx := rng.Intn(2)
			f := pinned[idx]
			t := rng.Intn(len(tc.pts)) + 1
			pinned[idx] = t
			tc.qs[i] = query{typ: 1, f: f, t: t}
		} else {
			v := rng.Intn(len(tc.pts)) + 1
			tc.qs[i] = query{typ: 2, v: v}
			has2 = true
		}
	}
	if !has2 {
		tc.qs[0] = query{typ: 2, v: 1}
	}
	return buildCase(tc)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotFields := strings.Fields(strings.TrimSpace(out.String()))
	expFields := strings.Fields(strings.TrimSpace(exp))
	if len(gotFields) != len(expFields) {
		return fmt.Errorf("expected %d numbers got %d", len(expFields), len(gotFields))
	}
	for i := 0; i < len(expFields); i++ {
		g, err := strconv.ParseFloat(gotFields[i], 64)
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		e, _ := strconv.ParseFloat(expFields[i], 64)
		if math.Abs(g-e) > 1e-4*math.Max(1, math.Abs(e)) {
			return fmt.Errorf("expected %v got %v", exp, out.String())
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
