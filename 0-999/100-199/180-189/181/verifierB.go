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
)

const coordLimit = 1000

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, input := range tests {
		expected, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %q, got %q\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-181B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "181B.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func generateTests() []string {
	var tests []string

	// Minimal configuration with one valid triple
	tests = append(tests, pointsToInput([][2]int{
		{0, 0}, {2, 0}, {1, 0},
	}))

	// No valid triples
	tests = append(tests, pointsToInput([][2]int{
		{0, 0}, {1, 0}, {2, 1}, {3, 1},
	}))

	// Dense line of evenly spaced points
	tests = append(tests, pointsToInput(linePoints(-50, -50, 2, 0, 51)))

	// Small grid with many combinations
	tests = append(tests, pointsToInput(gridPoints(-4, 4, 2, -4, 4, 2)))

	// Maximum size structured case
	tests = append(tests, pointsToInput(maxCase(3000)))

	rng := rand.New(rand.NewSource(2024))
	for i := 0; i < 25; i++ {
		size := 3 + rng.Intn(70)
		tests = append(tests, pointsToInput(randomWithMidpoints(rng, size)))
	}
	for i := 0; i < 8; i++ {
		size := 200 + rng.Intn(2801) // up to 3000
		tests = append(tests, pointsToInput(randomUnique(rng, size)))
	}
	return tests
}

func linePoints(x0, y0, dx, dy, count int) [][2]int {
	pts := make([][2]int, count)
	for i := 0; i < count; i++ {
		pts[i] = [2]int{x0 + i*dx, y0 + i*dy}
	}
	return pts
}

func gridPoints(xStart, xEnd, xStep, yStart, yEnd, yStep int) [][2]int {
	var pts [][2]int
	for x := xStart; x <= xEnd; x += xStep {
		for y := yStart; y <= yEnd; y += yStep {
			pts = append(pts, [2]int{x, y})
		}
	}
	return pts
}

func maxCase(n int) [][2]int {
	pts := make([][2]int, n)
	for i := 0; i < n; i++ {
		x := -coordLimit + (i % (2*coordLimit + 1))
		y := -coordLimit + (i / (2*coordLimit + 1))
		if y > coordLimit {
			y = coordLimit
		}
		pts[i] = [2]int{x, y}
	}
	return pts
}

func randomUnique(rng *rand.Rand, n int) [][2]int {
	pts := make([][2]int, 0, n)
	seen := make(map[[2]int]struct{}, n)
	for len(pts) < n {
		p := [2]int{
			rng.Intn(2*coordLimit+1) - coordLimit,
			rng.Intn(2*coordLimit+1) - coordLimit,
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		pts = append(pts, p)
	}
	return pts
}

func randomWithMidpoints(rng *rand.Rand, n int) [][2]int {
	pts := make([][2]int, 0, n)
	seen := make(map[[2]int]struct{})

	for len(pts) < n {
		if n-len(pts) >= 3 && rng.Intn(3) == 0 && addTriplet(rng, &pts, seen) {
			continue
		}
		p := [2]int{
			rng.Intn(2*coordLimit+1) - coordLimit,
			rng.Intn(2*coordLimit+1) - coordLimit,
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		pts = append(pts, p)
	}
	return pts
}

func addTriplet(rng *rand.Rand, pts *[][2]int, seen map[[2]int]struct{}) bool {
	for tries := 0; tries < 40; tries++ {
		mid := [2]int{
			rng.Intn(1601) - 800,
			rng.Intn(1601) - 800,
		}
		delta := [2]int{
			rng.Intn(401) - 200,
			rng.Intn(401) - 200,
		}
		if delta[0] == 0 && delta[1] == 0 {
			continue
		}
		p := [2]int{mid[0] - delta[0], mid[1] - delta[1]}
		q := [2]int{mid[0] + delta[0], mid[1] + delta[1]}
		if !inBounds(mid) || !inBounds(p) || !inBounds(q) {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		if _, ok := seen[q]; ok {
			continue
		}
		if _, ok := seen[mid]; ok {
			continue
		}
		seen[p] = struct{}{}
		seen[q] = struct{}{}
		seen[mid] = struct{}{}
		*pts = append(*pts, p, mid, q)
		return true
	}
	return false
}

func inBounds(p [2]int) bool {
	return p[0] >= -coordLimit && p[0] <= coordLimit && p[1] >= -coordLimit && p[1] <= coordLimit
}

func pointsToInput(points [][2]int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(points))
	for _, pt := range points {
		fmt.Fprintf(&b, "%d %d\n", pt[0], pt[1])
	}
	return b.String()
}
