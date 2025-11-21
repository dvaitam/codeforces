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

type testCase struct {
	input string
	t     int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", out, "1999E.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return out, nil
}

func runProgram(bin string, input string) (string, error) {
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

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(lines))
	}
	return lines, nil
}

func deterministicTests() []testCase {
	tests := []testCase{
		formatTest([][2]int64{{1, 3}, {2, 4}}),
		formatTest([][2]int64{{199999, 200000}, {1, 200000}}),
		formatTest([][2]int64{{5, 6}}),
	}
	return tests
}

func formatTest(pairs [][2]int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(pairs)))
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return testCase{input: sb.String(), t: len(pairs)}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		t := rnd.Intn(20) + 1
		pairs := make([][2]int64, t)
		for j := 0; j < t; j++ {
			l := rnd.Int63n(200000-1) + 1
			r := l + rnd.Int63n(200000-l) + 1
			pairs[j] = [2]int64{l, r}
		}
		tests = append(tests, formatTest(pairs))
	}
	// add max case
	tests = append(tests, formatTest([][2]int64{{1, 200000}}))
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(200)...)

	for idx, tc := range tests {
		expOut, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVals, err := parseOutput(expOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		for i := 0; i < tc.t; i++ {
			if gotVals[i] != expVals[i] {
				fmt.Fprintf(os.Stderr, "case %d test %d mismatch: expected %s got %s\n", idx+1, i+1, expVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
