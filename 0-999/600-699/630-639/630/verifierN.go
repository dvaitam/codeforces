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

func solve(a, b, c int) (float64, float64) {
	da := float64(b*b - 4*a*c)
	d := math.Sqrt(da)
	af := float64(a)
	bf := float64(b)
	x1 := (-bf + d) / (2 * af)
	x2 := (-bf - d) / (2 * af)
	if x1 < x2 {
		x1, x2 = x2, x1
	}
	return x1, x2
}

func generateCase(rng *rand.Rand) (string, string) {
	var a, b, c int
	for {
		a = rng.Intn(2001) - 1000
		if a != 0 {
			break
		}
	}
	for {
		b = rng.Intn(2001) - 1000
		c = rng.Intn(2001) - 1000
		if float64(b*b-4*a*c) > 0 {
			break
		}
	}
	x1, x2 := solve(a, b, c)
	input := fmt.Sprintf("%d %d %d\n", a, b, c)
	expected := fmt.Sprintf("%.6f\n%.6f", x1, x2)
	return input, expected
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
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierN.go /path/to/binary")
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
