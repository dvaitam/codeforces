package main

import (
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

const refSource = "0-999/900-999/950-959/958/958C3.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC3.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expVal, err := parseScalar(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVal, err := parseScalar(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}

		if gotVal != expVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, expVal, gotVal, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "958C3-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseScalar(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("sample-like-1", []int{3, 4, 7, 2}, 3, 10),
		buildCase("sample-like-2", []int{16, 3, 24, 13, 9, 8, 7, 5, 12, 12}, 5, 15),
		buildCase("small-edge", []int{1, 1}, 2, 2),
		buildCase("all-same", repeatValue(5, 20), 10, 7),
		buildCase("max-k", repeatSequence([]int{1, 2, 3, 4, 5}, 10), 50, 13),
		buildCase("descending", []int{50, 40, 30, 20, 10, 5, 3}, 3, 11),
	}

	tests = append(tests, stressCase())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1)))
	}
	return tests
}

func buildCase(name string, arr []int, k, p int) testCase {
	n := len(arr)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, k, p)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func repeatValue(val, n int) []int {
	res := make([]int, n)
	for i := range res {
		res[i] = val
	}
	return res
}

func repeatSequence(seq []int, times int) []int {
	res := make([]int, 0, len(seq)*times)
	for i := 0; i < times; i++ {
		res = append(res, seq...)
	}
	return res
}

func stressCase() testCase {
	n := 1000
	k := 50
	p := 37
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = (i*17)%1000000 + 1
	}
	return buildCase("stress", arr, k, p)
}

func randomCase(rng *rand.Rand, name string) testCase {
	n := rng.Intn(500) + 1
	k := rng.Intn(min(n, 50)) + 1
	p := rng.Intn(99) + 2
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1_000_000) + 1
	}
	return buildCase(name, arr, k, p)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
