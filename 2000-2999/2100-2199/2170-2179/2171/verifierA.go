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
	out := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", out, "2171A.go")
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

func deterministicTests() []testCase {
	return []testCase{
		{input: "88 94 95\n"},
		{input: "100 80 81\n"},
		{input: "98 99 98\n"},
		{input: "95 86 85\n"},
		{input: "80 80 80\n"},
		{input: "90 100 95\n"},
	}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, count)
	for i := 0; i < count; i++ {
		g := rnd.Intn(21) + 80
		c := rnd.Intn(21) + 80
		l := rnd.Intn(21) + 80
		tests[i] = testCase{input: fmt.Sprintf("%d %d %d\n", g, c, l)}
	}
	return tests
}

func parseOutput(out string) (string, error) {
	line := strings.TrimSpace(out)
	if line == "" {
		return "", fmt.Errorf("empty output")
	}
	fields := strings.Fields(line)
	if len(fields) == 2 {
		if fields[0] != "final" {
			return "", fmt.Errorf("invalid format %q", line)
		}
		if _, err := strconv.Atoi(fields[1]); err != nil {
			return "", fmt.Errorf("invalid score %q", fields[1])
		}
		return line, nil
	}
	if len(fields) == 2 && fields[0] == "check" && fields[1] == "again" {
		return line, nil
	}
	if line == "check again" {
		return line, nil
	}
	return "", fmt.Errorf("invalid output %q", line)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		expLine, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotLine, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if gotLine != expLine {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %q got %q\n", idx+1, expLine, gotLine)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
