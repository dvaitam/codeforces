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

type testCase struct {
	name  string
	input string
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to locate verifier path")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "reference build failed:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d (%s): %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d (%s): %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d (%s)\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, tc.name, exp, got, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, func(), error) {
	tmp, err := os.CreateTemp("", "cf-2053I1-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2053I1.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.Remove(tmp.Name())
	}
	return tmp.Name(), cleanup, nil
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

// ----------------- test generation -----------------

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(2053001))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, trivialTests())
	tests = append(tests, mixedSmallTest())
	tests = append(tests, randomPack("random-small", rng, 8, 30))
	tests = append(tests, randomPack("random-mid", rng, 6, 2000))
	tests = append(tests, randomPack("random-large", rng, 3, 70000))
	tests = append(tests, heavySingleCase(rng))

	return tests
}

func sampleTest() testCase {
	input := "4\n4\n1 2 3 4\n4\n2 -3 2 2\n10\n2 -7 6 3 -1 4 2 -5 8 -4\n20\n4 -2 4 3 -2 1 5 2 3 6 -5 -1 -4 -2 -3 5 -3 1 -4 1\n"
	return testCase{name: "sample", input: input}
}

func trivialTests() testCase {
	cases := [][]int64{
		{5},
		{0, 0, 0, 0},
		{3, -1, -1, -1},
		{10, -3, 2, 1},
		{-2, 5, -1, -2, 0},
	}
	return packCases("trivial", cases)
}

func mixedSmallTest() testCase {
	cases := [][]int64{
		{1, -1, 1, -1, 1},
		{4, -2, 4, -2, 4, -2},
		{10, -5, 3, -3, 2, -2, 1, -1, 1, -1},
	}
	return packCases("mixed-small", cases)
}

func randomPack(name string, rng *rand.Rand, t int, maxN int) testCase {
	cases := make([][]int64, 0, t)
	for i := 0; i < t; i++ {
		n := 1 + rng.Intn(maxN)
		cases = append(cases, randomArray(rng, n))
	}
	return packCases(name, cases)
}

func heavySingleCase(rng *rand.Rand) testCase {
	n := 150000 + rng.Intn(20000)
	cases := [][]int64{randomArray(rng, n)}
	return packCases("heavy-single", cases)
}

// randomArray builds an array with sum >= max(abs(ai)) to satisfy constraints.
func randomArray(rng *rand.Rand, n int) []int64 {
	arr := make([]int64, n)
	var sum int64
	var maxAbs int64

	for i := 0; i < n; i++ {
		val := int64(rng.Intn(2000000)) - 1_000_000
		arr[i] = val
		if val < 0 {
			// keep negatives modest to ease adjustments
			arr[i] = val / 4
		}
		if arr[i] == 0 && rng.Intn(4) == 0 {
			arr[i] = 1
		}
		if abs := abs64(arr[i]); abs > maxAbs {
			maxAbs = abs
		}
		sum += arr[i]
	}

	if sum < maxAbs {
		need := maxAbs - sum + int64(rng.Intn(10))
		if need > 900_000_000 {
			need = 900_000_000
		}
		arr[0] += need
		sum += need
		if abs := abs64(arr[0]); abs > maxAbs {
			maxAbs = abs
		}
	}
	if sum < 0 {
		add := -sum + int64(rng.Intn(10))
		if arr[0]+add > 1_000_000_000 {
			add = 1_000_000_000 - arr[0]
		}
		arr[0] += add
		sum += add
	}
	// Final guard to ensure constraint.
	if maxAbs > sum {
		diff := maxAbs - sum
		if diff > 0 {
			if arr[0]+diff > 1_000_000_000 {
				diff = 1_000_000_000 - arr[0]
			}
			arr[0] += diff
		}
	}
	return arr
}

func packCases(name string, cases [][]int64) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, arr := range cases {
		fmt.Fprintf(&b, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{name: name, input: b.String()}
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
