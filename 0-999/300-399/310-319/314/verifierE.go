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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	outPath := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "314E.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return outPath, nil
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(n int, s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 32)
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	sb.WriteString(s)
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []string {
	var tests []string
	tests = append(tests, buildInput(1, "a"))
	tests = append(tests, buildInput(2, "aa"))
	tests = append(tests, buildInput(2, "a?"))
	tests = append(tests, buildInput(2, "??"))
	tests = append(tests, buildInput(4, "abba"))
	tests = append(tests, buildInput(4, "a??b"))
	tests = append(tests, buildInput(6, "abcdef"))
	tests = append(tests, buildInput(6, "a?b?c?"))
	tests = append(tests, buildInput(8, "?a?b?c?d"))
	tests = append(tests, buildInput(100000, strings.Repeat("?", 100000)))
	return tests
}

func randomString(n int, rnd *rand.Rand) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if rnd.Intn(3) == 0 {
			buf[i] = '?'
		} else {
			c := byte('a' + rnd.Intn(26))
			if c == 'x' {
				c = 'y'
			}
			buf[i] = c
		}
	}
	return string(buf)
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%20 == 0:
			n = rnd.Intn(100000) + 1
		case i%5 == 0:
			n = rnd.Intn(1000) + 1
		default:
			n = rnd.Intn(50) + 1
		}
		tests = append(tests, buildInput(n, randomString(n, rnd)))
	}
	return tests
}

func parseOutput(out string) (uint64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		if len(fields) == 0 {
			return 0, fmt.Errorf("empty output")
		}
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
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
	tests = append(tests, randomTests(300)...)

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
