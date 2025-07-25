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
	cmd := exec.Command("go", "build", "-o", oracle, "867A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcasesA.txt: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}

	for i := 1; i <= t; i++ {
		var n int
		var s string
		if _, err := fmt.Fscan(reader, &n); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: read n: %v\n", i, err)
			os.Exit(1)
		}
		if _, err := fmt.Fscan(reader, &s); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: read s: %v\n", i, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n%s\n", n, s)
		want, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "test %d failed:\nexpected: %s\ngot: %s\ninput:\n%s", i, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
