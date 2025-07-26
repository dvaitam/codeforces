package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func buildOracle() (string, error) {
	oracle := "oracleB"
	cmd := exec.Command("go", "build", "-o", oracle, "1164B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + oracle, nil
}

func run(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
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

	expect, err := run(oracle)
	if err != nil {
		fmt.Fprintln(os.Stderr, "oracle runtime error:", err)
		os.Exit(1)
	}

	for i := 1; i <= 100; i++ {
		got, err := run(bin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
