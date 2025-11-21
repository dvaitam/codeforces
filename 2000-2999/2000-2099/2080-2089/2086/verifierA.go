package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
	count int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("unable to locate verifier file path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2086A-")
	if err != nil {
		return "", nil, err
	}
	oraclePath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oraclePath, "2086A.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return oraclePath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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
	return strings.TrimSpace(stdout.String()), nil
}

func makeInput(values []int64) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(values)))
	sb.WriteByte('\n')
	for _, v := range values {
		sb.WriteString(strconv.FormatInt(v, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateTests() []testCase {
	const maxN int64 = 100000000
	tests := []testCase{
		{
			name:  "single_min",
			input: makeInput([]int64{1}),
			count: 1,
		},
		{
			name:  "small_sequence",
			input: makeInput([]int64{1, 2, 3, 4, 5}),
			count: 5,
		},
		{
			name:  "edge_values",
			input: makeInput([]int64{1, maxN, maxN, 1, maxN}),
			count: 5,
		},
	}

	alternating := make([]int64, 10000)
	for i := range alternating {
		if i%2 == 0 {
			alternating[i] = maxN
		} else {
			alternating[i] = 1
		}
	}
	tests = append(tests, testCase{
		name:  "max_t_alternating_edges",
		input: makeInput(alternating),
		count: len(alternating),
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		t := rng.Intn(100) + 1
		values := make([]int64, t)
		for j := 0; j < t; j++ {
			values[j] = int64(rng.Intn(100000000)) + 1
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: makeInput(values),
			count: t,
		})
	}

	huge := make([]int64, 10000)
	for i := range huge {
		huge[i] = int64(rng.Intn(100000000)) + 1
	}
	tests = append(tests, testCase{
		name:  "max_t_random",
		input: makeInput(huge),
		count: len(huge),
	})

	return tests
}

func compareOutputs(expected, actual string, count int) error {
	expFields := strings.Fields(expected)
	if len(expFields) != count {
		return fmt.Errorf("oracle produced %d answers, expected %d", len(expFields), count)
	}
	actFields := strings.Fields(actual)
	if len(actFields) != count {
		return fmt.Errorf("expected %d answers, but got %d", count, len(actFields))
	}
	for i := 0; i < count; i++ {
		if expFields[i] != actFields[i] {
			return fmt.Errorf("mismatch at case %d: expected %s got %s", i+1, expFields[i], actFields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expected, err := runBinary(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		actual, err := runBinary(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := compareOutputs(expected, actual, tc.count); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.name, err, tc.input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
