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
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "1637B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func nextCase(scanner *bufio.Scanner) (string, bool) {
	for scanner.Scan() {
		nLine := strings.TrimSpace(scanner.Text())
		if nLine == "" {
			continue
		}
		if !scanner.Scan() {
			return "", false
		}
		arrLine := strings.TrimSpace(scanner.Text())
		input := fmt.Sprintf("1\n%s\n%s\n", nLine, arrLine)
		return input, true
	}
	return "", false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
	for {
		in, ok := nextCase(scanner)
		if !ok {
			break
		}
		idx++
		exp, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
