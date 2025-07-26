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

func buildOracle() (string, error) {
	ref := "./oracleA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "775A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return ref, nil
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func lineToInput(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return "", fmt.Errorf("bad test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	a, err := strconv.Atoi(fields[2])
	if err != nil {
		return "", err
	}
	if len(fields) != 3+n*m {
		return "", fmt.Errorf("expected %d numbers got %d", 3+n*m, len(fields))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, a))
	idx := 3
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			sb.WriteString(fields[idx])
			idx++
			if j+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String(), nil
}

func loadTests(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var tests []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tests = append(tests, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests, err := loadTests("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, line := range tests {
		input, err := lineToInput(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
