package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	a, b, m    int
	vx, vy, vz float64
}

// Embedded testcases (previously from testcasesD.txt) to keep verifier self contained.
const rawTestcasesD = `
4 10 9 -3 -3 4
8 10 2 4 -5 3
5 9 4 -2 -2 4
9 8 7 5 -4 1
3 9 7 -5 -5 1
10 9 8 0 -4 5
10 9 5 4 -1 2
3 7 4 1 -1 3
5 9 1 1 -4 1
2 6 7 -4 -2 1
1 4 8 -3 -1 1
1 2 7 -5 -5 5
8 9 7 0 -5 1
1 5 2 1 -1 2
3 9 10 -5 -3 1
9 10 2 2 -1 2
8 10 1 -3 -1 0
5 1 2 2 -1 1
2 5 10 1 -2 2
1 2 10 -2 -1 1
1 9 4 -1 -1 4
1 3 6 -1 -1 5
10 9 10 2 -5 0
10 6 10 2 -3 3
9 10 6 5 -2 4
1 2 4 -2 -4 5
6 2 4 -3 -5 4
1 10 7 2 -1 3
10 8 5 1 -5 2
6 10 1 -1 -1 4
1 9 9 -2 -2 2
1 9 1 1 -1 0
6 8 5 -1 -1 3
10 3 4 -2 -5 1
10 1 9 5 -1 2
6 7 10 -4 -2 3
5 3 6 -3 -1 0
3 6 3 -2 -3 0
2 6 3 -3 -2 4
2 7 7 2 -2 1
3 8 7 5 -2 4
2 7 6 2 -2 5
6 10 1 -2 -1 5
8 10 10 4 -4 0
1 2 2 2 -2 3
2 4 7 -2 -1 1
4 9 2 -1 -5 5
1 10 10 1 -3 1
5 7 4 0 -1 0
1 5 5 2 -1 1
8 8 6 2 -2 0
9 8 3 5 -1 4
`

func loadTestcases() ([]testCase, error) {
	fields := strings.Fields(rawTestcasesD)
	if len(fields)%6 != 0 {
		return nil, fmt.Errorf("unexpected token count %d (want multiple of 6)", len(fields))
	}
	cases := make([]testCase, 0, len(fields)/6)
	for i := 0; i < len(fields); i += 6 {
		a, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("parse a at token %d (%q): %w", i+1, fields[i], err)
		}
		b, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("parse b at token %d (%q): %w", i+2, fields[i+1], err)
		}
		m, err := strconv.Atoi(fields[i+2])
		if err != nil {
			return nil, fmt.Errorf("parse m at token %d (%q): %w", i+3, fields[i+2], err)
		}
		vx, err := strconv.ParseFloat(fields[i+3], 64)
		if err != nil {
			return nil, fmt.Errorf("parse vx at token %d (%q): %w", i+4, fields[i+3], err)
		}
		vy, err := strconv.ParseFloat(fields[i+4], 64)
		if err != nil {
			return nil, fmt.Errorf("parse vy at token %d (%q): %w", i+5, fields[i+4], err)
		}
		vz, err := strconv.ParseFloat(fields[i+5], 64)
		if err != nil {
			return nil, fmt.Errorf("parse vz at token %d (%q): %w", i+6, fields[i+5], err)
		}
		cases = append(cases, testCase{a: a, b: b, m: m, vx: vx, vy: vy, vz: vz})
	}
	return cases, nil
}

// solve203DCase uses the analytical unfolding formula (same approach as the reference solution).
func solve203DCase(tc testCase) (float64, float64) {
	aF := float64(tc.a)
	bF := float64(tc.b)

	t := -float64(tc.m) / tc.vy

	X := aF/2.0 + tc.vx*t
	Z := tc.vz * t

	Xrem := math.Mod(X, 2*aF)
	if Xrem < 0 {
		Xrem += 2 * aF
	}
	x0 := Xrem
	if Xrem > aF {
		x0 = 2*aF - Xrem
	}

	Zrem := math.Mod(Z, 2*bF)
	if Zrem < 0 {
		Zrem += 2 * bF
	}
	z0 := Zrem
	if Zrem > bF {
		z0 = 2*bF - Zrem
	}

	return x0, z0
}

func parseOutput(s string) (float64, float64, error) {
	parts := strings.Fields(s)
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("output should contain two numbers")
	}
	x0, err1 := strconv.ParseFloat(parts[0], 64)
	z0, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid floats")
	}
	return x0, z0, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		expX, expZ := solve203DCase(tc)
		input := fmt.Sprintf("%d %d %d %.10g %.10g %.10g\n", tc.a, tc.b, tc.m, tc.vx, tc.vy, tc.vz)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		gotX, gotZ, err := parseOutput(strings.TrimSpace(out.String()))
		if err != nil {
			fmt.Printf("test %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if math.Abs(gotX-expX) > 1e-6 || math.Abs(gotZ-expZ) > 1e-6 {
			fmt.Printf("test %d failed\nexpected: %.6f %.6f\n got: %.6f %.6f\n", idx+1, expX, expZ, gotX, gotZ)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
