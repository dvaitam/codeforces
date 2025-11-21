package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	refSourceB3  = "0-999/600-699/690-699/690/690B3.go"
	randomInputs = 80
)

type testCase struct {
	N      int
	points [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB3.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceB3)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicInputs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomInputs; i++ {
		tests = append(tests, randomInput(rng))
	}

	for idx, input := range tests {
		expect, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fail("candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input)
		}
		if normalize(expect) != normalize(got) {
			fail("test %d mismatch\ninput:\n%s\nexpected: %s\ngot: %s", idx+1, input, expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "690B3-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
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

func deterministicInputs() []string {
	var cases []string
	base := testCase{
		N: 7,
		points: [][2]int{
			{1, 1}, {1, 4}, {2, 2}, {3, 3},
			{4, 5}, {5, 2}, {6, 6}, {7, 3},
		},
	}
	cases = append(cases, buildInput([]testCase{base}))

	tc2 := []testCase{
		{N: 20, points: sequentialPoints(20, 12)},
		{N: 50, points: sequentialPoints(50, 30)},
	}
	cases = append(cases, buildInput(tc2))

	tc3 := testCase{N: 100000, points: sequentialPoints(100000, 1500)}
	cases = append(cases, buildInput([]testCase{tc3}))

	seeded := rand.New(rand.NewSource(1337))
	tc4 := []testCase{
		structuredRandomCase(seeded, 500, 200),
		structuredRandomCase(seeded, 2000, 500),
	}
	cases = append(cases, buildInput(tc4))
	return cases
}

func randomInput(rng *rand.Rand) string {
	testCount := rng.Intn(3) + 1
	var tcs []testCase
	totalM := 0
	for i := 0; i < testCount; i++ {
		remain := 200000 - totalM
		if remain < 8 {
			break
		}
		n := randomN(rng)
		maxCap := remain
		caseCap := 2000
		if rng.Intn(5) == 0 {
			caseCap = 6000
		}
		if maxCap > caseCap {
			maxCap = caseCap
		}
		maxPossible := n * n
		if maxCap > maxPossible {
			maxCap = maxPossible
		}
		if maxCap < 8 {
			continue
		}
		m := randRange(rng, 8, maxCap)
		points := randomPoints(rng, n, m)
		tcs = append(tcs, testCase{N: n, points: points})
		totalM += m
	}
	if len(tcs) == 0 {
		tcs = append(tcs, testCase{N: 5, points: sequentialPoints(5, 8)})
	}
	return buildInput(tcs)
}

func structuredRandomCase(rng *rand.Rand, n, count int) testCase {
	if count < 8 {
		count = 8
	}
	if count > n*n {
		count = n * n
	}
	points := randomPoints(rng, n, count)
	return testCase{N: n, points: points}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	totalM := 0
	for _, tc := range cases {
		m := len(tc.points)
		totalM += m
		fmt.Fprintf(&sb, "%d %d\n", tc.N, m)
		for _, p := range tc.points {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
	}
	sb.WriteString("0 0\n")
	return sb.String()
}

func sequentialPoints(n, count int) [][2]int {
	if count > n*n {
		count = n * n
	}
	points := make([][2]int, 0, count)
	x, y := 1, 1
	for len(points) < count && y <= n {
		points = append(points, [2]int{x, y})
		x++
		if x > n {
			x = 1
			y++
		}
	}
	return points
}

func randomPoints(rng *rand.Rand, n, count int) [][2]int {
	points := make([][2]int, 0, count)
	used := make(map[int64]struct{}, count*2)
	for len(points) < count {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		key := (int64(x) << 32) | int64(y)
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		points = append(points, [2]int{x, y})
	}
	sort.Slice(points, func(i, j int) bool {
		if points[i][0] == points[j][0] {
			return points[i][1] < points[j][1]
		}
		return points[i][0] < points[j][0]
	})
	return points
}

func randomN(rng *rand.Rand) int {
	switch rng.Intn(5) {
	case 0:
		return 5 + rng.Intn(15) // 5..19
	case 1:
		return 20 + rng.Intn(80) // 20..99
	case 2:
		return 100 + rng.Intn(900) // 100..999
	case 3:
		return 1000 + rng.Intn(9000) // 1000..9999
	default:
		return 10000 + rng.Intn(90001) // 10000..100000
	}
}

func randRange(rng *rand.Rand, lo, hi int) int {
	if hi < lo {
		lo, hi = hi, lo
	}
	if lo == hi {
		return lo
	}
	return lo + rng.Intn(hi-lo+1)
}

func normalize(out string) string {
	return strings.Join(strings.Fields(out), " ")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
