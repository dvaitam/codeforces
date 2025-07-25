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

type Pt struct{ x, y float64 }

func dist(a, b Pt) float64 { dx := a.x - b.x; dy := a.y - b.y; return math.Hypot(dx, dy) }

func cross(a, b, c Pt) float64 { return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x) }

func convexHull(pts []Pt) []Pt {
	if len(pts) == 0 {
		return nil
	}
	p := make([]Pt, len(pts))
	copy(p, pts)
	sort.Slice(p, func(i, j int) bool {
		if p[i].x == p[j].x {
			return p[i].y < p[j].y
		}
		return p[i].x < p[j].x
	})
	n := len(p)
	hull := make([]Pt, 0, 2*n)
	for i := 0; i < n; i++ {
		for len(hull) >= 2 && cross(hull[len(hull)-2], hull[len(hull)-1], p[i]) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p[i])
	}
	l := len(hull)
	for i := n - 2; i >= 0; i-- {
		for len(hull) > l && cross(hull[len(hull)-2], hull[len(hull)-1], p[i]) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p[i])
	}
	if len(hull) > 1 {
		hull = hull[:len(hull)-1]
	}
	return hull
}

func pointInPoly(poly []Pt, p Pt) bool {
	n := len(poly)
	if n == 0 {
		return false
	}
	pos, neg := false, false
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		val := cross(a, b, p)
		if val > 0 {
			pos = true
		} else if val < 0 {
			neg = true
		}
		if pos && neg {
			return false
		}
	}
	return true
}

func solveD(reds, blues []Pt) float64 {
	if len(blues) < 3 {
		return math.Inf(1)
	}
	hull := convexHull(blues)
	for _, r := range reds {
		if !pointInPoly(hull, r) {
			return math.Inf(1)
		}
	}
	best := 0.0
	for _, r := range reds {
		mind := math.Inf(1)
		for _, b := range blues {
			d := dist(r, b)
			if d < mind {
				mind = d
			}
		}
		if mind > best {
			best = mind
		}
	}
	return best
}

func genCase(rng *rand.Rand) (string, float64) {
	blues := []Pt{{0, 0}, {10, 0}, {0, 10}, {10, 10}}
	n := rng.Intn(3) + 1
	reds := make([]Pt, n)
	for i := 0; i < n; i++ {
		reds[i] = Pt{float64(rng.Intn(9) + 1), float64(rng.Intn(9) + 1)}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(blues))
	for _, r := range reds {
		fmt.Fprintf(&sb, "%.0f %.0f\n", r.x, r.y)
	}
	for _, b := range blues {
		fmt.Fprintf(&sb, "%.0f %.0f\n", b.x, b.y)
	}
	ans := solveD(reds, blues)
	if math.IsInf(ans, 1) {
		return sb.String(), math.Inf(1)
	}
	return sb.String(), ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func equalFloat(a, b float64) bool { return math.Abs(a-b) <= 1e-4*math.Max(1, math.Abs(b)) }

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, expect := genCase(rng)
		outStr, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		if math.IsInf(expect, 1) {
			if strings.TrimSpace(outStr) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:\n%s", t, outStr, in)
				os.Exit(1)
			}
		} else {
			var got float64
			if _, err := fmt.Sscan(outStr, &got); err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: bad output %q\n", t, outStr)
				os.Exit(1)
			}
			if !equalFloat(got, expect) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %.5f got %.5f\ninput:\n%s", t, expect, got, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
