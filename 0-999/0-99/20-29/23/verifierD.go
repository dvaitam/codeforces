package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func cp(x0, y0, x1, y1, x2, y2 float64) float64 {
	return (x1-x0)*(y2-y0) - (x2-x0)*(y1-y0)
}

func check(x1, y1, x2, y2, x3, y3 float64) (bool, [8]float64) {
	sqr := func(x float64) float64 { return x * x }
	a1 := (x2 - x1) * 2.0
	b1 := (y2 - y1) * 2.0
	c1 := sqr(2*x1-x2) + sqr(2*y1-y2) - sqr(x1) - sqr(y1)
	a2 := (x3 - 2*x2 + x1) * 2.0
	b2 := (y3 - 2*y2 + y1) * 2.0
	c2 := sqr(x1) + sqr(y1) - sqr(x3-2*x2+2*x1) - sqr(y3-2*y2+2*y1)
	if a1*b2 == a2*b1 {
		return false, [8]float64{}
	}
	Y1 := (c2*a1 - c1*a2) / (b1*a2 - b2*a1)
	X1 := (c2*b1 - c1*b2) / (a1*b2 - a2*b1)
	X2 := 2*x1 - X1
	Y2 := 2*y1 - Y1
	X3 := 2*x2 - 2*x1 + X1
	Y3 := 2*y2 - 2*y1 + Y1
	X4 := 2*x3 - 2*x2 + 2*x1 - X1
	Y4 := 2*y3 - 2*y2 + 2*y1 - Y1
	v1 := cp(X1, Y1, X2, Y2, X3, Y3)
	v2 := cp(X2, Y2, X3, Y3, X4, Y4)
	v3 := cp(X3, Y3, X4, Y4, X1, Y1)
	v4 := cp(X4, Y4, X1, Y1, X2, Y2)
	if (v1 < 0 && v2 < 0 && v3 < 0 && v4 < 0) || (v1 > 0 && v2 > 0 && v3 > 0 && v4 > 0) {
		return true, [8]float64{X1, Y1, X2, Y2, X3, Y3, X4, Y4}
	}
	return false, [8]float64{}
}

func solveOne(x1, y1, x2, y2, x3, y3 float64) (bool, [8]float64) {
	if ok, p := check(x1, y1, x2, y2, x3, y3); ok {
		return true, p
	}
	if ok, p := check(x1, y1, x3, y3, x2, y2); ok {
		return true, p
	}
	if ok, p := check(x2, y2, x1, y1, x3, y3); ok {
		return true, p
	}
	return false, [8]float64{}
}

func solveD(points [][6]float64) string {
	var sb strings.Builder
	for idx, pt := range points {
		if idx > 0 {
			sb.WriteByte('\n')
		}
		found, res := solveOne(pt[0], pt[1], pt[2], pt[3], pt[4], pt[5])
		if found {
			sb.WriteString("YES\n")
			fmt.Fprintf(&sb, "%.9f %.9f %.9f %.9f %.9f %.9f %.9f %.9f", res[0], res[1], res[2], res[3], res[4], res[5], res[6], res[7])
		} else {
			sb.WriteString("NO\n\n")
			continue
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	pts := make([][6]float64, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		for j := 0; j < 6; j++ {
			pts[i][j] = float64(rng.Intn(11))
		}
		fmt.Fprintf(&sb, "%d %d %d %d %d %d\n", int(pts[i][0]), int(pts[i][1]), int(pts[i][2]), int(pts[i][3]), int(pts[i][4]), int(pts[i][5]))
	}
	return sb.String(), solveD(pts)
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
