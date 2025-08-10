package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	input string
	x, y  float64
}

func rotate(k, xi, yi int) (float64, float64) {
	x := float64(xi)
	y := float64(yi)
	theta := float64(k) * math.Pi / 180.0
	c := math.Cos(theta)
	s := math.Sin(theta)
	x2 := x*c - y*s
	y2 := x*s + y*c
	return x2, y2
}

func genTests() []TestCase {
	r := rand.New(rand.NewSource(9))
	cases := make([]TestCase, 100)
	for i := 0; i < 100; i++ {
		k := r.Intn(360)
		xi := r.Intn(201) - 100
		yi := r.Intn(201) - 100
		input := fmt.Sprintf("%d\n%d %d\n", k, xi, yi)
		x2, y2 := rotate(k, xi, yi)
		cases[i] = TestCase{input: input, x: x2, y: y2}
	}
	return cases
}

func run(bin, in string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var gx, gy float64
		if _, err := fmt.Sscanf(got, "%f %f", &gx, &gy); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse output %q: %v\n", i+1, got, err)
			os.Exit(1)
		}
		const eps = 1e-6
		if math.Abs(gx-tc.x) > eps || math.Abs(gy-tc.y) > eps {
			fmt.Fprintf(os.Stderr, "test %d failed: input %q expected %.10f %.10f got %q\n", i+1, tc.input, tc.x, tc.y, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
