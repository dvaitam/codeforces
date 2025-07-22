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
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "421A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBinary(bin, input string) (string, error) {
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

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	var t int
	fmt.Sscan(scanner.Text(), &t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "unexpected EOF reading case %d\n", i+1)
			os.Exit(1)
		}
		line1 := strings.TrimSpace(scanner.Text())
		var n, a, b int
		fmt.Sscan(line1, &n, &a, &b)

		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing Arthur list for case %d\n", i+1)
			os.Exit(1)
		}
		line2 := strings.TrimSpace(scanner.Text())
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing Alexander list for case %d\n", i+1)
			os.Exit(1)
		}
		line3 := strings.TrimSpace(scanner.Text())

		input := line1 + "\n" + line2 + "\n" + line3 + "\n"

		expect, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}

		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}

		if got != expect {
			fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", t)
}
