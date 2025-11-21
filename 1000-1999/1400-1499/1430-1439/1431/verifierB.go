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
	out := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", out, "1431B.go")
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

func parseOutput(out string, expectedLines int) ([]int, error) {
	if strings.TrimSpace(out) == "" && expectedLines == 0 {
		return nil, nil
	}
	lines := strings.Fields(out)
	if len(lines) != expectedLines {
		return nil, fmt.Errorf("expected %d answers, got %d", expectedLines, len(lines))
	}
	res := make([]int, expectedLines)
	for i, line := range lines {
		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", line)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative answer %d at case %d", val, i+1)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	tests := []testCase{
		formatTest([]string{"v", "vv", "w", "vw", "wv"}),
		formatTest([]string{"vvvv", "wwww", "vwwv", "vvwvv"}),
		formatTest([]string{"v" + strings.Repeat("w", 99)}),
	}
	return tests
}

func formatTest(stringsList []string) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(stringsList)))
	for _, s := range stringsList {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), t: len(stringsList)}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		t := rnd.Intn(100) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for j := 0; j < t; j++ {
			length := rnd.Intn(100) + 1
			for k := 0; k < length; k++ {
				if rnd.Intn(2) == 0 {
					sb.WriteByte('v')
				} else {
					sb.WriteByte('w')
				}
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{input: sb.String(), t: t})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
				fmt.Fprintf(os.Stderr, "case %d test %d mismatch: expected %d got %d\n", idx+1, i+1, expVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
