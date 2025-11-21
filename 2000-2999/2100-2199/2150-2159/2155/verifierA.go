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

const refSource = "2000-2999/2100-2199/2150-2159/2155/2155A.go"

type testCase struct {
	name  string
	input string
	nVals []int
}

func runProgram(bin, input string) (string, error) {
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

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "single", input: "1\n1\n", nVals: []int{1}},
		{name: "small", input: "2\n2\n5\n", nVals: []int{2, 5}},
		{name: "edge", input: "1\n1000000\n", nVals: []int{1000000}},
	}
}

func randomTests(count int) []testCase {
	tests := make([]testCase, 0, count)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		t := rng.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		nVals := make([]int, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(1_000_000) + 1
			nVals[j] = n
			sb.WriteString(fmt.Sprintf("%d\n", n))
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: sb.String(),
			nVals: nVals,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refBin, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	tests := append(manualTests(), randomTests(100)...)
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutput(refOut, len(tc.nVals))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, len(tc.nVals))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		for i := 0; i < len(tc.nVals); i++ {
			expected := int64(2 * (tc.nVals[i] - 1))
			if candVals[i] != expected {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, expected, candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
			if refVals[i] != expected {
				fmt.Fprintf(os.Stderr, "reference produced unexpected value on test %d (%s) case %d: expected %d got %d\n",
					idx+1, tc.name, i+1, expected, refVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
