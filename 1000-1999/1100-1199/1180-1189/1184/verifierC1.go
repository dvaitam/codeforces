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

const refSource = "./1184C1.go"

type point struct {
	x int
	y int
}

type testCase struct {
	name  string
	input string
	n     int
	pts   []point
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refPoint, err := parsePoint(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}
		if !isInteriorPoint(refPoint, tc.pts) {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s)\n", idx+1, tc.name)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candPoint, err := parsePoint(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if !isInteriorPoint(candPoint, tc.pts) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): point (%d %d) is not the unique interior point\ninput:\n%s\n", idx+1, tc.name, candPoint.x, candPoint.y, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1184C1-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
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

func parsePoint(out string) (point, error) {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return point{}, fmt.Errorf("expected two integers, got %d tokens", len(fields))
	}
	x, err := strconv.Atoi(fields[0])
	if err != nil {
		return point{}, fmt.Errorf("invalid x coordinate %q: %v", fields[0], err)
	}
	y, err := strconv.Atoi(fields[1])
	if err != nil {
		return point{}, fmt.Errorf("invalid y coordinate %q: %v", fields[1], err)
	}
	return point{x, y}, nil
}

func isInteriorPoint(target point, pts []point) bool {
	minX, maxX := 1<<31-1, -1<<31
	minY, maxY := 1<<31-1, -1<<31
	for _, p := range pts {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	return target.x > minX && target.x < maxX && target.y > minY && target.y < maxY
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("n2-simple", 2, 0, 4),
		buildCase("n3-middle", 3, 1, 7),
		buildCase("n5-large-square", 5, 2, 10),
		buildCase("n10-max", 10, 5, 20),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(9) + 2
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1), n))
	}
	return tests
}

func buildCase(name string, n int, offset int, side int) testCase {
	points := make([]point, 0, 4*n+1)
	minX := offset
	maxX := offset + side
	minY := offset
	maxY := offset + side
	for i := 0; i < n; i++ {
		points = append(points, point{minX, minY + i})
		points = append(points, point{maxX, minY + i})
		points = append(points, point{minX + i, minY})
		points = append(points, point{minX + i, maxY})
	}
	interior := point{(minX + maxX) / 2, (minY + maxY) / 2}
	points = append(points, interior)
	points = deduplicate(points)
	input := formatInput(n, points)
	return testCase{name: name, input: input, n: n, pts: points}
}

func randomCase(rng *rand.Rand, name string, n int) testCase {
	side := rng.Intn(41) + 5
	offsetX := rng.Intn(10)
	offsetY := rng.Intn(10)
	minX := offsetX
	maxX := offsetX + side
	minY := offsetY
	maxY := offsetY + side
	points := make([]point, 0, 4*n+1)
	genSide := func(x1, y1, x2, y2 int) {
		if x1 == x2 {
			for i := 0; i < n; i++ {
				y := y1 + rng.Intn(abs(y2-y1)+1)
				points = append(points, point{x1, y})
			}
		} else {
			for i := 0; i < n; i++ {
				x := x1 + rng.Intn(abs(x2-x1)+1)
				points = append(points, point{x, y1})
			}
		}
	}
	genSide(minX, minY, minX, maxY)
	genSide(maxX, minY, maxX, maxY)
	genSide(minX, minY, maxX, minY)
	genSide(minX, maxY, maxX, maxY)
	points = deduplicate(points)
	for len(points) < 4*n {
		points = append(points, point{minX, minY})
	}
	interior := point{rng.Intn(maxX-minX-1) + minX + 1, rng.Intn(maxY-minY-1) + minY + 1}
	points = append(points, interior)
	input := formatInput(n, points)
	return testCase{name: name, input: input, n: n, pts: points}
}

func deduplicate(points []point) []point {
	seen := make(map[point]struct{})
	out := make([]point, 0, len(points))
	for _, p := range points {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func formatInput(n int, points []point) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, p := range points {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
