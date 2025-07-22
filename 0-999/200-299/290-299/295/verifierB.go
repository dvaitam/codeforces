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
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "295B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		p := 0
		n := atoi(parts[p])
		p++
		need := n*n + n
		if len(parts) != 1+need {
			fmt.Printf("test %d incomplete\n", idx)
			os.Exit(1)
		}
		matrix := make([][]string, n)
		for i := 0; i < n; i++ {
			row := make([]string, n)
			for j := 0; j < n; j++ {
				row[j] = parts[p]
				p++
			}
			matrix[i] = row
		}
		perm := parts[p:]
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			input.WriteString(strings.Join(matrix[i], " ") + "\n")
		}
		input.WriteString(strings.Join(perm, " ") + "\n")

		expectCmd := exec.Command(oracle)
		expectCmd.Stdin = strings.NewReader(input.String())
		var exp bytes.Buffer
		expectCmd.Stdout = &exp
		if err := expectCmd.Run(); err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(exp.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
