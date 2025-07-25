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

func solve(n int, l, v1, v2 float64, k int) float64 {
	g := n / k
	if n%k != 0 {
		g++
	}
	fg := float64(g)
	numerator := l*v1 + (2*v2*fg-v2)*l
	denominator := (2*v2*fg-v2)*v1 + v2*v2
	return numerator / denominator
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(100) + 1
	k := rng.Intn(n) + 1
	l := float64(rng.Intn(1000) + 1)
	v1 := float64(rng.Intn(100) + 1)
	v2 := float64(rng.Intn(100) + int(v1) + 2) // ensure v2 > v1
	input := fmt.Sprintf("%d %.0f %.0f %.0f %d\n", n, l, v1, v2, k)
	expected := solve(n, l, v1, v2, k)
	return input, expected
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
