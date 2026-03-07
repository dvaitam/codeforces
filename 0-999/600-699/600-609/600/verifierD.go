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
	"time"
)

const (
	totalTests = 120
	tolerance  = 1e-6
)

type circleCase struct {
	x1 int64
	y1 int64
	r1 int64
	x2 int64
	y2 int64
	r2 int64
}

func circleIntersectionArea(tc circleCase) float64 {
	dx := float64(tc.x1 - tc.x2)
	dy := float64(tc.y1 - tc.y2)
	d := math.Sqrt(dx*dx + dy*dy)
	R := float64(tc.r1)
	r := float64(tc.r2)
	pi := math.Pi
	if d >= R+r {
		return 0
	}
	if d+r <= R {
		return pi * r * r
	}
	if d+R <= r {
		return pi * R * R
	}
	a := math.Acos((R*R+d*d-r*r)/(2*R*d)) * 2
	b := math.Acos((r*r+d*d-R*R)/(2*r*d)) * 2
	return 0.5 * (R*R*(a-math.Sin(a)) + r*r*(b-math.Sin(b)))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	if strings.HasSuffix(candidate, ".go") {
		tmp, err := os.CreateTemp("", "verifierD-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("go", "build", "-o", tmp.Name(), candidate).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		candidate = tmp.Name()
	}

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		expected := circleIntersectionArea(tc)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if !closeEnough(expected, candVal) {
			fmt.Printf("test %d failed: expected %.10f, got %.10f\n", idx+1, expected, candVal)
			printInput(input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func formatInput(tc circleCase) []byte {
	return []byte(fmt.Sprintf("%d %d %d %d %d %d\n", tc.x1, tc.y1, tc.r1, tc.x2, tc.y2, tc.r2))
}

func parseOutput(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", fields[0], err)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("non-finite value %v", val)
	}
	return val, nil
}

func closeEnough(expected, actual float64) bool {
	diff := math.Abs(expected - actual)
	allowed := tolerance * math.Max(1.0, math.Abs(expected))
	return diff <= allowed+1e-12
}

func generateTests() []circleCase {
	tests := []circleCase{
		{0, 0, 1, 0, 0, 1},
		{0, 0, 5, 10, 0, 5},
		{0, 0, 10, 1, 1, 1},
		{100, 100, 50, 120, 100, 50},
		{-1_000_000_000, -1_000_000_000, 1_000_000_000, 1_000_000_000, 1_000_000_000, 1_000_000_000},
		{0, 0, 1_000_000_000, 1_000_000_000, 0, 1_000_000_000},
		{0, 0, 123456789, 987654321, -987654321, 1},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		tests = append(tests, randomCase(rnd))
	}
	// deterministic stress cases
	tests = append(tests, circleCase{0, 0, 1_000_000_000, 0, 1, 999_999_999})
	tests = append(tests, circleCase{500_000_000, -500_000_000, 750_000_000, -500_000_000, 500_000_000, 250_000_000})
	tests = append(tests, circleCase{-123456789, 987654321, 345678901, 123456789, -987654321, 234567890})
	tests = append(tests, circleCase{1, 2, 3, 4, 5, 6})
	tests = append(tests, circleCase{-999_999_999, 999_999_999, 999_999_999, 999_999_999, -999_999_999, 1})

	return tests
}

func randomCase(rnd *rand.Rand) circleCase {
	return circleCase{
		x1: randInt64(rnd, -1_000_000_000, 1_000_000_000),
		y1: randInt64(rnd, -1_000_000_000, 1_000_000_000),
		r1: randInt64(rnd, 1, 1_000_000_000),
		x2: randInt64(rnd, -1_000_000_000, 1_000_000_000),
		y2: randInt64(rnd, -1_000_000_000, 1_000_000_000),
		r2: randInt64(rnd, 1, 1_000_000_000),
	}
}

func randInt64(rnd *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + rnd.Int63n(hi-lo+1)
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
