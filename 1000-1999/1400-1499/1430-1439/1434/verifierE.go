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
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", out, "1434E.go")
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

func parseOutput(out string) (string, error) {
	line := strings.TrimSpace(out)
	if line == "" {
		return "", fmt.Errorf("empty output")
	}
	line = strings.ToUpper(line)
	if line != "YES" && line != "NO" {
		return "", fmt.Errorf("invalid verdict %q", line)
	}
	return line, nil
}

func deterministicTests() []testCase {
	tests := []testCase{
		formatTest([][]int{{1}}),
		formatTest([][]int{{1, 2}}),
		formatTest([][]int{{1, 2, 4}, {10, 11}}),
		formatTest([][]int{{1, 3, 6, 10}, {2, 5, 9}}),
	}
	return tests
}

func formatTest(seqs [][]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(seqs)))
	for _, seq := range seqs {
		sb.WriteString(fmt.Sprintf("%d\n", len(seq)))
		for i, v := range seq {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String()}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	totalLen := 0
	for len(tests) < count && totalLen < 100000 {
		n := rnd.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if totalLen >= 100000 {
				n = i
				break
			}
			length := rnd.Intn(20) + 1
			if totalLen+length > 100000 {
				length = 100000 - totalLen
			}
			totalLen += length
			sb.WriteString(fmt.Sprintf("%d\n", length))
			cur := 0
			for j := 0; j < length; j++ {
				cur += rnd.Intn(5) + 1
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(cur))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{input: sb.String()})
	}
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
		expAns, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotAns, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if gotAns != expAns {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s\n", idx+1, expAns, gotAns)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
