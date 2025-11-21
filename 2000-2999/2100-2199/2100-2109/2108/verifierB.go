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

const refSource = "2000-2999/2100-2199/2100-2109/2108/2108B.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-2108B-*")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func parseOutput(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	return fields, nil
}

func fixedTests() []string {
	return []string{
		"5\n1\n8\n10\n100\n999999999\n",
		"4\n0\n2\n3\n7\n",
	}
}

func randomInput(rng *rand.Rand) string {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		x := rng.Int63n(1_000_000) // keep moderate for speed
		sb.WriteString(fmt.Sprintf("%d\n", x))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refPath, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// a couple of larger x values to check higher search region
	tests = append(tests, "3\n999999937\n999999991\n1000000000\n")
	for i := 0; i < 20; i++ {
		tests = append(tests, randomInput(rng))
	}

	for idx, input := range tests {
		var t int
		if _, err := fmt.Fscan(strings.NewReader(input), &t); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse test count for case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseOutput(expOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse failed on case %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse failed on case %d: %v\noutput:\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}

		for i := 0; i < t; i++ {
			if exp[i] != got[i] {
				fmt.Fprintf(os.Stderr, "mismatch on case %d test %d: expected %s got %s\ninput:\n%s", idx+1, i+1, exp[i], got[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed")
}
