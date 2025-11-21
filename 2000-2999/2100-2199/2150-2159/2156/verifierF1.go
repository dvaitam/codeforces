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

const refSource = "2000-2999/2100-2199/2150-2159/2156/2156F1.go"

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

func parseOutput(out string, expected int) ([][]int, error) {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	res := make([][]int, 0, expected)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		seq := make([]int, len(fields))
		for i, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", f)
			}
			seq[i] = val
		}
		res = append(res, seq)
	}
	if len(res) != expected {
		return nil, fmt.Errorf("expected %d sequences, got %d", expected, len(res))
	}
	return res, nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "single", input: "1\n1\n1\n2\n1 2\n1 3\n", nVals: []int{1}},
		{name: "two_cases", input: "2\n2\n1 2\n1 3\n3\n1 2 3\n4\n1 2 3 4\n1 5 6 7\n1 2 3 4\n", nVals: []int{2, 3}},
	}
}

func randomTests(count int) []testCase {
	tests := make([]testCase, 0, count)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		nVals := make([]int, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(4) + 1
			nVals[j] = n
			sb.WriteString(fmt.Sprintf("%d\n", n))
			for k := 0; k < n; k++ {
				lenArr := rng.Intn(4) + 1
				sb.WriteString(fmt.Sprintf("%d", lenArr))
				for x := 0; x < lenArr; x++ {
					sb.WriteString(fmt.Sprintf(" %d", rng.Intn(10)))
				}
				sb.WriteByte('\n')
			}
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
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
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
		refSeqs, err := parseOutput(refOut, len(tc.nVals))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candSeqs, err := parseOutput(candOut, len(tc.nVals))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		for caseIdx := range tc.nVals {
			if len(refSeqs[caseIdx]) != len(candSeqs[caseIdx]) {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d length mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, caseIdx+1, len(refSeqs[caseIdx]), len(candSeqs[caseIdx]), tc.input, refOut, candOut)
				os.Exit(1)
			}
			for i := range refSeqs[caseIdx] {
				if refSeqs[caseIdx][i] != candSeqs[caseIdx][i] {
					fmt.Fprintf(os.Stderr, "test %d (%s) case %d position %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
						idx+1, tc.name, caseIdx+1, i+1, refSeqs[caseIdx][i], candSeqs[caseIdx][i], tc.input, refOut, candOut)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
