package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	randomTests = 200
	maxT        = 20
)

type testInput struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1], "candidate2097C")
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candCleanup != nil {
		defer candCleanup()
	}

	oracle, oracleCleanup, err := prepareOracle()
	if err != nil {
		fmt.Println("failed to prepare reference solution:", err)
		return
	}
	defer oracleCleanup()

	tests := deterministicTests()
	total := 0
	for idx, test := range tests {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("deterministic test %d failed: %v\ninput:\n%s", idx+1, err, test.input)
			return
		}
		total++
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTests; i++ {
		test := randomTest(rng)
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("random test %d failed: %v\ninput:\n%s", i+1, err, test.input)
			return
		}
		total++
	}

	for idx, test := range largeTests() {
		if err := runTest(test, candidate, oracle); err != nil {
			fmt.Printf("large test %d failed: %v\ninput length: %d bytes\n", idx+1, err, len(test.input))
			return
		}
		total++
	}

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path, prefix string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle() (string, func(), error) {
	dir := sourceDir()
	src := filepath.Join(dir, "2097C.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2097C_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runTest(test testInput, candidate, oracle string) error {
	expectOut, err := runBinary(oracle, test.input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	gotOut, err := runBinary(candidate, test.input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}

	expect, err := parseOutput(expectOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse oracle output: %v", err)
	}
	got, err := parseOutput(gotOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse contestant output: %v", err)
	}

	for i := 0; i < test.t; i++ {
		if expect[i] != got[i] {
			return fmt.Errorf("case %d: expected %d got %d", i+1, expect[i], got[i])
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func parseOutput(out string, t int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, 0, t)
	for len(res) < t {
		var x int64
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return nil, fmt.Errorf("need %d integers, got %d (%v)", t, len(res), err)
		}
		res = append(res, x)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d integers, output has extra data", t)
	}
	return res, nil
}

func deterministicTests() []testInput {
	// Cases crafted to hit divisibility / non-divisibility branches.
	return []testInput{
		buildInput([]caseSpec{
			{n: 6, x: 2, y: 2, vx: 5, vy: 2},
			{n: 6, x: 2, y: 2, vx: 20, vy: 8},
			{n: 4, x: 1, y: 2, vx: 1, vy: 1},
			{n: 4, x: 1, y: 1, vx: 1, vy: 2},
			{n: 4, x: 1, y: 1, vx: 2, vy: 1},
			{n: 6, x: 2, y: 3, vx: 2, vy: 3},
		}),
		buildInput([]caseSpec{
			{n: 1_000_000_000, x: 123456789, y: 23456789, vx: 999_999_937, vy: 17},
			{n: 1_000_000_000, x: 1, y: 1, vx: 1_000_000_000, vy: 1_000_000_000},
			{n: 1_000_000_000, x: 400_000_000, y: 300_000_000, vx: 700_000_000, vy: 900_000_000},
		}),
	}
}

type caseSpec struct {
	n, x, y, vx, vy int64
}

func buildInput(cases []caseSpec) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", c.n, c.x, c.y, c.vx, c.vy)
	}
	return testInput{input: sb.String(), t: len(cases)}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(maxT) + 1
	cases := make([]caseSpec, t)
	for i := 0; i < t; i++ {
		cases[i] = randomCase(rng)
	}
	return buildInput(cases)
}

func randomCase(rng *rand.Rand) caseSpec {
	n := int64(rng.Intn(1_000_000_000-3) + 3)
	// Ensure x + y < n
	x := int64(rng.Int63n(n-1) + 1)
	yLimit := n - x - 1
	if yLimit < 1 {
		// force smaller x
		x = 1
		yLimit = n - x - 1
	}
	y := int64(rng.Int63n(yLimit) + 1)
	vx := int64(rng.Int63n(1_000_000_000) + 1)
	vy := int64(rng.Int63n(1_000_000_000) + 1)
	return caseSpec{n: n, x: x, y: y, vx: vx, vy: vy}
}

func largeTests() []testInput {
	// Stress with maximum constraints and varied directions.
	return []testInput{
		buildInput([]caseSpec{
			{n: 1_000_000_000, x: 500_000_000, y: 499_999_999, vx: 1_000_000_000, vy: 1},
			{n: 1_000_000_000, x: 1, y: 1, vx: 999_999_999, vy: 999_999_999},
		}),
		buildInput([]caseSpec{
			{n: 1_000_000_000, x: 123456789, y: 87654321, vx: 987654321, vy: 123456789},
			{n: 1_000_000_000, x: 999_999_9, y: 42, vx: 1_000_000_000, vy: 500_000_000},
			{n: 1_000_000_000, x: 42, y: 999_999_9, vx: 500_000_000, vy: 1_000_000_000},
		}),
	}
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
