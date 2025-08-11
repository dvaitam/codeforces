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

type point struct{ x, y float64 }

func pointLineDistance(p, a, b point) float64 {
	num := math.Abs((p.x-a.x)*(b.y-a.y) - (p.y-a.y)*(b.x-a.x))
	den := math.Hypot(b.x-a.x, b.y-a.y)
	if den == 0 {
		return math.Hypot(p.x-a.x, p.y-a.y)
	}
	return num / den
}

func solveCase(pts []point) string {
	n := len(pts)
	ans := math.MaxFloat64
	for i := 0; i < n; i++ {
		a := pts[(i-1+n)%n]
		b := pts[i]
		c := pts[(i+1)%n]

		altA := pointLineDistance(a, b, c)
		altB := pointLineDistance(b, a, c)
		altC := pointLineDistance(c, a, b)

		d := math.Min(altA, math.Min(altB, altC)) / 2
		if d < ans {
			ans = d
		}
	}
	return fmt.Sprintf("%.10f\n", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 4
	angles := make([]float64, n)
	for i := range angles {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	pts := make([]point, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, ang := range angles {
		r := rng.Float64()*10 + 1
		x := r * math.Cos(ang)
		y := r * math.Sin(ang)
		xi := int(x)
		yi := int(y)
		pts[i] = point{float64(xi), float64(yi)}
		fmt.Fprintf(&sb, "%d %d\n", xi, yi)
	}
	exp := solveCase(pts)
	return sb.String(), exp
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	var g, e float64
	fmt.Sscan(got, &g)
	fmt.Sscan(strings.TrimSpace(exp), &e)
	if math.Abs(g-e) > 1e-6 {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(exp), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
