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
	cmd := exec.Command("go", "build", "-o", oracle, "1496A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(bin, input string) (string, string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), stderr.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for idx := 1; idx <= t; idx++ {
		var n, k int
		var s string
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: read nk: %v\n", idx, err)
			os.Exit(1)
		}
		if _, err := fmt.Fscan(reader, &s); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: read s: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
		want, _, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, stderr, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
