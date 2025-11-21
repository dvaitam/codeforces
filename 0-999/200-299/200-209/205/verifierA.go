package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "205A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(times []int64) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(times)))
	sb.WriteByte('\n')
	for i, t := range times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(t, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomTest(n int, maxVal int64, rnd *rand.Rand) string {
	if n == 0 {
		n = 1
	}
	if maxVal <= 0 {
		maxVal = 1
	}
	times := make([]int64, n)
	for i := 0; i < n; i++ {
		times[i] = rnd.Int63n(maxVal) + 1
	}
	return buildInput(times)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	var tests []string
	// Deterministic tests
	tests = append(tests, buildInput([]int64{5}))
	tests = append(tests, buildInput([]int64{5, 3}))
	tests = append(tests, buildInput([]int64{4, 4}))
	tests = append(tests, buildInput([]int64{7, 2, 5, 1, 9}))
	tests = append(tests, buildInput([]int64{10, 10, 10, 10}))

	// Large deterministic case with many duplicates
	large := make([]int64, 100000)
	for i := range large {
		if i == 54321 {
			large[i] = 1
		} else {
			large[i] = 2
		}
	}
	tests = append(tests, buildInput(large))

	// Randomized tests
	const randomCases = 300
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomCases; i++ {
		n := rnd.Intn(1000) + 1
		maxVal := int64(1_000_000_000)
		tests = append(tests, randomTest(n, maxVal, rnd))
	}

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}

		if expOut != strings.TrimSpace(gotOut) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", idx+1, expOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
