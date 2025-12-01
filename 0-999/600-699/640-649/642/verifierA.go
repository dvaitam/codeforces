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
	outPath := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, "642A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(bin string, input []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.Bytes(), nil
}

func deterministicTests() [][]byte {
	var tests [][]byte
	var sb strings.Builder
	sb.WriteString("1 1\n")
	sb.WriteString("250 1\n")
	sb.WriteString("-1\n")
	sb.WriteString("-1 -1\n")
	tests = append(tests, []byte(sb.String()))

	sb.Reset()
	sb.WriteString("1 1\n")
	sb.WriteString("250 2\n")
	sb.WriteString("0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 -1\n-1 -1\n")
	tests = append(tests, []byte(sb.String()))

	return tests
}

func randomTests(count int) [][]byte {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([][]byte, 0, count)
	for i := 0; i < count; i++ {
		t := rnd.Intn(5) + 1
		p := rnd.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", t, p))
		for j := 0; j < p; j++ {
			tl := rnd.Intn(1000) + 250
			testsCount := rnd.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", tl, testsCount))
		}
		events := rnd.Intn(3) + 1
		for e := 0; e < events; e++ {
			sb.WriteString(fmt.Sprintf("%d ", rnd.Intn(p)))
		}
		sb.WriteString("-1\n")
		sb.WriteString("-1 -1\n")
		tests = append(tests, []byte(sb.String()))
	}
	return tests
}

func compareOutputs(exp, got []byte) error {
	expStr := strings.TrimSpace(string(exp))
	gotStr := strings.TrimSpace(string(got))
	if expStr != gotStr {
		return fmt.Errorf("expected output differs")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(50)...)

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(expOut, gotOut); err != nil {
			fmt.Fprintf(os.Stderr, "case %d mismatch: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
