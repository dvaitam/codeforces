package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource  = "./1184C2.go"
	refBinary  = "ref1184C2.bin"
	totalTests = 80
)

type point struct {
	x int
	y int
}

type testCase struct {
	n      int
	r      int
	points []point
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if candVal != refVal {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, refVal, candVal, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref-1184C2-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1184C2.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.r))
	for _, p := range tc.points {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 1, r: 1, points: []point{{0, 0}}},
		{n: 2, r: 1, points: []point{{0, 0}, {1, 0}}},
		{n: 3, r: 2, points: []point{{0, 0}, {1, 1}, {2, 0}}},
		{n: 4, r: 3, points: []point{{-1, -1}, {1, 1}, {2, -2}, {-2, 2}}},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-2 {
		n := rng.Intn(40) + 1
		r := rng.Intn(50) + 1
		points := make([]point, n)
		used := make(map[[2]int]bool)
		for i := 0; i < n; i++ {
			for {
				x := rng.Intn(201) - 100
				y := rng.Intn(201) - 100
				key := [2]int{x, y}
				if !used[key] {
					used[key] = true
					points[i] = point{x, y}
					break
				}
			}
		}
		tests = append(tests, testCase{n: n, r: r, points: points})
	}
	tests = append(tests,
		heavyCase(1000, 500, rand.New(rand.NewSource(1))),
		heavyCase(300000, 1000000, rand.New(rand.NewSource(2))),
	)
	return tests
}

func heavyCase(n int, r int, rng *rand.Rand) testCase {
	points := make([]point, n)
	used := make(map[[2]int]bool, n)
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(2_000_001) - 1_000_000
			y := rng.Intn(2_000_001) - 1_000_000
			key := [2]int{x, y}
			if !used[key] {
				used[key] = true
				points[i] = point{x, y}
				break
			}
		}
	}
	return testCase{n: n, r: r, points: points}
}
