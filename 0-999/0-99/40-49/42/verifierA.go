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

func computeVolume(n int, V float64, a, b []float64) float64 {
	sumA := 0.0
	for _, v := range a {
		sumA += v
	}
	maxPan := V / sumA
	maxIng := math.Inf(1)
	for i := 0; i < n; i++ {
		ratio := b[i] / a[i]
		if ratio < maxIng {
			maxIng = ratio
		}
	}
	x := maxPan
	if maxIng < x {
		x = maxIng
	}
	return sumA * x
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(5) + 1
	V := float64(rng.Intn(20) + 1)
	a := make([]float64, n)
	b := make([]float64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %.0f\n", n, V))
	for i := 0; i < n; i++ {
		a[i] = float64(rng.Intn(10) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%.0f", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		b[i] = float64(rng.Intn(20))
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%.0f", b[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), computeVolume(n, V, a, b)
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
		return fmt.Errorf("bad output: %v", err)
	}
	diff := math.Abs(got - expected)
	tol := 1e-4 * math.Max(1, math.Abs(expected))
	if diff > tol {
		return fmt.Errorf("expected %.5f got %.5f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
