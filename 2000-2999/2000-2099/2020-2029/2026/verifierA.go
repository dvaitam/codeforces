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

type testCase struct {
	X int
	Y int
	K int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s", err, candOut)
		os.Exit(1)
	}

	refSegments, err := parseSegments(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n%s", err, refOut)
		os.Exit(1)
	}
	if err := validateAll(tests, refSegments); err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\n", err)
		os.Exit(1)
	}

	candSegments, err := parseSegments(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\n%s", err, candOut)
		os.Exit(1)
	}
	if err := validateAll(tests, candSegments); err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2026A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2026A.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(X, Y, K int) {
		tests = append(tests, testCase{X: X, Y: Y, K: K})
	}
	add(1, 1, 1)
	add(3, 4, 1)
	add(4, 3, 3)
	add(3, 4, 3)
	add(1000, 1000, 1414)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		X := rng.Intn(1000) + 1
		Y := rng.Intn(1000) + 1
		if X > Y {
			X, Y = Y, X
		}
		if Y > X {
			Y = X
		}
		K := int(math.Sqrt(float64(X*X + Y*Y)))
		if K == 0 {
			continue
		}
		tests = append(tests, testCase{X: X, Y: Y, K: K})
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.X, tc.Y, tc.K))
	}
	return sb.String()
}

type segment struct {
	ax, ay int
	bx, by int
	cx, cy int
	dx, dy int
}

func parseSegments(out string, t int) ([]segment, error) {
	fields := strings.Fields(out)
	if len(fields) != t*8 {
		return nil, fmt.Errorf("expected %d numbers, got %d", t*8, len(fields))
	}
	segs := make([]segment, t)
	for i := 0; i < t; i++ {
		offset := i * 8
		vals := make([]int, 8)
		for j := 0; j < 8; j++ {
			num, err := strconv.Atoi(fields[offset+j])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q at position %d", fields[offset+j], offset+j+1)
			}
			vals[j] = num
		}
		segs[i] = segment{
			ax: vals[0], ay: vals[1],
			bx: vals[2], by: vals[3],
			cx: vals[4], cy: vals[5],
			dx: vals[6], dy: vals[7],
		}
	}
	return segs, nil
}

func validateAll(tests []testCase, segs []segment) error {
	if len(tests) != len(segs) {
		return fmt.Errorf("number of tests and segments mismatch")
	}
	for i := range tests {
		if err := validateSegment(tests[i], segs[i]); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	return nil
}

func validateSegment(tc testCase, seg segment) error {
	if !inBounds(seg.ax, seg.ay, tc.X, tc.Y) ||
		!inBounds(seg.bx, seg.by, tc.X, tc.Y) ||
		!inBounds(seg.cx, seg.cy, tc.X, tc.Y) ||
		!inBounds(seg.dx, seg.dy, tc.X, tc.Y) {
		return fmt.Errorf("points out of bounds")
	}
	if lengthSquared(seg.ax, seg.ay, seg.bx, seg.by) < tc.K*tc.K {
		return fmt.Errorf("segment AB shorter than required")
	}
	if lengthSquared(seg.cx, seg.cy, seg.dx, seg.dy) < tc.K*tc.K {
		return fmt.Errorf("segment CD shorter than required")
	}
	v1x := seg.bx - seg.ax
	v1y := seg.by - seg.ay
	v2x := seg.dx - seg.cx
	v2y := seg.dy - seg.cy
	if v1x == 0 && v1y == 0 {
		return fmt.Errorf("segment AB has zero length")
	}
	if v2x == 0 && v2y == 0 {
		return fmt.Errorf("segment CD has zero length")
	}
	if v1x*v2x+v1y*v2y != 0 {
		return fmt.Errorf("segments not perpendicular")
	}
	return nil
}

func inBounds(x, y, X, Y int) bool {
	return 0 <= x && x <= X && 0 <= y && y <= Y
}

func lengthSquared(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}
