package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "66E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	var lines []string
	process := func(lines []string) {
		if len(lines) == 0 {
			return
		}
		idx++
		input := strings.Join(lines, "\n") + "\n"
		expected, err := run(oracle, input)
		if err != nil {
			fmt.Printf("case %d oracle error: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed:\nexpected: %s\n got: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			process(lines)
			lines = nil
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	process(lines)
	fmt.Printf("All %d tests passed\n", idx)
}
