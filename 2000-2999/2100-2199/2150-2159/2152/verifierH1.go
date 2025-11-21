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

const refSource = "2000-2999/2100-2199/2150-2159/2152/2152H1.go"

type testCase struct {
	name  string
	input string
	t     int
	qs    []int64
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
		{
			name:  "simple_line",
			input: "1\n3\n1 2 5\n2 3 7\n2\n0\n5\n",
			t:     1,
			qs:    []int64{0, 5},
		},
		{
			name:  "single_node",
			input: "1\n1\n2\n10\n20\n",
			t:     1,
			qs:    []int64{10, 20},
		},
		{
			name:  "two_test_cases",
			input: "2\n2\n1 2 3\n1\n0\n3\n1 2 4\n1 3 6\n1\n2\n",
			t:     2,
			qs:    []int64{0, 2},
		},
	}
}

func randomTests(count int) []testCase {
	tests := make([]testCase, 0, count)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		testCases := rng.Intn(2) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", testCases))
		qCounts := make([]int, testCases)
		var totalQueries int
		for tc := 0; tc < testCases; tc++ {
			n := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("%d\n", n))
			for j := 0; j < n-1; j++ {
				u := j + 1
				v := j + 2
				w := rng.Intn(10) + 1
				sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
			}
			q := rng.Intn(5) + 1
			qCounts[tc] = q
			totalQueries += q
			sb.WriteString(fmt.Sprintf("%d\n", q))
			for j := 0; j < q; j++ {
				l := rng.Intn(20)
				sb.WriteString(fmt.Sprintf("%d\n", l))
			}
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: sb.String(),
			t:     totalQueries,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH1.go /path/to/binary")
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

	tests := append(manualTests(), randomTests(50)...)
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutput(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		for i := 0; i < tc.t; i++ {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at answer %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
