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
	refSource  = "1000-1999/1600-1699/1660-1669/1662/1662K.go"
	totalTests = 80
	tolerance  = 1e-4
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierK.go /path/to/binary")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if !closeEnough(refVal, candVal) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %.10f, got %.10f\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, tc.name, refVal, candVal, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref1662K-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1662K.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
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
	return val, nil
}

func closeEnough(exp, got float64) bool {
	diff := math.Abs(exp - got)
	den := math.Max(1.0, math.Abs(exp))
	return diff <= tolerance*den+1e-9
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "axis_triangle", input: "0 0\n5 0\n0 5\n"},
		{name: "colinear", input: "-5 0\n0 0\n5 0\n"},
		{name: "sample_like", input: "0 0\n5 0\n3 3\n"},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests {
		tests = append(tests, randomCase(rng, len(tests)+1))
	}
	return tests
}

func randomCase(rng *rand.Rand, idx int) testCase {
	coords := make([][2]int, 3)
	used := make(map[[2]int]bool)
	for i := 0; i < 3; i++ {
		for {
			x := rng.Intn(20001) - 10000
			y := rng.Intn(20001) - 10000
			key := [2]int{x, y}
			if !used[key] {
				used[key] = true
				coords[i] = key
				break
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", coords[i][0], coords[i][1]))
	}
	return testCase{name: fmt.Sprintf("rand_%d", idx), input: sb.String()}
}
