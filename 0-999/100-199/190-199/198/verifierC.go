package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	xx, yy, vv        float64
	x, y, v, rrGlobal float64
	tx, ty            float64
)

func getPos(m float64) float64 {
	qx := x + (tx-x)*m
	qy := y + (ty-y)*m
	return math.Hypot(qx, qy)
}

func getExt(x0, y0 float64) (float64, float64) {
	c := math.Hypot(x0, y0)
	a := rrGlobal
	b := math.Sqrt(c*c - a*a)
	return b, math.Acos(a / c)
}

func check(t float64) float64 {
	l2, r2 := 0.0, 1.0
	for i := 0; i < 200; i++ {
		m1 := (2*l2 + r2) / 3
		m2 := (l2 + 2*r2) / 3
		if getPos(m1) <= getPos(m2) {
			r2 = m2
		} else {
			l2 = m1
		}
	}
	if getPos(l2) >= rrGlobal {
		return math.Hypot(x-tx, y-ty)
	}
	spang := math.Abs(math.Atan2(y, x) - math.Atan2(ty, tx))
	if 2*math.Pi-spang < spang {
		spang = 2*math.Pi - spang
	}
	b1, ang1 := getExt(x, y)
	b2, ang2 := getExt(tx, ty)
	spang -= (ang1 + ang2)
	return b1 + b2 + spang*rrGlobal
}

func solveCase(xp, yp, vp, x1, y1, v1, r float64) float64 {
	xx, yy, vv = xp, yp, vp
	x, y, v, rrGlobal = x1, y1, v1, r
	stang := math.Atan2(yy, xx)
	rr := math.Hypot(xx, yy)
	l, rtime := 0.0, 1e7
	for i := 0; i < 200; i++ {
		m := (l + rtime) / 2
		nang := stang + vv*m/rr
		tx = math.Cos(nang) * rr
		ty = math.Sin(nang) * rr
		if check(m) <= v*m {
			rtime = m
		} else {
			l = m
		}
	}
	return rtime
}

func generateCase(rng *rand.Rand) (string, string) {
	xp := rng.Float64()*20 - 10
	yp := rng.Float64()*20 - 10
	vp := rng.Float64()*5 + 1
	x1 := rng.Float64()*20 - 10
	y1 := rng.Float64()*20 - 10
	r := rng.Float64()*5 + 1
	v1 := vp + rng.Float64()*5 + 0.1
	if math.Hypot(x1, y1) <= r {
		x1 += r
		y1 += r
	}
	if math.Hypot(xp, yp) <= r {
		xp += r
		yp += r
	}
	inp := fmt.Sprintf("%.4f %.4f %.4f\n%.4f %.4f %.4f %.4f\n", xp, yp, vp, x1, y1, v1, r)
	ans := solveCase(xp, yp, vp, x1, y1, v1, r)
	return inp, fmt.Sprintf("%.6f", ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if _, err := fmt.Sscanf(got, "%f", new(float64)); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	// compare with tolerance
	var gotVal float64
	fmt.Sscanf(got, "%f", &gotVal)
	var expVal float64
	fmt.Sscanf(expected, "%f", &expVal)
	if math.Abs(gotVal-expVal) > 1e-3 {
		return fmt.Errorf("expected %.6f got %s", expVal, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
