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
	"sort"
	"strconv"
	"strings"
	"time"
)

const tolerance = 1e-6

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refVal, err := parseFloat(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVal, err := parseFloat(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if !almostEqual(refVal, gotVal, tolerance) {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %.12f got %.12f\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
				idx+1, refVal, gotVal, tc.input, refOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "1510F_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "1510F.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseFloat(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single number, got %d tokens", len(fields))
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", fields[0], err)
	}
	return val, nil
}

func almostEqual(a, b, tol float64) bool {
	diff := math.Abs(a - b)
	if diff <= tol {
		return true
	}
	div := math.Max(1.0, math.Abs(a))
	return diff <= tol*div
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 80)...)
	tests = append(tests, stressTests(rng)...)
	return tests
}

func manualTests() []testCase {
	polygons := []struct {
		points [][2]float64
		l      float64
	}{
		{
			points: [][2]float64{{0, 0}, {1, 0}, {0, 1}},
			l:      4,
		},
		{
			points: [][2]float64{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			l:      5,
		},
		{
			points: [][2]float64{{0, 0}, {2, -1}, {3, 0}, {4, 3}, {-1, 4}},
			l:      17,
		},
	}
	var tests []testCase
	for _, poly := range polygons {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %.6f\n", len(poly.points), poly.l))
		for _, pt := range poly.points {
			sb.WriteString(fmt.Sprintf("%.6f %.6f\n", pt[0], pt[1]))
		}
		tests = append(tests, testCase{input: sb.String()})
	}
	return tests
}

func randomTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	for i := 0; i < batches; i++ {
		n := rng.Intn(50) + 3
		radius := rng.Float64()*90 + 10
		points := generateConvexPolygon(rng, n, radius)
		perimeter := polygonPerimeter(points)
		l := perimeter + rng.Float64()*500 + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %.10f\n", n, l))
		for _, pt := range points {
			sb.WriteString(fmt.Sprintf("%.10f %.10f\n", pt[0], pt[1]))
		}
		tests = append(tests, testCase{input: sb.String()})
	}
	return tests
}

func stressTests(rng *rand.Rand) []testCase {
	var tests []testCase
	for i := 0; i < 5; i++ {
		n := 10000
		points := generateConvexPolygon(rng, n, 100000)
		perimeter := polygonPerimeter(points)
		l := perimeter + rng.Float64()*8e5 + 1e-3
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %.10f\n", n, l))
		for _, pt := range points {
			sb.WriteString(fmt.Sprintf("%.10f %.10f\n", pt[0], pt[1]))
		}
		tests = append(tests, testCase{input: sb.String()})
	}
	return tests
}

func generateConvexPolygon(rng *rand.Rand, n int, maxRadius float64) [][2]float64 {
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	points := make([][2]float64, n)
	for i := 0; i < n; i++ {
		r := rng.Float64()*maxRadius + 1
		points[i] = [2]float64{r * math.Cos(angles[i]), r * math.Sin(angles[i])}
	}
	return convexHull(points)
}

func polygonPerimeter(points [][2]float64) float64 {
	per := 0.0
	for i := 0; i < len(points); i++ {
		j := (i + 1) % len(points)
		per += math.Hypot(points[i][0]-points[j][0], points[i][1]-points[j][1])
	}
	return per
}

func convexHull(points [][2]float64) [][2]float64 {
	sort.Slice(points, func(i, j int) bool {
		if points[i][0] == points[j][0] {
			return points[i][1] < points[j][1]
		}
		return points[i][0] < points[j][0]
	})
	n := len(points)
	if n <= 1 {
		return points
	}
	half := func(points [][2]float64) [][2]float64 {
		h := make([][2]float64, 0, n)
		for _, p := range points {
			for len(h) >= 2 && cross(h[len(h)-2], h[len(h)-1], p) <= 0 {
				h = h[:len(h)-1]
			}
			h = append(h, p)
		}
		return h
	}
	lower := half(points)
	upper := half(reversePoints(points))
	lower = lower[:len(lower)-1]
	upper = upper[:len(upper)-1]
	return append(lower, upper...)
}

func reversePoints(points [][2]float64) [][2]float64 {
	res := make([][2]float64, len(points))
	copy(res, points)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func cross(a, b, c [2]float64) float64 {
	return (b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])
}
