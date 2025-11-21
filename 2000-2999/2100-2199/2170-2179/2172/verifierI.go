package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type point struct {
	x int
	y int
}

type testCase struct {
	name string
	n    int
	r    int
	pts  []point
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2172I-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleI")
	cmd := exec.Command("go", "build", "-o", bin, "2172I.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.r))
	for _, p := range tc.pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "single_point", n: 1, r: 5, pts: []point{{0, 0}}},
		{name: "two_opposite", n: 2, r: 10, pts: []point{{3, 4}, {-3, -4}}},
		{name: "triangle", n: 3, r: 15, pts: []point{{4, 0}, {-2, 3}, {-2, -3}}},
		{name: "square", n: 4, r: 20, pts: []point{{6, 0}, {0, 6}, {-6, 0}, {0, -6}}},
	}
}

func randomPoint(r int, rng *rand.Rand) point {
	limit := int(float64(r) * 0.9)
	for {
		x := rng.Intn(2*limit+1) - limit
		y := rng.Intn(2*limit+1) - limit
		if x == 0 && y == 0 {
			continue
		}
		if float64(x*x+y*y) <= 0.81*float64(r*r) {
			return point{x, y}
		}
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 10)
	for i := 0; i < 10; i++ {
		n := rng.Intn(10) + 1
		r := rng.Intn(100) + 1
		pts := make([]point, n)
		for j := 0; j < n; j++ {
			pts[j] = randomPoint(r, rng)
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_small_%d", i+1),
			n:    n,
			r:    r,
			pts:  pts,
		})
	}
	for i := 0; i < 10; i++ {
		n := rng.Intn(5000) + 1
		r := rng.Intn(1_000_000) + 1
		if n > 200000 {
			n = 200000
		}
		pts := make([]point, n)
		for j := 0; j < n; j++ {
			pts[j] = randomPoint(r, rng)
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_large_%d", i+1),
			n:    n,
			r:    r,
			pts:  pts,
		})
	}
	return tests
}

func stressTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 123))
	tests := make([]testCase, 0, 3)
	tests = append(tests, testCase{
		name: "stress_line",
		n:    200000,
		r:    1_000_000,
		pts: func() []point {
			pts := make([]point, 200000)
			for i := 0; i < len(pts); i++ {
				x := -900000 + i*10
				if x > 900000 {
					x = 900000
				}
				pts[i] = point{x, 0}
			}
			return pts
		}(),
	})
	tests = append(tests, testCase{
		name: "stress_circle",
		n:    200000,
		r:    1_000_000,
		pts: func() []point {
			pts := make([]point, 200000)
			for i := 0; i < len(pts); i++ {
				angle := 2 * math.Pi * float64(i) / 200000.0
				rad := 0.9 * float64(1_000_000)
				x := int(rad * math.Cos(angle))
				y := int(rad * math.Sin(angle))
				pts[i] = point{x, y}
			}
			return pts
		}(),
	})
	tests = append(tests, testCase{
		name: "stress_clumped",
		n:    200000,
		r:    1_000_000,
		pts: func() []point {
			pts := make([]point, 200000)
			for i := 0; i < len(pts); i++ {
				pts[i] = randomPoint(1_000_000, rng)
			}
			return pts
		}(),
	})
	return tests
}

func compareFloats(expected, actual string) error {
	const eps = 1e-6
	expVal, err := strconv.ParseFloat(expected, 64)
	if err != nil {
		return fmt.Errorf("failed to parse expected value %q: %v", expected, err)
	}
	actVal, err := strconv.ParseFloat(actual, 64)
	if err != nil {
		return fmt.Errorf("failed to parse actual value %q: %v", actual, err)
	}
	diff := math.Abs(expVal - actVal)
	den := math.Max(1.0, math.Abs(expVal))
	if diff/den > eps {
		return fmt.Errorf("values differ: expected %.12f got %.12f (diff %.12g)", expVal, actVal, diff/den)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	tests = append(tests, stressTests()...)

	for idx, tc := range tests {
		input := buildInput(tc)
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		if err := compareFloats(expected, actual); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sexpected:\n%s\nactual:\n%s\n", idx+1, tc.name, err, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
