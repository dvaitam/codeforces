package main

import (
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
	oracle := filepath.Join(dir, "oracleB2")
	cmd := exec.Command("go", "build", "-o", oracle, "1254B2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin string, data []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
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
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	data, err := os.ReadFile("testcasesB2.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}

	expect, err := run(oracle, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle error: %v\n", err)
		os.Exit(1)
	}
	got, err := run(bin, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	expTok := strings.Fields(expect)
	gotTok := strings.Fields(got)
	if len(expTok) != len(gotTok) {
		fmt.Fprintf(os.Stderr, "token count mismatch: expected %d got %d\n", len(expTok), len(gotTok))
		os.Exit(1)
	}
	for i := range expTok {
		if expTok[i] != gotTok[i] {
			fmt.Fprintf(os.Stderr, "mismatch at token %d: expected %s got %s\n", i+1, expTok[i], gotTok[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(expTok))
}
