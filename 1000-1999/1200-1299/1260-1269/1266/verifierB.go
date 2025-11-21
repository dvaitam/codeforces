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
	cmd := exec.Command("go", "build", "-o", out, "1266B.go")
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

func parseOutput(out string, expectedLines int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expectedLines {
		return nil, fmt.Errorf("expected %d outputs, got %d", expectedLines, len(lines))
	}
	for _, line := range lines {
		lineUpper := strings.ToUpper(line)
		if lineUpper != "YES" && lineUpper != "NO" {
			return nil, fmt.Errorf("invalid verdict %q", line)
		}
	}
	return lines, nil
}

func deterministicTests() []testCase {
	tests := []testCase{
		formatTest([]int64{1, 2, 3, 4, 5, 6, 7}),
		formatTest([]int64{14, 15, 16, 21, 28, 29}),
		formatTest([]int64{100, 101, 102, 103, 104, 105, 106}),
	}
	return tests
}

func formatTest(values []int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(values)))
	for i, v := range values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), t: len(values)}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		t := rnd.Intn(1000) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for j := 0; j < t; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			choice := rnd.Intn(10)
			var x int64
			switch choice {
			case 0:
				x = int64(rnd.Intn(20) + 1)
			case 1:
				x = int64(rnd.Intn(1000) + 1)
			case 2:
				x = int64(rnd.Intn(14) + 1)
			case 3:
				x = 14 * int64(rnd.Intn(1000000)+1)
			default:
				x = rand.Int63n(1_000_000_000_000_000_000) + 1
			}
			sb.WriteString(strconv.FormatInt(x, 10))
		}
		sb.WriteByte('\n')
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
			if strings.ToUpper(gotVals[i]) != strings.ToUpper(expVals[i]) {
				fmt.Fprintf(os.Stderr, "case %d query %d mismatch: expected %s got %s\n", idx+1, i+1, expVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
