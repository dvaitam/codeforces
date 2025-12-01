package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2169E.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2169E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
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

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(21692169))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeCase([]point{
		{0, 0, 0},
		{0, 5, 10},
		{5, 0, 10},
	}))
	tests = append(tests, makeCase([]point{
		{0, 0, 100},
		{1, 0, 1},
		{0, 1, 1},
		{1, 1, 1},
	}))

	for i := 0; i < 25; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(8)+1, 30))
	}

	tests = append(tests, randomCase(rng, 80, 2000))
	tests = append(tests, structuredLineCase(10000, true))
	tests = append(tests, structuredLineCase(10000, false))
	tests = append(tests, maxCase())

	return tests
}

type point struct {
	x int64
	y int64
	c int64
}

func makeCase(points []point) testCase {
	return singleCase(points)
}

func sampleTest() testCase {
	return testCase{
		input: "4\n" +
			"1\n1\n1\n42\n" +
			"4\n1 2 10 5\n0 5 10 5\n1 1 1 1\n" +
			"3\n6 7 8\n3 3 3\n9 9 9\n" +
			"2\n1000000000000000000 10\n10 1000000000000000000\n12345 54321\n",
	}
}

func singleCase(points []point) testCase {
	return multiCase([][]point{points})
}

func multiCase(instances [][]point) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(instances))
	for _, pts := range instances {
		fmt.Fprintln(&b, len(pts))
		writePoints(&b, pts, func(p point) int64 { return p.x })
		writePoints(&b, pts, func(p point) int64 { return p.y })
		writePoints(&b, pts, func(p point) int64 { return p.c })
	}
	return testCase{input: b.String()}
}

func writePoints(b *strings.Builder, pts []point, f func(point) int64) {
	for i, p := range pts {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", f(p))
	}
	b.WriteByte('\n')
}

func randomCase(rng *rand.Rand, tCount int, maxTotal int) testCase {
	if tCount < 1 {
		tCount = 1
	}
	var instances [][]point
	remaining := maxTotal
	for i := 0; i < tCount && remaining > 0; i++ {
		n := rng.Intn(min(remaining, 5000)) + 1
		instances = append(instances, uniquePoints(rng, n))
		remaining -= n
	}
	if len(instances) == 0 {
		instances = append(instances, uniquePoints(rng, 1))
	}
	return multiCase(instances)
}

func uniquePoints(rng *rand.Rand, n int) []point {
	seen := make(map[[2]int64]struct{}, n)
	points := make([]point, 0, n)
	for len(points) < n {
		x := rng.Int63n(1_000_000_000_000_000)
		y := rng.Int63n(1_000_000_000_000_000)
		key := [2]int64{x, y}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		c := rng.Int63n(1_000_000_000) + 1
		points = append(points, point{x: x, y: y, c: c})
	}
	return points
}

func structuredLineCase(n int, vertical bool) testCase {
	points := make([]point, n)
	for i := 0; i < n; i++ {
		if vertical {
			points[i] = point{x: 0, y: int64(i), c: int64(i % 5)}
		} else {
			points[i] = point{x: int64(i), y: 0, c: int64(i%7 + 1)}
		}
	}
	return singleCase(points)
}

func maxCase() testCase {
	n := 300000
	points := make([]point, n)
	for i := 0; i < n; i++ {
		points[i] = point{
			x: int64(i),
			y: int64(i * 2),
			c: int64((i % 1000) + 1),
		}
	}
	return singleCase(points)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
