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
	"time"
)

const (
	randomTests = 200
	maxK        = 1_000_000_000
	maxT        = 1000
)

type testInput struct {
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1])
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

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candidate2071A_%d", time.Now().UnixNano()))
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
	src := filepath.Join(dir, "2071A.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2071A_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runTest(test testInput, candidate, oracle string) error {
	oracleOut, err := runBinary(oracle, test.input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	candOut, err := runBinary(candidate, test.input)
	if err != nil {
		return fmt.Errorf("contestant runtime error: %v", err)
	}

	expect, err := parseOutputs(oracleOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse oracle output: %v", err)
	}
	got, err := parseOutputs(candOut, test.t)
	if err != nil {
		return fmt.Errorf("failed to parse contestant output: %v", err)
	}

	for i := 0; i < test.t; i++ {
		if expect[i] != got[i] {
			word := "NO"
			if expect[i] {
				word = "YES"
			}
			return fmt.Errorf("case %d: expected %s", i+1, word)
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutputs(output string, expected int) ([]bool, error) {
	reader := strings.NewReader(output)
	res := make([]bool, 0, expected)
	for len(res) < expected {
		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			return nil, fmt.Errorf("need %d tokens, got %d (%v)", expected, len(res), err)
		}
		token = strings.TrimSpace(token)
		if token == "" {
			return nil, fmt.Errorf("empty token")
		}
		if strings.EqualFold(token, "YES") {
			res = append(res, true)
		} else if strings.EqualFold(token, "NO") {
			res = append(res, false)
		} else {
			return nil, fmt.Errorf("invalid token %q", token)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("expected %d tokens, output has extra data", expected)
	}
	return res, nil
}

func deterministicTests() []testInput {
	tests := []testInput{
		buildTestInput([]int64{1, 2, 3, 10, 1000000000}),
		buildTestInput(makeArithmetic(1, 3, 20)),
		buildTestInput(makeRepeated(42, 100)),
	}
	tests = append(tests, buildTestInput(makeSequential(1, maxT)))
	return tests
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(maxT) + 1
	values := make([]int64, t)
	for i := range values {
		values[i] = rng.Int63n(maxK) + 1
	}
	return buildTestInput(values)
}

func buildTestInput(values []int64) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(values))
	for _, v := range values {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	return testInput{
		input: sb.String(),
		t:     len(values),
	}
}

func makeArithmetic(start, step int64, count int) []int64 {
	res := make([]int64, count)
	for i := 0; i < count; i++ {
		val := start + int64(i)*step
		if val > maxK {
			val = start
		}
		res[i] = val
	}
	return res
}

func makeSequential(start int64, count int) []int64 {
	res := make([]int64, count)
	val := start
	for i := 0; i < count; i++ {
		if val > maxK {
			val = 1
		}
		res[i] = val
		val++
	}
	return res
}

func makeRepeated(value int64, count int) []int64 {
	res := make([]int64, count)
	if value < 1 || value > maxK {
		value = 1
	}
	for i := 0; i < count; i++ {
		res[i] = value
	}
	return res
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
