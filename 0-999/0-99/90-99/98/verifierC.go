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

// solve implements the correct reference solver for 98C.
// Given a, b, l (sides of a room and length of a board),
// compute the minimum width that can pass through an L-shaped corridor.
func solve(A, B, L float64) string {
	f := func(theta float64) float64 {
		return A*math.Sin(theta) + B*math.Cos(theta) - L*math.Sin(theta)*math.Cos(theta)
	}

	intervals := 10000
	minW := L

	for i := 0; i < intervals; i++ {
		left := (math.Pi / 2) * float64(i) / float64(intervals)
		right := (math.Pi / 2) * float64(i+1) / float64(intervals)

		for j := 0; j < 80; j++ {
			m1 := left + (right-left)/3.0
			m2 := right - (right-left)/3.0
			if f(m1) < f(m2) {
				right = m2
			} else {
				left = m1
			}
		}
		if f(left) < minW {
			minW = f(left)
		}
	}

	if minW > 1e-8 {
		return fmt.Sprintf("%.9f", minW)
	}
	return "My poor head =("
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
