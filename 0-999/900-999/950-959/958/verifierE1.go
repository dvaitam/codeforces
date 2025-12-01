package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	refSource        = "958E1.go"
	tempOraclePrefix = "oracle-958E1-"
	randomTestsCount = 120
	maxCoord         = 10000
)

type point struct {
	x int
	y int
}

type testCase struct {
	name  string
	R, B  int
	ships []point
	bases []point
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)

	for idx, tc := range tests {
		input := formatInput(tc)

		expOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expected, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, gotOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if expected != got {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %v got %v\n", idx+1, tc.name, expected, got)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE1")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.R, tc.B)
	for _, p := range tc.ships {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	for _, p := range tc.bases {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	return sb.String()
}

func parseAnswer(out string) (bool, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return false, fmt.Errorf("empty output")
	}
	switch strings.ToLower(fields[0]) {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected answer %q", fields[0])
	}
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name:  "mismatched_counts",
			R:     2,
			B:     3,
			ships: []point{{0, 0}, {2, 0}},
			bases: []point{{0, 1}, {2, 1}, {4, 1}},
		},
		{
			name:  "vertical_pairs",
			R:     3,
			B:     3,
			ships: []point{{0, 0}, {2, 0}, {4, 0}},
			bases: []point{{0, 2}, {2, 2}, {4, 2}},
		},
		{
			name:  "crossing_required",
			R:     2,
			B:     2,
			ships: []point{{0, 0}, {2, 0}},
			bases: []point{{2, 2}, {0, -2}},
		},
		{
			name:  "triangle_configuration",
			R:     3,
			B:     3,
			ships: []point{{0, 0}, {4, 0}, {2, 3}},
			bases: []point{{0, 4}, {4, 4}, {2, 7}},
		},
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		R := rng.Intn(10) + 1
		B := rng.Intn(10) + 1
		ships, bases := randomPoints(R, B, rng)
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			R:     R,
			B:     B,
			ships: ships,
			bases: bases,
		})
	}
	return tests
}

func randomPoints(R, B int, rng *rand.Rand) ([]point, []point) {
	total := R + B
	points := make([]point, 0, total)
	for len(points) < total {
		p := point{
			x: rng.Intn(2*maxCoord+1) - maxCoord,
			y: rng.Intn(2*maxCoord+1) - maxCoord,
		}
		if exists(points, p) || collinearWithAny(points, p) {
			continue
		}
		points = append(points, p)
	}
	return points[:R], points[R:]
}

func exists(points []point, p point) bool {
	for _, q := range points {
		if q == p {
			return true
		}
	}
	return false
}

func collinearWithAny(points []point, p point) bool {
	n := len(points)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if cross(points[i], points[j], p) == 0 {
				return true
			}
		}
	}
	return false
}

func cross(a, b, c point) int {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}
