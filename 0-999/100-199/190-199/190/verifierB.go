package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func solve(a, b, r1, x, y, r2 float64) string {
	dx := a - x
	dy := b - y
	d := math.Hypot(dx, dy)
	var res float64
	if d > r1+r2 {
		res = (d - r1 - r2) / 2
	} else if d+math.Min(r1, r2) < math.Max(r1, r2) {
		res = (math.Max(r1, r2) - math.Min(r1, r2) - d) / 2
	} else {
		res = 0
	}
	return fmt.Sprintf("%.12f", res)
}

func generateCases() []testCase {
	rand.Seed(2)
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		a := float64(rand.Intn(21) - 10)
		b := float64(rand.Intn(21) - 10)
		x := float64(rand.Intn(21) - 10)
		y := float64(rand.Intn(21) - 10)
		// ensure centers not equal
		if a == x && b == y {
			x += 1
		}
		r1 := float64(rand.Intn(10) + 1)
		r2 := float64(rand.Intn(10) + 1)
		buf := bytes.Buffer{}
		fmt.Fprintf(&buf, "%.0f %.0f %.0f\n", a, b, r1)
		fmt.Fprintf(&buf, "%.0f %.0f %.0f\n", x, y, r2)
		cases[i] = testCase{input: buf.String(), expected: solve(a, b, r1, x, y, r2)}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
