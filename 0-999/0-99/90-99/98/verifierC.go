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

func f(x, A, B, L float64) float64 {
	y := math.Sqrt(L*L - x*x)
	return (A*x + B*y - x*y) / L
}

func solve(A, B, L float64) string {
	if L <= B {
		v := math.Min(A, L)
		return fmt.Sprintf("%.8f", v)
	} else if L <= A {
		v := math.Min(B, L)
		return fmt.Sprintf("%.8f", v)
	}
	left, right := 0.0, L
	var ma, mb float64
	for i := 0; i < 500; i++ {
		ma = (left*3 + right) / 4
		mb = (left + right*3) / 4
		fa := f(ma, A, B, L)
		fb := f(mb, A, B, L)
		if fa < fb {
			right = mb
		} else {
			left = ma
		}
	}
	ff := f((left+right)/2, A, B, L)
	ff = math.Min(ff, L)
	ff = math.Min(ff, A)
	if ff < 1e-8 {
		return "My poor head =("
	}
	return fmt.Sprintf("%.8f", ff)
}

func genCase(rng *rand.Rand) string {
	A := float64(rng.Intn(10000) + 1)
	B := float64(rng.Intn(10000) + 1)
	L := float64(rng.Intn(10000) + 1)
	return fmt.Sprintf("%f %f %f\n", A, B, L)
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
	if math.Abs(got-exp) > 1e-6 {
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
