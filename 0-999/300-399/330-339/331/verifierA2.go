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

const refSource = "0-999/300-399/330-339/331/331A2.go"

type testCase struct {
	name  string
	input string
	arr   []int64
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refSum, err := parseSum(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		sumReported, removedIdx, err := parseCandidateOutput(candOut, len(tc.arr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid output: %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if err := validateSolution(tc.arr, sumReported, removedIdx); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): solution invalid: %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if sumReported != refSum {
			fmt.Fprintf(os.Stderr, "test %d (%s): reported sum %d, expected %d\ninput:\n%soutput:\n%s", idx+1, tc.name, sumReported, refSum, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier331A2-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseSum(output string) (int64, error) {
	fields := strings.Fields(output)
	if len(fields) < 1 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid sum %q", fields[0])
	}
	return val, nil
}

func parseCandidateOutput(output string, n int) (int64, []int, error) {
	fields := strings.Fields(output)
	if len(fields) < 2 {
		return 0, nil, fmt.Errorf("expected at least two integers")
	}
	sumVal, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid sum %q", fields[0])
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k %q", fields[1])
	}
	if k < 0 || k > n {
		return 0, nil, fmt.Errorf("invalid k value %d", k)
	}
	if len(fields) != 2+k {
		return 0, nil, fmt.Errorf("expected %d indices, got %d tokens total", k, len(fields)-2)
	}
	remove := make([]int, k)
	for i := 0; i < k; i++ {
		idx, err := strconv.Atoi(fields[2+i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid index %q", fields[2+i])
		}
		remove[i] = idx
	}
	return sumVal, remove, nil
}

func validateSolution(arr []int64, reportedSum int64, removeIdx []int) error {
	n := len(arr)
	removed := make([]bool, n)
	for _, idx := range removeIdx {
		if idx < 1 || idx > n {
			return fmt.Errorf("index %d out of range", idx)
		}
		if removed[idx-1] {
			return fmt.Errorf("duplicate index %d", idx)
		}
		removed[idx-1] = true
	}
	remaining := make([]int64, 0, n-len(removeIdx))
	sum := int64(0)
	for i, val := range arr {
		if !removed[i] {
			remaining = append(remaining, val)
			sum += val
		}
	}
	if len(remaining) < 2 {
		return fmt.Errorf("less than two trees remain")
	}
	if remaining[0] != remaining[len(remaining)-1] {
		return fmt.Errorf("first (%d) and last (%d) remaining values differ", remaining[0], remaining[len(remaining)-1])
	}
	if sum != reportedSum {
		return fmt.Errorf("reported sum %d, actual %d", reportedSum, sum)
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeTestCase("sample1", []int64{1, 2, 3, 1, 2}),
		makeTestCase("sample2", []int64{1, -2, 3, 1, -2}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		n := rng.Intn(100) + 2
		arr := make([]int64, n)
		common := int64(rng.Intn(100))
		arr[0], arr[1] = common, common
		for j := 2; j < n; j++ {
			arr[j] = int64(rng.Intn(200) - 100)
		}
		tests = append(tests, makeTestCase(fmt.Sprintf("random_%d", i+1), arr))
	}
	return tests
}

func makeTestCase(name string, arr []int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String(), arr: arr}
}
