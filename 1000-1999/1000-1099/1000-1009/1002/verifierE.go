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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func expectedOps(n int) []string {
	ops := make([]string, 0, 2*n+3)
	ops = append(ops, fmt.Sprintf("X %d", n+1))
	for i := 1; i <= n+1; i++ {
		ops = append(ops, fmt.Sprintf("H %d", i))
	}
	ops = append(ops, "ORACLE")
	for i := 1; i <= n; i++ {
		ops = append(ops, fmt.Sprintf("H %d", i))
	}
	return ops
}

func checkCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	lines := strings.FieldsFunc(strings.TrimSpace(out), func(r rune) bool { return r == '\n' || r == '\r' })
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("failed to parse op count: %v", err)
	}
	ops := expectedOps(n)
	if k != len(ops) {
		return fmt.Errorf("expected %d operations, got %d", len(ops), k)
	}
	if len(lines)-1 != len(ops) {
		return fmt.Errorf("expected %d lines, got %d", len(ops), len(lines)-1)
	}
	for i, op := range ops {
		if strings.TrimSpace(lines[i+1]) != op {
			return fmt.Errorf("line %d: expected %q got %q", i+2, op, strings.TrimSpace(lines[i+1]))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesE.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad N on line %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		idx++
		if err := checkCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
