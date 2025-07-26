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
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "1168C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runCase(bin, oracle, input string) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing test case %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		q, _ := strconv.Atoi(parts[1])
		need := 2 + n + 2*q
		if len(parts) != need {
			fmt.Fprintf(os.Stderr, "test %d expected %d numbers got %d\n", i+1, need, len(parts))
			os.Exit(1)
		}
		arr := strings.Join(parts[2:2+n], " ")
		queries := parts[2+n:]
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		sb.WriteString(arr)
		sb.WriteByte('\n')
		for j := 0; j < len(queries); j += 2 {
			sb.WriteString(queries[j])
			sb.WriteByte(' ')
			sb.WriteString(queries[j+1])
			sb.WriteByte('\n')
		}
		if err := runCase(bin, oracle, sb.String()); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
