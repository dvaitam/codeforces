package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD1")
	cmd := exec.Command("go", "build", "-o", oracle, "177D1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return oracle, nil
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

func parseInts(fields []string) ([]int, error) {
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func buildInput(n, m, c int, a, b []int) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%d %d %d\n", n, m, c)
	for i := 0; i < n; i++ {
		if i > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(strconv.Itoa(a[i]))
	}
	builder.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(strconv.Itoa(b[i]))
	}
	builder.WriteByte('\n')
	return builder.String()
}

func parseOutput(out string, expectedLen int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expectedLen {
		return nil, fmt.Errorf("expected %d numbers, got %d", expectedLen, len(fields))
	}
	return parseInts(fields)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1<<25)

	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		caseNum++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "case %d: not enough data\n", caseNum)
			os.Exit(1)
		}
		header, err := parseInts(fields[:3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		n, m, c := header[0], header[1], header[2]
		expectedCount := 3 + n + m
		if len(fields) != expectedCount {
			fmt.Fprintf(os.Stderr, "case %d: expected %d numbers, got %d\n", caseNum, expectedCount, len(fields))
			os.Exit(1)
		}
		aVals, err := parseInts(fields[3 : 3+n])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		bVals, err := parseInts(fields[3+n:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", caseNum, err)
			os.Exit(1)
		}

		input := buildInput(n, m, c, aVals, bVals)

		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", caseNum, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(binary, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}

		expVals, err := parseOutput(expOut, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", caseNum, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", caseNum, err)
			os.Exit(1)
		}

		match := true
		for i := 0; i < n; i++ {
			if gotVals[i] != expVals[i] {
				match = false
				break
			}
		}
		if !match {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", caseNum, expVals, gotVals)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}
