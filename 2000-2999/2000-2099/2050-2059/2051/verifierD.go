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
	n   int
	x   int64
	y   int64
	arr []int64
}

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/solution")
		os.Exit(1)
	}
	target := args[0]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expected, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\n%s", err, expected)
		os.Exit(1)
	}
	actual, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s", err, actual)
		os.Exit(1)
	}

	expTokens := parseTokens(expected)
	actTokens := parseTokens(actual)
	if len(expTokens) != len(actTokens) {
		fmt.Fprintf(os.Stderr, "token count mismatch: expected %d got %d\n", len(expTokens), len(actTokens))
		fmt.Fprintf(os.Stderr, "Input:\n%s", input)
		os.Exit(1)
	}
	for i := range expTokens {
		if expTokens[i] != actTokens[i] {
			fmt.Fprintf(os.Stderr, "mismatch on answer %d: expected %s got %s\n", i+1, expTokens[i], actTokens[i])
			fmt.Fprintf(os.Stderr, "Input:\n%s", input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2051D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, "2051D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
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
		return stdout.String() + stderr.String(), fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.x, tc.y))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 4, x: 8, y: 10, arr: []int64{4, 6, 3, 6}},
		{n: 6, x: 2, y: 27, arr: []int64{4, 9, 6, 3, 4, 5}},
		{n: 3, x: 8, y: 10, arr: []int64{3, 2, 1}},
		{n: 6, x: 8, y: 8, arr: []int64{1, 1, 2, 2, 2, 3}},
		{n: 5, x: 1, y: 15, arr: []int64{1, 2, 3, 4, 5}},
		{n: 3, x: 3, y: 6, arr: []int64{3, 2, 1}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	sumN := 0
	for len(tests) < cap(tests) && sumN < 180000 {
		n := rng.Intn(10) + 3
		if sumN+n > 200000 {
			break
		}
		arr := make([]int64, n)
		var total int64
		for i := 0; i < n; i++ {
			val := int64(rng.Intn(1_000_000_000) + 1)
			arr[i] = val
			total += val
		}
		if total == 0 {
			continue
		}
		low := rng.Int63n(total) + 1
		high := low + int64(rng.Intn(1000))
		if high > total {
			high = total
		}
		tc := testCase{
			n:   n,
			x:   low,
			y:   high,
			arr: arr,
		}
		tests = append(tests, tc)
		sumN += n
	}
	return tests
}

func parseTokens(out string) []string {
	return strings.Fields(strings.TrimSpace(out))
}
