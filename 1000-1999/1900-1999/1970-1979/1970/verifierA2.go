package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input string
	s     string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleA2")
	cmd := exec.Command("go", "build", "-o", out, "1970A2.go")
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
	for i := 0; i < len(line); i++ {
		if line[i] != '(' && line[i] != ')' {
			return "", fmt.Errorf("invalid character %q", line[i])
		}
	}
	return line, nil
}

func isBalanced(seq string) bool {
	bal := 0
	for i := 0; i < len(seq); i++ {
		if seq[i] == '(' {
			bal++
		} else {
			bal--
			if bal < 0 {
				return false
			}
		}
	}
	return bal == 0
}

func balancedShuffle(t string) string {
	n := len(t)
	type column struct {
		prefix int
		pos    int
		ch     byte
	}
	cols := make([]column, n)
	bal := 0
	for i := 0; i < n; i++ {
		cols[i] = column{prefix: bal, pos: i + 1, ch: t[i]}
		if t[i] == '(' {
			bal++
		} else {
			bal--
		}
	}
	sort.Slice(cols, func(i, j int) bool {
		if cols[i].prefix != cols[j].prefix {
			return cols[i].prefix < cols[j].prefix
		}
		return cols[i].pos > cols[j].pos
	})
	var sb strings.Builder
	for _, c := range cols {
		sb.WriteByte(c.ch)
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		newTestCase("()"),
		newTestCase("(())"),
		newTestCase("()()"),
		newTestCase("(()())"),
		newTestCase("(((())))"),
	}
}

func newTestCase(s string) testCase {
	return testCase{
		input: s + "\n",
		s:     s,
	}
}

func randomBalanced(rnd *rand.Rand, maxPairs int) string {
	pairs := rnd.Intn(maxPairs) + 1
	openRemain := pairs
	closeRemain := pairs
	bal := 0
	var sb strings.Builder
	for openRemain > 0 || closeRemain > 0 {
		if openRemain == 0 {
			sb.WriteByte(')')
			closeRemain--
			bal--
			continue
		}
		if closeRemain == 0 {
			sb.WriteByte('(')
			openRemain--
			bal++
			continue
		}
		if bal == 0 {
			sb.WriteByte('(')
			openRemain--
			bal++
			continue
		}
		if rnd.Intn(openRemain+closeRemain) < openRemain {
			sb.WriteByte('(')
			openRemain--
			bal++
		} else {
			sb.WriteByte(')')
			closeRemain--
			bal--
		}
	}
	return sb.String()
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		s := randomBalanced(rnd, 500) // length up to 1000
		tests = append(tests, newTestCase(s))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/binary")
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
		expSeq, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotSeq, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if len(gotSeq) != len(tc.s) {
			fmt.Fprintf(os.Stderr, "case %d: length mismatch, expected %d got %d\n", idx+1, len(tc.s), len(gotSeq))
			os.Exit(1)
		}
		if !isBalanced(gotSeq) {
			fmt.Fprintf(os.Stderr, "case %d: candidate output is not balanced\n", idx+1)
			os.Exit(1)
		}
		if balancedShuffle(gotSeq) != tc.s {
			fmt.Fprintf(os.Stderr, "case %d: candidate output does not map to input under balanced shuffle\n", idx+1)
			os.Exit(1)
		}
		if gotSeq != expSeq {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s\n", idx+1, expSeq, gotSeq)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
