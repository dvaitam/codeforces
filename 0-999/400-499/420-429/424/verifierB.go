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
	"strings"
	"time"
)

type point struct {
	x, y int
	k    int
}

type testCase struct {
	id    string
	input string
}

type expectedResult struct {
	impossible bool
	value      float64
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]
	if candidate == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}

	baseDir := currentDir()
	refBin, err := buildReference(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := evaluate(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		got, err := evaluate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if !compareResults(exp, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %s\nInput:\n%sExpected: %s\nGot: %s\n", tc.id, tc.input, formatResult(exp), formatResult(got))
			os.Exit(1)
		}
		if (i+1)%10 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d/%d tests...\n", i+1, len(tests))
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine caller file")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref424B.bin")
	cmd := exec.Command("go", "build", "-o", out, "424B.go")
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func evaluate(target, input string) (expectedResult, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return expectedResult{}, err
	}
	return parseOutput(out)
}

func parseOutput(out string) (expectedResult, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return expectedResult{}, fmt.Errorf("empty output")
	}
	reader := strings.NewReader(out)
	var token string
	if _, err := fmt.Fscan(reader, &token); err != nil {
		return expectedResult{}, fmt.Errorf("failed to read token: %v", err)
	}
	if token == "-1" {
		return expectedResult{impossible: true}, nil
	}
	var value float64
	if _, err := fmt.Sscan(token, &value); err != nil {
		return expectedResult{}, fmt.Errorf("failed to parse float: %v", err)
	}
	return expectedResult{value: value}, nil
}

func compareResults(exp, got expectedResult) bool {
	if exp.impossible || got.impossible {
		return exp.impossible == got.impossible
	}
	diff := math.Abs(exp.value - got.value)
	allowed := math.Max(1e-6, math.Abs(exp.value)*1e-6)
	return diff <= allowed
}

func formatResult(res expectedResult) string {
	if res.impossible {
		return "-1"
	}
	return fmt.Sprintf("%.10f", res.value)
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		id: "already-megacity",
		input: formatInput(1_000_000, []point{
			{x: 1, y: 0, k: 1},
		}),
	})
	tests = append(tests, testCase{
		id: "exact-hit",
		input: formatInput(999_000, []point{
			{x: 3, y: 4, k: 2_000},
		}),
	})
	tests = append(tests, testCase{
		id: "need-two-points",
		input: formatInput(500_000, []point{
			{x: 6, y: 8, k: 200_000},
			{x: 8, y: 15, k: 400_000},
		}),
	})
	tests = append(tests, testCase{
		id: "impossible",
		input: formatInput(100, []point{
			{x: 1000, y: 1000, k: 1},
			{x: -1000, y: -1000, k: 1},
		}),
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTestCase(rng, fmt.Sprintf("rand-%02d", i+1)))
	}
	return tests
}

func randomTestCase(rng *rand.Rand, id string) testCase {
	n := rng.Intn(20) + 1
	s := rng.Intn(900_000) + 50_000
	pts := make([]point, 0, n)
	used := make(map[[2]int]bool)
	for len(pts) < n {
		x := rng.Intn(20001) - 10000
		y := rng.Intn(20001) - 10000
		if x == 0 && y == 0 {
			continue
		}
		key := [2]int{x, y}
		if used[key] {
			continue
		}
		used[key] = true
		k := rng.Intn(999_999) + 1
		pts = append(pts, point{x: x, y: y, k: k})
	}
	// sometimes lower n with big k to ensure reachable
	if rng.Float64() < 0.3 {
		s = rng.Intn(200_000) + 1
	}
	return testCase{
		id:    id,
		input: formatInput(s, pts),
	}
}

func formatInput(s int, pts []point) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(pts), s)
	for _, p := range pts {
		fmt.Fprintf(&sb, "%d %d %d\n", p.x, p.y, p.k)
	}
	return sb.String()
}
