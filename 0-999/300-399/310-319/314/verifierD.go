package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceD = "314D.go"
	refBinaryD = "ref314D.bin"
	totalTests = 80
	tolerance  = 1e-6
)

type point struct {
	x int64
	y int64
}

type testCase struct {
	points []point
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}

		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if !closeEnough(refVal, candVal) {
			fmt.Printf("test %d failed: expected %.10f, got %.10f\n", idx+1, refVal, candVal)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryD, refSourceD)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryD), nil
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

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tc.points)))
	sb.WriteByte('\n')
	for _, p := range tc.points {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return []byte(sb.String())
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

func generateTests() []testCase {
	tests := []testCase{
		{points: []point{{0, 0}}},
		{points: []point{{0, 0}, {1, 0}}},
		{points: []point{{1, 2}, {-3, 4}, {5, -6}}},
		{points: []point{{-1, -1}, {-2, -3}, {-4, -5}}},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-3 {
		n := rnd.Intn(60) + 1
		points := make([]point, n)
		for i := 0; i < n; i++ {
			points[i] = point{
				x: randInt64(rnd, -1_000_000_000, 1_000_000_000),
				y: randInt64(rnd, -1_000_000_000, 1_000_000_000),
			}
		}
		tests = append(tests, testCase{points: points})
	}

	tests = append(tests, testCase{points: buildGridPoints(300)})
	tests = append(tests, testCase{points: buildCirclePoints(400)})
	tests = append(tests, testCase{points: buildRandomLarge(100000)})

	return tests
}

func randInt64(rnd *rand.Rand, lo, hi int64) int64 {
	return lo + rnd.Int63n(hi-lo+1)
}

func buildGridPoints(k int) []point {
	points := make([]point, 0, k*k)
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			points = append(points, point{x: int64(i), y: int64(j)})
		}
	}
	return points
}

func buildCirclePoints(n int) []point {
	points := make([]point, n)
	for i := 0; i < n; i++ {
		angle := float64(i) / float64(n) * 2 * math.Pi
		points[i] = point{
			x: int64(math.Round(1e6 * math.Cos(angle))),
			y: int64(math.Round(1e6 * math.Sin(angle))),
		}
	}
	return points
}

func buildRandomLarge(n int) []point {
	points := make([]point, n)
	rnd := rand.New(rand.NewSource(1234567))
	for i := 0; i < n; i++ {
		points[i] = point{
			x: randInt64(rnd, -1_000_000_000, 1_000_000_000),
			y: randInt64(rnd, -1_000_000_000, 1_000_000_000),
		}
	}
	return points
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
