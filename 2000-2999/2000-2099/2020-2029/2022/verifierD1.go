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
	roles []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2022D1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD1")
	cmd := exec.Command("go", "build", "-o", outPath, "2022D1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		n := len(tc.roles)
		sb.WriteString(fmt.Sprintf("%d manual\n", n))
		for i, v := range tc.roles {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{roles: []int{0, -1, 1}},
		{roles: []int{-1, 1, 0, 0}},
		{roles: []int{1, 1, -1, 1, 1}},
		{roles: []int{0, 0, 0, -1, 0, 0}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	totalN := 0
	for len(tests) < 30 && totalN < 80000 {
		n := rng.Intn(20) + 3
		if totalN+n > 80000 {
			n = 80000 - totalN
		}
		roles := make([]int, n)
		pos := rng.Intn(n)
		for i := 0; i < n; i++ {
			if i == pos {
				roles[i] = -1
				continue
			}
			if rng.Intn(2) == 0 {
				roles[i] = 0
			} else {
				roles[i] = 1
			}
		}
		tests = append(tests, testCase{roles: roles})
		totalN += n
	}
	for len(tests) < 40 {
		n := rng.Intn(1000) + 1000
		roles := make([]int, n)
		pos := rng.Intn(n)
		for i := 0; i < n; i++ {
			if i == pos {
				roles[i] = -1
				continue
			}
			if rng.Intn(2) == 0 {
				roles[i] = 0
			} else {
				roles[i] = 1
			}
		}
		tests = append(tests, testCase{roles: roles})
		totalN += n
		if totalN >= 100000 {
			break
		}
	}
	return tests
}

func stressTests() []testCase {
	n := 100000
	roles := make([]int, n)
	for i := range roles {
		roles[i] = 1
	}
	roles[n/2] = -1
	return []testCase{{roles: roles}}
}

func compareOutputs(expected, actual string, count int) error {
	exp := strings.Fields(expected)
	act := strings.Fields(actual)
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d answers, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("expected %d answers, got %d", count, len(act))
	}
	for i := 0; i < count; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at test %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	tests = append(tests, stressTests()...)
	input := buildInput(tests)

	expected, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	actual, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := compareOutputs(expected, actual, len(tests)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nInput:\n%s\nExpected:\n%s\nActual:\n%s\n", err, input, expected, actual)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
