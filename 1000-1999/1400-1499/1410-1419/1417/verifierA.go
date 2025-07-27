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
	oracle := filepath.Join(dir, "oracleA.bin")
	cmd := exec.Command("go", "build", "-o", oracle, "1417A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	if len(fields) != 2+n {
		return "", fmt.Errorf("expected %d numbers got %d", 2+n, len(fields))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fields[2+i])
	}
	sb.WriteByte('\n')
	return sb.String(), nil
}

func loadTests() ([]string, error) {
	file, err := os.Open("testcasesA.txt")
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests, err := loadTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, line := range tests {
		input, err := lineToInput(line)
		if err != nil {
			fmt.Printf("bad test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", idx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
