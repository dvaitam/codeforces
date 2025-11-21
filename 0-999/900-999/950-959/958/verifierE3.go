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

const refSource = "0-999/900-999/950-959/958/958E3.go"

type point struct {
	X int64
	Y int64
}

type instance struct {
	n     int
	ships []point
	bases []point
}

type testCase struct {
	name  string
	input string
	inst  instance
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE3.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := validateSolution(tc.inst, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		out, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if err := validateSolution(tc.inst, out); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "958E3-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func validateSolution(inst instance, output string) error {
	permutation, err := parsePermutation(output, inst.n)
	if err != nil {
		return err
	}
	for i := 0; i < inst.n; i++ {
		if permutation[i] < 0 || permutation[i] >= inst.n {
			return fmt.Errorf("base index out of range at line %d", i+1)
		}
	}
	for i := 0; i < inst.n; i++ {
		for j := i + 1; j < inst.n; j++ {
			if segmentsIntersect(inst.ships[i], inst.bases[permutation[i]], inst.ships[j], inst.bases[permutation[j]]) {
				return fmt.Errorf("segments %d and %d intersect", i+1, j+1)
			}
		}
	}
	return nil
}

func parsePermutation(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	used := make([]bool, n)
	res := make([]int, n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d: %v", f, i+1, err)
		}
		if val < 1 || val > n {
			return nil, fmt.Errorf("value %d at position %d out of range [1,%d]", val, i+1, n)
		}
		if used[val-1] {
			return nil, fmt.Errorf("base %d is assigned multiple times", val)
		}
		used[val-1] = true
		res[i] = val - 1
	}
	for idx, ok := range used {
		if !ok {
			return nil, fmt.Errorf("base %d is unused", idx+1)
		}
	}
	return res, nil
}

func segmentsIntersect(p1, q1, p2, q2 point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 {
		return true
	}
	if o1 == 0 && onSegment(p1, q1, p2) {
		return true
	}
	if o2 == 0 && onSegment(p1, q1, q2) {
		return true
	}
	if o3 == 0 && onSegment(p2, q2, p1) {
		return true
	}
	if o4 == 0 && onSegment(p2, q2, q1) {
		return true
	}
	return false
}

func orientation(a, b, c point) int {
	val := (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
	switch {
	case val > 0:
		return 1
	case val < 0:
		return -1
	default:
		return 0
	}
}

func onSegment(a, b, c point) bool {
	return c.X >= min64(a.X, b.X) && c.X <= max64(a.X, b.X) &&
		c.Y >= min64(a.Y, b.Y) && c.Y <= max64(a.Y, b.Y)
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("n1-basic",
			[]point{{0, 0}},
			[]point{{1, 1}},
		),
		buildCase("n2-simple",
			[]point{{0, 0}, {2, 0}},
			[]point{{0, 2}, {2, 2}},
		),
		buildCase("n3-triangle",
			[]point{{-2, 0}, {0, 3}, {3, 1}},
			[]point{{-3, 4}, {1, -2}, {4, 5}},
		),
		buildCase("zigzag",
			[]point{{-4, -1}, {-1, 4}, {2, -2}, {5, 3}},
			[]point{{-5, 3}, {-2, -4}, {3, 5}, {6, -1}},
		),
	}

	rng := rand.New(rand.NewSource(958e3001))
	tests = append(tests, randomCase(rng, "random-5", 5))
	tests = append(tests, randomCase(rng, "random-8", 8))
	tests = append(tests, randomCase(rng, "random-12", 12))
	for i := 0; i < 30; i++ {
		n := rng.Intn(40) + 2
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1), n))
	}
	return tests
}

func buildCase(name string, ships, bases []point) testCase {
	if len(ships) != len(bases) {
		panic("ships and bases length mismatch")
	}
	total := append(clonePoints(ships), clonePoints(bases)...)
	if err := checkGeneralPosition(total); err != nil {
		panic(fmt.Sprintf("invalid points for case %s: %v", name, err))
	}
	var sb strings.Builder
	sb.Grow((len(ships)*2 + 1) * 24)
	fmt.Fprintf(&sb, "%d\n", len(ships))
	for _, pt := range ships {
		fmt.Fprintf(&sb, "%d %d\n", pt.X, pt.Y)
	}
	for _, pt := range bases {
		fmt.Fprintf(&sb, "%d %d\n", pt.X, pt.Y)
	}
	inst := instance{
		n:     len(ships),
		ships: clonePoints(ships),
		bases: clonePoints(bases),
	}
	return testCase{name: name, input: sb.String(), inst: inst}
}

func clonePoints(src []point) []point {
	dst := make([]point, len(src))
	copy(dst, src)
	return dst
}

func checkGeneralPosition(points []point) error {
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			if points[i] == points[j] {
				return fmt.Errorf("duplicate point at indices %d and %d", i, j)
			}
			for k := j + 1; k < len(points); k++ {
				if orientation(points[i], points[j], points[k]) == 0 {
					return fmt.Errorf("points %d, %d, %d are colinear", i, j, k)
				}
			}
		}
	}
	return nil
}

func randomCase(rng *rand.Rand, name string, n int) testCase {
	total := 2 * n
	points := make([]point, 0, total)
	for len(points) < total {
		x := rng.Intn(20001) - 10000
		y := rng.Intn(20001) - 10000
		candidate := point{X: int64(x), Y: int64(y)}
		if duplicatePoint(candidate, points) {
			continue
		}
		if formsColinear(candidate, points) {
			continue
		}
		points = append(points, candidate)
	}
	ships := clonePoints(points[:n])
	bases := clonePoints(points[n:])
	return buildCase(name, ships, bases)
}

func duplicatePoint(pt point, points []point) bool {
	for _, p := range points {
		if p == pt {
			return true
		}
	}
	return false
}

func formsColinear(pt point, points []point) bool {
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			if orientation(points[i], points[j], pt) == 0 {
				return true
			}
		}
	}
	return false
}
