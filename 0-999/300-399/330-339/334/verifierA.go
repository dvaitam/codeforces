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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func checkOutput(n int, output string) error {
	output = strings.TrimSpace(output)
	lines := strings.Split(output, "\n")
	if len(lines) != n {
		return fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	used := make([]bool, n*n+1)
	target := n * (n*n + 1) / 2
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != n {
			return fmt.Errorf("line %d: expected %d numbers, got %d", i+1, n, len(fields))
		}
		sum := 0
		for _, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("line %d: bad integer %q", i+1, f)
			}
			if v < 1 || v > n*n {
				return fmt.Errorf("line %d: value %d out of range", i+1, v)
			}
			if used[v] {
				return fmt.Errorf("number %d repeated", v)
			}
			used[v] = true
			sum += v
		}
		if sum != target {
			return fmt.Errorf("line %d: sum %d expected %d", i+1, sum, target)
		}
	}
	for v := 1; v <= n*n; v++ {
		if !used[v] {
			return fmt.Errorf("number %d missing", v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesA.txt: %v\n", err)
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
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n", n)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if err := checkOutput(n, out); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
