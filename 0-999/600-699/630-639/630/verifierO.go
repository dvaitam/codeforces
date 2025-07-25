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

func solve(px, py, vx, vy float64, a, b, c, d int) string {
	t := math.Sqrt(vx*vx + vy*vy)
	vx /= t
	vy /= t
	vx1 := vy
	vy1 := -vx
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%.10f %.10f\n", px+vx*float64(b), py+vy*float64(b)))
	sb.WriteString(fmt.Sprintf("%.10f %.10f\n", px-vx1*float64(a)/2, py-vy1*float64(a)/2))
	sb.WriteString(fmt.Sprintf("%.10f %.10f\n", px-vx1*float64(c)/2, py-vy1*float64(c)/2))
	sb.WriteString(fmt.Sprintf("%.10f %.10f\n", px-vx1*float64(c)/2-vx*float64(d), py-vy1*float64(c)/2-vy*float64(d)))
	sb.WriteString(fmt.Sprintf("%.10f %.10f\n", px+vx1*float64(c)/2-vx*float64(d), py+vy1*float64(c)/2-vy*float64(d)))
	sb.WriteString(fmt.Sprintf("%.10f %.10f\n", px+vx1*float64(c)/2, py+vy1*float64(c)/2))
	sb.WriteString(fmt.Sprintf("%.10f %.10f", px+vx1*float64(a)/2, py+vy1*float64(a)/2))
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	px := rng.Float64()*2000 - 1000
	py := rng.Float64()*2000 - 1000
	vx := rng.Float64()*2000 - 1000
	vy := rng.Float64()*2000 - 1000
	if vx == 0 && vy == 0 {
		vx = 1
	}
	a := rng.Intn(1000) + 1
	c := rng.Intn(a-1) + 1
	b := rng.Intn(1000) + 1
	d := rng.Intn(1000) + 1
	input := fmt.Sprintf("%.0f %.0f %.0f %.0f %d %d %d %d\n", px, py, vx, vy, a, b, c, d)
	output := solve(px, py, vx, vy, a, b, c, d)
	return input, output
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
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected\n%s\ngot\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierO.go /path/to/binary")
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
