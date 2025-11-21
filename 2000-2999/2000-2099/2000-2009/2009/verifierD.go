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
)

const (
	refSource2009D = "2009D.go"
	refBinary2009D = "ref2009D.bin"
	maxTests       = 160
	maxTotalN      = 200000
)

type point struct {
	x int
	y int
}

type testCase struct {
	points []point
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on case %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2009D, refSource2009D)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2009D), nil
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

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", len(tc.points))
		for _, p := range tc.points {
			fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2009))
	var tests []testCase
	totalN := 0

	add := func(points []point) {
		tests = append(tests, testCase{points: points})
		totalN += len(points)
	}

	// Small hand-crafted cases
	add([]point{{0, 0}, {1, 0}, {0, 1}})
	add([]point{{0, 0}, {2, 0}, {1, 1}, {3, 1}})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain == 0 {
			break
		}
		maxN := 5000
		if remain < maxN {
			maxN = remain
		}
		n := rnd.Intn(maxN-2) + 3
		points := make([]point, 0, n)
		used := make(map[[2]int]struct{})
		for len(points) < n {
			x := rnd.Intn(n + 1)
			y := rnd.Intn(2)
			key := [2]int{x, y}
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			points = append(points, point{x: x, y: y})
		}
		add(points)
	}
	return tests
}
