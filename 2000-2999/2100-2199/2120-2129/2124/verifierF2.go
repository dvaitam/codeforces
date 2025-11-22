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
	maxT        = 8
)

type testInput struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1], "candidate2124F2")
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
	src := filepath.Join(dir, "2124F2.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2124F2_%d", time.Now().UnixNano()))
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
			return fmt.Errorf("case %d: expected %s got %s", i+1, expect[i], got[i])
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

func parseOutput(out string, t int) ([]string, error) {
	reader := strings.NewReader(out)
	res := make([]string, 0, t)
	for len(res) < t {
		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			return nil, fmt.Errorf("need %d tokens, got %d (%v)", t, len(res), err)
		}
		res = append(res, token)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d tokens, output has extra data", t)
	}
	return res, nil
}

func deterministicTests() []testInput {
	// Hand-crafted fixed tests that touch common branches (empty constraints, forbidden values, partial forbids).
	custom := buildInput([]caseSpec{
		{n: 3, m: 0, cons: nil},
		{n: 3, m: 3, cons: [][2]int{{1, 1}, {2, 1}, {3, 1}}},
		{n: 3, m: 1, cons: [][2]int{{1, 2}}},
		{n: 4, m: 4, cons: [][2]int{{1, 1}, {2, 1}, {3, 1}, {4, 1}}},
		{n: 5, m: 5, cons: [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 1}}},
		{n: 6, m: 2, cons: [][2]int{{2, 3}, {4, 2}}},
	})

	// A single near-maximum length case to ensure big limits are handled.
	big := buildInput([]caseSpec{
		{n: 5000, m: 1, cons: [][2]int{{69, 420}}},
	})

	return []testInput{
		custom,
		big,
	}
}

type caseSpec struct {
	n, m int
	cons [][2]int
}

func buildInput(cases []caseSpec) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		fmt.Fprintf(&sb, "%d %d\n", c.n, c.m)
		for _, p := range c.cons {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
	}
	return testInput{input: sb.String(), t: len(cases)}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(maxT) + 1
	cases := make([]caseSpec, t)
	for i := 0; i < t; i++ {
		cases[i] = randomCase(rng, 60, 100)
	}
	return buildInput(cases)
}

func randomCase(rng *rand.Rand, maxN, maxM int) caseSpec {
	n := rng.Intn(maxN) + 1
	maxPossible := n * n
	mLimit := maxM
	if maxPossible < mLimit {
		mLimit = maxPossible
	}
	m := rng.Intn(mLimit + 1)

	seen := make(map[[2]int]struct{}, m)
	cons := make([][2]int, 0, m)
	for len(cons) < m {
		i := rng.Intn(n) + 1
		x := rng.Intn(n) + 1
		key := [2]int{i, x}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		cons = append(cons, key)
	}
	return caseSpec{n: n, m: m, cons: cons}
}

func largeTests() []testInput {
	// Near-constraint-size cases to stress performance.
	c1 := caseSpec{
		n:    400,
		m:    500,
		cons: gridConstraints(400, 500),
	}
	c2 := caseSpec{
		n:    800,
		m:    800,
		cons: gridConstraints(800, 800),
	}
	return []testInput{
		buildInput([]caseSpec{c1}),
		buildInput([]caseSpec{c2}),
	}
}

func gridConstraints(n, m int) [][2]int {
	cons := make([][2]int, 0, m)
	for i := 1; i <= n && len(cons) < m; i++ {
		for x := 1; x <= n && len(cons) < m; x += 2 {
			cons = append(cons, [2]int{i, x})
		}
	}
	return cons
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
