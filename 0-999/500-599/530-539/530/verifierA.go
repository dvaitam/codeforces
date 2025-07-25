package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(a, b, c float64) []float64 {
	d := b*b - 4*a*c
	sqrtD := math.Sqrt(d)
	x1 := (-b - sqrtD) / (2 * a)
	x2 := (-b + sqrtD) / (2 * a)
	if d == 0 {
		return []float64{x1}
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	return []float64{x1, x2}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		var a, b, c int
		for {
			a = rng.Intn(21) - 10
			if a != 0 {
				break
			}
		}
		b = rng.Intn(21) - 10
		c = rng.Intn(21) - 10
		for float64(b*b-4*a*c) < 0 {
			for {
				a = rng.Intn(21) - 10
				if a != 0 {
					break
				}
			}
			b = rng.Intn(21) - 10
			c = rng.Intn(21) - 10
		}
		input := fmt.Sprintf("%d %d %d\n", a, b, c)
		want := expected(float64(a), float64(b), float64(c))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", i+1, len(want), len(fields))
			os.Exit(1)
		}
		for j, f := range fields {
			got, err := strconv.ParseFloat(f, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid float output\n", i+1)
				os.Exit(1)
			}
			diff := math.Abs(got - want[j])
			tol := 1e-4 * math.Max(1, math.Abs(want[j]))
			if diff > tol {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %.10f got %.10f\n", i+1, want[j], got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
