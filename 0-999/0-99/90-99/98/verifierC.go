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

func solve(A, B, L float64) string {
	// Always do ternary search to find minimum of f(x) = (A*x + B*y - x*y)/L
	// where y = sqrt(L^2 - x^2), over x in [0, L].
	// This gives the maximum width W of a coffin that can pass the L-shaped turn.
	f := func(x float64) float64 {
		y := math.Sqrt(math.Max(0, L*L-x*x))
		return (A*x + B*y - x*y) / L
	}

	left, right := 0.0, L
	for i := 0; i < 500; i++ {
		ma := (left*3 + right) / 4.0
		mb := (left + right*3) / 4.0
		if f(ma) < f(mb) {
			right = mb
		} else {
			left = ma
		}
	}
	ff := f((left + right) / 2.0)
	// Cap at corridor widths and stick length (width <= length)
	ff = math.Min(ff, L)
	ff = math.Min(ff, A)
	ff = math.Min(ff, B)
	if ff < 1e-8 {
		return "My poor head =("
	}
	return fmt.Sprintf("%.9f", ff)
}

func genCase(rng *rand.Rand) string {
	A := rng.Intn(10000) + 1
	B := rng.Intn(10000) + 1
	L := rng.Intn(10000) + 1
	return fmt.Sprintf("%d %d %d\n", A, B, L)
}

func runCase(bin string, input string, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	if expected == "My poor head =(" {
		if gotStr != expected {
			return fmt.Errorf("expected %s got %s", expected, gotStr)
		}
		return nil
	}
	var got float64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	var exp float64
	fmt.Sscan(expected, &exp)
	if math.Abs(got-exp) > 1e-4 {
		return fmt.Errorf("expected %s got %s", expected, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		var A, B, L float64
		fmt.Sscan(input, &A, &B, &L)
		expected := solve(A, B, L)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
