package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func buildOracle(src, outName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, outName)
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierX.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle("1250D.go", "oracleD")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}
	cmdO := exec.Command(oracle)
	cmdO.Stdin = bytes.NewReader(data)
	outO, err := cmdO.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
		os.Exit(1)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}
	if string(bytes.TrimSpace(out)) != string(bytes.TrimSpace(outO)) {
		fmt.Println("output mismatch")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
