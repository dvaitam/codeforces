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
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "932F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(fields []string) (string, error) {
	if len(fields) < 1 {
		return "", fmt.Errorf("empty line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	need := 1 + 2*n + 2*(n-1)
	if len(fields) != need {
		return "", fmt.Errorf("field count mismatch")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	idx := 1
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fields[idx])
		idx++
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fields[idx])
		idx++
	}
	sb.WriteByte('\n')
	for i := 0; i < n-1; i++ {
		sb.WriteString(fmt.Sprintf("%s %s\n", fields[idx], fields[idx+1]))
		idx += 2
	}
	return sb.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	const testcasesRaw = `4 3 3 5 1 4 2 1 2 1 2 2 3 2 4
2 4 5 1 5 1 2
2 2 4 3 2 1 2
2 1 2 5 5 1 2
2 2 1 1 2 1 2
2 2 3 3 2 1 2
2 2 4 3 1 1 2
3 2 2 3 1 3 3 1 2 2 3
2 3 3 3 4 1 2
2 4 4 2 1 1 2
2 3 4 1 5 1 2
3 4 5 1 4 1 2 1 2 1 3
2 4 3 5 3 1 2
3 1 5 3 3 1 4 1 2 1 3
3 5 5 3 2 3 3 1 2 2 3
4 3 3 2 1 2 3 4 2 1 2 1 3 3 4
4 4 1 2 5 3 3 4 4 1 2 1 3 3 4
2 4 3 2 2 1 2
4 4 1 2 4 3 2 1 4 1 2 1 3 2 4
4 2 5 4 4 3 4 3 3 1 2 2 3 1 4
2 4 5 2 4 1 2
2 1 4 3 5 1 2
2 3 5 1 3 1 2
4 3 4 3 3 3 2 5 1 1 2 2 3 2 4
4 3 4 3 5 3 3 3 3 1 2 2 3 1 4
4 4 3 3 5 2 5 2 2 1 2 2 3 2 4
4 1 4 2 5 5 5 4 3 1 2 1 3 1 4
2 5 4 5 2 1 2
4 2 1 4 2 2 1 2 1 1 2 1 3 2 4
2 5 1 4 4 1 2
3 5 1 5 2 2 3 1 2 2 3
3 3 4 1 5 3 1 1 2 1 3
3 5 5 3 5 3 3 1 2 2 3
3 5 5 3 4 2 2 1 2 2 3
2 2 5 1 3 1 2
3 1 3 5 5 5 2 1 2 2 3
3 4 2 4 3 4 4 1 2 1 3
4 5 3 3 4 3 4 1 3 1 2 2 3 2 4
2 4 1 1 5 1 2
2 3 1 2 4 1 2
3 5 5 3 2 4 1 1 2 2 3
3 2 3 3 4 5 4 1 2 1 3
2 3 3 5 3 1 2
4 5 1 4 3 3 5 2 1 1 2 2 3 3 4
3 1 5 2 1 4 4 1 2 2 3
3 4 2 3 4 1 3 1 2 2 3
4 3 4 3 2 2 4 5 4 1 2 1 3 2 4
4 5 1 1 2 3 1 1 1 1 2 2 3 1 4
3 4 1 1 1 2 5 1 2 1 3
3 4 2 2 5 2 4 1 2 1 3
4 2 1 5 1 4 1 3 1 1 2 2 3 3 4
2 1 4 1 3 1 2
3 3 3 1 2 3 2 1 2 1 3
3 1 5 2 3 3 5 1 2 1 3
3 4 4 2 2 1 3 1 2 1 3
2 1 5 5 4 1 2
3 4 4 3 1 5 1 1 2 1 3
2 4 2 1 1 1 2
2 2 5 4 3 1 2
3 3 5 3 4 4 2 1 2 2 3
4 4 3 5 1 5 5 2 2 1 2 1 3 1 4
2 1 2 2 2 1 2
3 4 3 1 4 4 3 1 2 2 3
3 4 1 2 2 2 1 1 2 2 3
3 2 5 3 1 3 2 1 2 1 3
2 4 1 3 3 1 2
2 1 4 4 2 1 2
4 3 3 4 2 2 1 2 1 1 2 1 3 3 4
3 3 3 2 4 3 3 1 2 2 3
2 2 2 2 1 1 2
4 3 1 4 1 4 4 2 1 1 2 2 3 1 4
4 1 2 4 1 4 1 2 1 1 2 2 3 1 4
2 3 4 4 3 1 2
2 2 2 5 4 1 2
4 2 1 3 1 1 5 2 5 1 2 1 3 3 4
2 4 5 3 4 1 2
4 4 5 2 1 3 4 2 5 1 2 1 3 2 4
3 1 3 4 3 2 2 1 2 2 3
2 1 4 5 2 1 2
2 4 2 1 2 1 2
3 1 1 4 3 4 5 1 2 2 3
4 2 5 1 4 2 1 1 4 1 2 1 3 2 4
4 2 3 1 3 2 2 3 4 1 2 2 3 1 4
2 5 5 1 1 1 2
2 2 3 4 3 1 2
2 5 4 2 3 1 2
3 3 1 2 5 5 2 1 2 2 3
4 5 2 2 4 4 4 1 3 1 2 2 3 1 4
3 5 4 2 4 1 2 1 2 2 3
4 1 1 2 4 1 2 1 1 1 2 2 3 3 4
3 4 4 2 2 3 2 1 2 2 3
2 3 3 2 1 1 2
3 3 2 2 2 4 4 1 2 1 3
2 4 1 3 2 1 2
3 4 5 5 1 2 4 1 2 1 3
4 5 5 5 2 4 1 3 1 1 2 2 3 3 4
3 5 4 4 5 5 1 1 2 2 3
3 5 1 5 2 3 3 1 2 1 3
2 3 5 5 4 1 2
2 3 2 1 1 1 2`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		input, err := buildInput(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase %d: %v\n", idx, err)
			os.Exit(1)
		}
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
