package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const limitK = int64(1_000_000_000_000_000_000)

type pair struct {
	x int64
	y int64
}

type testInput struct {
	raw   string
	cases []pair
}

func buildReference() (string, error) {
	path := "./2085C_ref.bin"
	cmd := exec.Command("go", "build", "-o", path, "2085C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
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
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseInts(out string) ([]int64, error) {
	if strings.TrimSpace(out) == "" {
		return []int64{}, nil
	}
	fields := strings.Fields(out)
	res := make([]int64, len(fields))
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func buildInputFromPairs(pairs []pair) testInput {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(pairs)))
	sb.WriteByte('\n')
	for _, p := range pairs {
		sb.WriteString(strconv.FormatInt(p.x, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(p.y, 10))
		sb.WriteByte('\n')
	}
	copied := make([]pair, len(pairs))
	copy(copied, pairs)
	return testInput{raw: sb.String(), cases: copied}
}

func addManualTests() []testInput {
	tests := []testInput{}
	tests = append(tests, buildInputFromPairs([]pair{{x: 2, y: 5}}))
	tests = append(tests, buildInputFromPairs([]pair{{x: 6, y: 6}}))
	tests = append(tests, buildInputFromPairs([]pair{
		{x: 1, y: 1_000_000_000},
		{x: 1_000_000_000, y: 1},
		{x: 123456789, y: 987654321},
	}))
	return tests
}

func randomPairs(rng *rand.Rand, count int, maxVal int64) []pair {
	res := make([]pair, count)
	for i := 0; i < count; i++ {
		res[i] = pair{
			x: rng.Int63n(maxVal) + 1,
			y: rng.Int63n(maxVal) + 1,
		}
	}
	return res
}

func buildTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := addManualTests()

	for i := 0; i < 200; i++ {
		t := rng.Intn(4) + 1
		tests = append(tests, buildInputFromPairs(randomPairs(rng, t, 1000)))
	}

	for i := 0; i < 200; i++ {
		t := rng.Intn(5) + 1
		tests = append(tests, buildInputFromPairs(randomPairs(rng, t, 1_000_000_000)))
	}

	tests = append(tests, buildInputFromPairs(randomPairs(rng, 10000, 1_000_000_000)))

	return tests
}

func validateCase(x, y, cand, ref int64) error {
	if ref == -1 {
		if cand != -1 {
			return fmt.Errorf("expected -1 but got %d for x=%d y=%d", cand, x, y)
		}
		return nil
	}
	if cand == -1 {
		return fmt.Errorf("solution exists but got -1 for x=%d y=%d", x, y)
	}
	if cand < 0 || cand > limitK {
		return fmt.Errorf("k out of range (%d) for x=%d y=%d", cand, x, y)
	}
	sx := uint64(x) + uint64(cand)
	sy := uint64(y) + uint64(cand)
	if sx+sy != (sx ^ sy) {
		return fmt.Errorf("sum=%d xor=%d mismatch for x=%d y=%d k=%d", sx+sy, sx^sy, x, y, cand)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var refPath string
	fail := func(format string, args ...interface{}) {
		if refPath != "" {
			_ = os.Remove(refPath)
		}
		fmt.Fprintf(os.Stderr, format+"\n", args...)
		os.Exit(1)
	}

	var err error
	refPath, err = buildReference()
	if err != nil {
		fail("%v", err)
	}
	defer os.Remove(refPath)

	tests := buildTests()

	for idx, test := range tests {
		refOut, err := runProgram(refPath, test.raw)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		candOut, err := runProgram(bin, test.raw)
		if err != nil {
			fail("test %d: runtime error: %v\ninput:\n%s", idx+1, err, test.raw)
		}

		refVals, err := parseInts(refOut)
		if err != nil {
			fail("reference output parse error on test %d: %v", idx+1, err)
		}
		candVals, err := parseInts(candOut)
		if err != nil {
			fail("candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
		}

		if len(refVals) != len(test.cases) {
			fail("reference produced %d answers but expected %d on test %d", len(refVals), len(test.cases), idx+1)
		}
		if len(candVals) != len(test.cases) {
			fail("candidate produced %d answers but expected %d on test %d", len(candVals), len(test.cases), idx+1)
		}

		for i, tc := range test.cases {
			if err := validateCase(tc.x, tc.y, candVals[i], refVals[i]); err != nil {
				fail("test %d case %d failed: %v\ninput:\n%s", idx+1, i+1, err, test.raw)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
