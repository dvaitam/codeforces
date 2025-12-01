package main

import (
	"bufio"
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

const refSource = "./2074C.go"

type testCase struct {
	x int64
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2074C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{x: 2}, // smallest, power of two
		{x: 3}, // all ones
		{x: 4}, // power of two
		{x: 5}, // should have answer
		{x: 6},
		{x: 7}, // all ones
		{x: 8}, // power of two
		{x: 9},
		{x: 10},
		{x: 11},
		{x: 12},
		{x: 13},
		{x: 15}, // all ones
		{x: 16}, // power of two
		{x: 31}, // all ones
		{x: 32}, // power of two
		{x: 33},
		{x: 1000000000}, // upper limit
		{x: 999999937},
		{x: 123456789},
	}
}

func randomTest(rng *rand.Rand) testCase {
	return testCase{x: rng.Int63n(1_000_000_000-1) + 2}
}

func randomPowerOfTwo(rng *rand.Rand) testCase {
	shift := rng.Intn(29) + 1 // keep within limit (2^29 <= 1e9)
	return testCase{x: int64(1) << shift}
}

func randomAllOnes(rng *rand.Rand) testCase {
	shift := rng.Intn(29) + 1 // (1<<29)-1 < 1e9
	return testCase{x: (int64(1) << shift) - 1}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.FormatInt(tc.x, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		val, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer on test %d: %v", i+1, err)
		}
		res[i] = val
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func isTriangle(x, y, z int64) bool {
	return x+y > z && x+z > y && y+z > x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomPowerOfTwo(rng))
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomAllOnes(rng))
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refHas := want[i] != -1
		y := got[i]
		if y == -1 {
			if refHas {
				fmt.Fprintf(os.Stderr, "test %d failed: solution exists but candidate output -1\nx=%d\n", i+1, tc.x)
				os.Exit(1)
			}
			continue
		}
		if y <= 0 || y >= tc.x {
			fmt.Fprintf(os.Stderr, "test %d failed: invalid y range %d for x=%d\n", i+1, y, tc.x)
			os.Exit(1)
		}
		z := tc.x ^ y
		if !isTriangle(tc.x, y, z) {
			fmt.Fprintf(os.Stderr, "test %d failed: sides do not form non-degenerate triangle x=%d y=%d z=%d\n", i+1, tc.x, y, z)
			os.Exit(1)
		}
		if !refHas {
			fmt.Fprintf(os.Stderr, "test %d failed: reference says no solution but candidate produced valid-looking y\nx=%d y=%d\n", i+1, tc.x, y)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
