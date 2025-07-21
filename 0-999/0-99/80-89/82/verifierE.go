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

type pd struct{ X, Y float64 }

type line struct{ A, B, C float64 }

const eps = 1e-7

func getLine(a, b pd) line {
	A := b.Y - a.Y
	B := a.X - b.X
	C := -A*a.X - B*a.Y
	return line{A, B, C}
}

func cp(a, b pd) float64 { return a.X*b.Y - a.Y*b.X }

func intersect(a, b line) (pd, bool) {
	d := cp(pd{a.A, a.B}, pd{b.A, b.B})
	if math.Abs(d) < eps {
		return pd{}, false
	}
	x := -cp(pd{a.C, a.B}, pd{b.C, b.B}) / d
	y := -cp(pd{a.A, a.C}, pd{b.A, b.C}) / d
	return pd{x, y}, true
}

func equals(a, b pd) bool { return math.Hypot(a.X-b.X, a.Y-b.Y) < eps }

func getArea(p []pd) float64 {
	if len(p) < 3 {
		return 0
	}
	pts := append([]pd(nil), p...)
	var c pd
	n := float64(len(pts))
	for _, v := range pts {
		c.X += v.X / n
		c.Y += v.Y / n
	}
	for i := range pts {
		pts[i].X -= c.X
		pts[i].Y -= c.Y
	}
	type ap struct {
		ang float64
		pt  pd
	}
	arr := make([]ap, len(pts))
	for i, v := range pts {
		arr[i] = ap{math.Atan2(v.Y, v.X), v}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].ang < arr[j].ang })
	pts = pts[:0]
	for _, v := range arr {
		pts = append(pts, v.pt)
	}
	uniq := make([]pd, 0, len(pts))
	for _, v := range pts {
		if len(uniq) == 0 || !equals(uniq[len(uniq)-1], v) {
			uniq = append(uniq, v)
		}
	}
	pts = uniq
	res := 0.0
	for i := range pts {
		j := (i + 1) % len(pts)
		ax := pts[i].X - c.X
		ay := pts[i].Y - c.Y
		bx := pts[j].X - c.X
		by := pts[j].Y - c.Y
		res += ax*by - ay*bx
	}
	return math.Abs(res) / 2
}

func solveE(n int, h, f float64, segs [][2]float64) float64 {
	s := make([]pd, 0, n*2)
	for i := 0; i < n; i++ {
		l, r := segs[i][0], segs[i][1]
		if l < 0 && r > 0 {
			s = append(s, pd{l, 0}, pd{0, r})
		} else {
			s = append(s, pd{l, r})
		}
	}
	cu := pd{0, f}
	cd := pd{0, -f}
	ul := getLine(pd{0, h}, pd{1, h})
	dl := getLine(pd{0, -h}, pd{1, -h})
	ans := 0.0
	for i := range s {
		lp := pd{s[i].X, -h}
		rp := pd{s[i].Y, -h}
		lLine := getLine(cd, lp)
		rLine := getLine(cd, rp)
		ulp, _ := intersect(lLine, ul)
		urp, _ := intersect(rLine, ul)
		ans += h * (s[i].Y - s[i].X + urp.X - ulp.X)
	}
	ans *= 2
	for i := range s {
		for j := 0; j <= i; j++ {
			pts := make([]pd, 0, 16)
			lp1 := pd{s[i].X, -h}
			rp1 := pd{s[i].Y, -h}
			l1 := getLine(cd, lp1)
			r1 := getLine(cd, rp1)
			ulp1, _ := intersect(l1, ul)
			urp1, _ := intersect(r1, ul)
			lp2 := pd{s[j].X, h}
			rp2 := pd{s[j].Y, h}
			l2 := getLine(cu, lp2)
			r2 := getLine(cu, rp2)
			dlp2, _ := intersect(l2, dl)
			drp2, _ := intersect(r2, dl)
			if ulp1.X >= lp2.X && ulp1.X <= rp2.X {
				pts = append(pts, ulp1)
			}
			if urp1.X >= lp2.X && urp1.X <= rp2.X {
				pts = append(pts, urp1)
			}
			if dlp2.X >= lp1.X && dlp2.X <= rp1.X {
				pts = append(pts, dlp2)
			}
			if drp2.X >= lp1.X && drp2.X <= rp1.X {
				pts = append(pts, drp2)
			}
			if lp2.X >= ulp1.X && lp2.X <= urp1.X {
				pts = append(pts, lp2)
			}
			if rp2.X >= ulp1.X && rp2.X <= urp1.X {
				pts = append(pts, rp2)
			}
			if lp1.X >= dlp2.X && lp1.X <= drp2.X {
				pts = append(pts, lp1)
			}
			if rp1.X >= dlp2.X && rp1.X <= drp2.X {
				pts = append(pts, rp1)
			}
			if cll, ok := intersect(l1, l2); ok && cll.Y >= -h && cll.Y <= h {
				pts = append(pts, cll)
			}
			if clr, ok := intersect(l1, r2); ok && clr.Y >= -h && clr.Y <= h {
				pts = append(pts, clr)
			}
			if crl, ok := intersect(r1, l2); ok && crl.Y >= -h && crl.Y <= h {
				pts = append(pts, crl)
			}
			if crr, ok := intersect(r1, r2); ok && crr.Y >= -h && crr.Y <= h {
				pts = append(pts, crr)
			}
			if i == j {
				ans -= getArea(pts)
			} else {
				ans -= 2 * getArea(pts)
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(3) + 1
	h := float64(rng.Intn(9) + 1)
	f := h + float64(rng.Intn(20)+1)
	segs := make([][2]float64, n)
	used := [][2]float64{}
	for i := 0; i < n; i++ {
		for {
			l := float64(rng.Intn(10) - 5)
			r := l + float64(rng.Intn(3)+1)
			overlap := false
			for _, u := range used {
				if !(r <= u[0] || l >= u[1]) {
					overlap = true
					break
				}
			}
			if !overlap {
				segs[i] = [2]float64{l, r}
				used = append(used, segs[i])
				break
			}
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, int(h), int(f))
	for _, s := range segs {
		fmt.Fprintf(&b, "%d %d\n", int(s[0]), int(s[1]))
	}
	expected := solveE(n, h, f, segs)
	return b.String(), expected
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		input, expected := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var got float64
		fmt.Sscan(strings.TrimSpace(out), &got)
		if math.Abs(got-expected) > 1e-4 {
			fmt.Printf("case %d failed\ninput:\n%sexpected %.6f got %.6f\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
