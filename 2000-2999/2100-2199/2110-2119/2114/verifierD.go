package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const referenceSource = "2000-2999/2100-2199/2110-2119/2114/2114D.go"

type point struct {
	x int64
	y int64
}

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	candidatePath, err := candidatePathFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	candidate, candCleanup, err := prepareBinary(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	// Build the provided reference solution to ensure it compiles, but
	// verification relies on independent checking to support any valid
	// solution.
	refPath := referencePath()
	if _, refCleanup, err := buildReferenceBinary(refPath); err == nil {
		defer refCleanup()
	} else {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}

	tests := buildTests()
	for idx, tc := range tests {
		cases, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing generated input %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		expect := make([]int64, len(cases))
		for i, c := range cases {
			expect[i] = solveCase(c)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\nraw output:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		for i := range expect {
			if candVals[i] != expect[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\n", idx+1, tc.name, i+1, expect[i], candVals[i])
				fmt.Fprintln(os.Stderr, previewInput(tc.input))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("usage: go run verifierD.go /path/to/binary-or-source")
	}
	return os.Args[1], nil
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), "2114D.go")
	}
	return referenceSource
}

func prepareBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "candidate2114D")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", bin, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "reference2114D")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

type inputCase struct {
	points []point
}

func parseInput(in string) ([]inputCase, error) {
	r := bufio.NewReader(strings.NewReader(in))
	var t int
	if _, err := fmt.Fscan(r, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	if t < 1 || t > 10000 {
		return nil, fmt.Errorf("t out of bounds: %d", t)
	}
	cases := make([]inputCase, t)
	sumN := 0
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(r, &n); err != nil {
			return nil, fmt.Errorf("failed to read n for case %d: %v", i+1, err)
		}
		if n < 1 || n > 200000 {
			return nil, fmt.Errorf("n out of bounds on case %d: %d", i+1, n)
		}
		sumN += n
		if sumN > 200000 {
			return nil, fmt.Errorf("sum of n exceeds limit")
		}
		pts := make([]point, n)
		for j := 0; j < n; j++ {
			var x, y int64
			if _, err := fmt.Fscan(r, &x, &y); err != nil {
				return nil, fmt.Errorf("failed to read point %d of case %d: %v", j+1, i+1, err)
			}
			if x < 1 || x > 1_000_000_000 || y < 1 || y > 1_000_000_000 {
				return nil, fmt.Errorf("point out of bounds at case %d idx %d", i+1, j+1)
			}
			pts[j] = point{x: x, y: y}
		}
		cases[i] = inputCase{points: pts}
	}
	return cases, nil
}

func solveCase(c inputCase) int64 {
	pts := c.points
	n := len(pts)
	if n == 1 {
		return 1
	}

	const inf int64 = 1 << 60
	minX1, minX2 := inf, inf
	maxX1, maxX2 := int64(-1), int64(-1)
	cntMinX, cntMaxX := 0, 0
	minY1, minY2 := inf, inf
	maxY1, maxY2 := int64(-1), int64(-1)
	cntMinY, cntMaxY := 0, 0

	for _, p := range pts {
		x := p.x
		switch {
		case x < minX1:
			minX2 = minX1
			minX1 = x
			cntMinX = 1
		case x == minX1:
			cntMinX++
		case x < minX2:
			minX2 = x
		}

		switch {
		case x > maxX1:
			maxX2 = maxX1
			maxX1 = x
			cntMaxX = 1
		case x == maxX1:
			cntMaxX++
		case x > maxX2:
			maxX2 = x
		}

		y := p.y
		switch {
		case y < minY1:
			minY2 = minY1
			minY1 = y
			cntMinY = 1
		case y == minY1:
			cntMinY++
		case y < minY2:
			minY2 = y
		}
		switch {
		case y > maxY1:
			maxY2 = maxY1
			maxY1 = y
			cntMaxY = 1
		case y == maxY1:
			cntMaxY++
		case y > maxY2:
			maxY2 = y
		}
	}

	widthAll := maxX1 - minX1 + 1
	heightAll := maxY1 - minY1 + 1
	best := widthAll * heightAll

	for _, p := range pts {
		if n-1 == 0 {
			best = minInt64(best, 1)
			continue
		}
		if n-1 == 1 {
			best = minInt64(best, 1)
			continue
		}
		minXO := minX1
		if p.x == minX1 && cntMinX == 1 {
			minXO = minX2
		}
		maxXO := maxX1
		if p.x == maxX1 && cntMaxX == 1 {
			maxXO = maxX2
		}
		minYO := minY1
		if p.y == minY1 && cntMinY == 1 {
			minYO = minY2
		}
		maxYO := maxY1
		if p.y == maxY1 && cntMaxY == 1 {
			maxYO = maxY2
		}
		width := maxXO - minXO + 1
		height := maxYO - minYO + 1
		area := width * height
		best = minInt64(best, area)
	}
	return best
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func buildTests() []testCase {
	var tests []testCase

	tests = append(tests, testCase{
		name:  "single-point",
		t:     2,
		input: "2\n1\n5 5\n2\n1 1\n1 2\n",
	})
	tests = append(tests, testCase{
		name:  "line-and-rectangle",
		t:     2,
		input: "2\n4\n1 1\n1 2\n1 3\n1 4\n5\n1 1\n4 4\n1 4\n4 1\n2 2\n",
	})
	tests = append(tests, testCase{
		name:  "two-points",
		t:     1,
		input: "1\n2\n100 1\n101 1\n",
	})
	tests = append(tests, testCase{
		name:  "duplicate-extremes",
		t:     1,
		input: "1\n5\n1 1\n1 3\n2 2\n3 1\n3 3\n",
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i, 10, 12))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, i, 60, 6))
	}

	return tests
}

func randomTest(rng *rand.Rand, idx, maxN, maxT int) testCase {
	t := rng.Intn(maxT-1) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			x := rng.Int63n(1_000_000_000) + 1
			y := rng.Int63n(1_000_000_000) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
	}
	return testCase{
		name:  fmt.Sprintf("rand-%d-%d", maxN, idx+1),
		t:     t,
		input: sb.String(),
	}
}

func previewInput(in string) string {
	const limit = 500
	if len(in) <= limit {
		return in
	}
	return in[:limit] + "..."
}
