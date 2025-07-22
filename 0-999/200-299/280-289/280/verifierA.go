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

func expected(w, h, deg int) float64 {
	x := float64(w)
	y := float64(h)
	if x < y {
		x, y = y, x
	}
	a := float64(deg) / 180 * math.Pi
	if a > math.Pi/2 {
		a = math.Pi - a
	}
	diag := math.Sqrt(x*x + y*y)
	maxAngle := math.Asin(y/diag) * 2
	halfA := a / 2
	var ans float64
	if a >= maxAngle && a <= math.Pi-maxAngle {
		l := y / math.Sin(halfA)
		hh := l / 2 * math.Tan(halfA)
		ans = l * hh
	} else {
		l := y * math.Sin(halfA)
		ans = l * l / math.Tan(halfA)
		g := math.Cos(a) + math.Sin(a) + 1
		gg := math.Cos(a) - math.Sin(a) + 1
		xv := 0.0
		if gg != 0 {
			xv = (x - (x+y)/g*math.Sin(a)) / gg
		}
		yv := 0.0
		if g != 0 {
			yv = (x+y)/g - xv
		}
		s1 := xv * math.Sin(halfA) * xv * math.Cos(halfA) * 2
		s2 := yv * math.Sin(halfA) * yv * math.Cos(halfA) * 2
		s3 := yv * math.Cos(halfA) * 2 * xv * math.Cos(halfA) * 2
		ans = s1 + s2 + s3
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, float64) {
	w := rng.Intn(1_000_000) + 1
	h := rng.Intn(1_000_000) + 1
	a := rng.Intn(181)
	input := fmt.Sprintf("%d %d %d\n", w, h, a)
	return input, expected(w, h, a)
}

func runCase(bin, input string, expected float64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	diff := math.Abs(got - expected)
	tol := 1e-6 * math.Max(1, math.Abs(expected))
	if diff > tol {
		return fmt.Errorf("expected %.7f got %.7f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
