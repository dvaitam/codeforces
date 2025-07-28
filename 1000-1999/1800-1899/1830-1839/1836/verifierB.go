package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseB struct {
	n int64
	k int64
	g int64
}

func parseTestsB() ([]testCaseB, error) {
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var tests []testCaseB
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("invalid test case: %s", line)
		}
		n, err1 := strconv.ParseInt(fields[0], 10, 64)
		k, err2 := strconv.ParseInt(fields[1], 10, 64)
		g, err3 := strconv.ParseInt(fields[2], 10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("invalid numbers: %s", line)
		}
		tests = append(tests, testCaseB{n, k, g})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func expectedB(n, k, g int64) int64 {
	tVal := (g - 1) / 2
	maxSave := tVal * n
	total := k * g
	if maxSave > total {
		maxSave = total
	}
	return (maxSave / g) * g
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestsB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.k, tc.g)
		expect := expectedB(tc.n, tc.k, tc.g)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: non-integer output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
