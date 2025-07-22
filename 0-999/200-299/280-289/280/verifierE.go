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

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func solve(xs []float64, q, a, b float64) ([]float64, float64) {
	n := len(xs)
	y := make([]float64, n)
	L := make([]float64, n)
	U := make([]float64, n)
	for i := 0; i < n; i++ {
		L[i] = 1 + float64(i)*a
		U[i] = q - float64(n-1-i)*a
	}
	for i := 0; i < n; i++ {
		y[i] = clamp(xs[i], L[i], U[i])
	}
	for it := 0; it < 200; it++ {
		maxd := 0.0
		for i := 0; i+1 < n; i++ {
			lo := y[i] + a
			hi := y[i] + b
			old := y[i+1]
			y[i+1] = clamp(y[i+1], lo, hi)
			if d := math.Abs(y[i+1] - old); d > maxd {
				maxd = d
			}
		}
		for i := n - 1; i > 0; i-- {
			hi := y[i] - a
			lo := y[i] - b
			old := y[i-1]
			y[i-1] = clamp(y[i-1], math.Max(lo, L[i-1]), math.Min(hi, U[i-1]))
			if d := math.Abs(y[i-1] - old); d > maxd {
				maxd = d
			}
		}
		if maxd < 1e-9 {
			break
		}
	}
	cost := 0.0
	for i := 0; i < n; i++ {
		d := y[i] - xs[i]
		cost += d * d
	}
	return y, cost
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	q := rng.Float64()*20 + 5
	a := rng.Float64() * 2
	b := a + rng.Float64()*2
	xs := make([]float64, n)
	for i := range xs {
		xs[i] = rng.Float64() * q
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %.6f %.6f %.6f\n", n, q, a, b))
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%.6f", v))
	}
	sb.WriteByte('\n')
	ys, cost := solve(xs, q, a, b)
	var exp strings.Builder
	for i, v := range ys {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%.9f", v))
	}
	exp.WriteByte('\n')
	exp.WriteString(fmt.Sprintf("%.9f", cost))
	exp.WriteByte('\n')
	return sb.String(), exp.String()
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
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
