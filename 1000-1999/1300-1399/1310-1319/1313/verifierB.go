package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	exePath := filepath.Join(dir, "oracleB")
	src := filepath.Join(dir, "1313B.go")
	cmd := exec.Command("go", "build", "-o", exePath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return exePath, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Int63n(1_000_000_000) + 1
	x := rng.Int63n(n) + 1
	y := rng.Int63n(n) + 1
	return fmt.Sprintf("1\n%d %d %d\n", n, x, y)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// generate tests
	tests := make([]string, 0, 120)
	tests = append(tests, "1\n1 1 1\n")
	tests = append(tests, fmt.Sprintf("1\n1000000000 1 1000000000\n"))
	for i := 0; i < 118; i++ {
		tests = append(tests, generateCase(rng))
	}

	for i, input := range tests {
		expect, err := runProg(filepath.Join("./", oracle), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error: %v\n", err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\ngot: %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
