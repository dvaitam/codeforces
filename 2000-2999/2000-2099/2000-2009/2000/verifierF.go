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

const referenceSource = "./2000F.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := manualTests()
	for len(tests) < 150 {
		tests = append(tests, randomTest(rng))
	}

	for idx, input := range tests {
		refOutRaw, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOutRaw)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOutRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOutRaw)
			os.Exit(1)
		}

		candOutRaw, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOutRaw)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOutRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOutRaw)
			os.Exit(1)
		}
		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d lines got %d\ninput:\n%sreference:\n%s\ncandidate:\n%s\n",
				idx+1, len(refVals), len(candVals), input, refOutRaw, candOutRaw)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, i+1, refVals[i], candVals[i], input, refOutRaw, candOutRaw)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2000F-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2000F.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
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
	return stdout.String(), nil
}

func manualTests() []string {
	return []string{
		// sample from statement
		"7\n1 4\n6 3\n1 5\n4 4\n5 10\n1 1\n1 1\n1 1\n1 1\n2 100\n1 2\n5 6\n3 11\n2 2\n3 3\n4 4\n3 25\n9 2\n4 3\n8 10\n4 18\n5 4\n8 5\n8 3\n6 2\n",
		"1\n1 1\n1 1\n",
		"1\n2 2\n1 1\n1 1\n",
		"2\n2 3\n1 2\n2 3\n1 1\n5 5\n",
	}
}

func randomTest(rng *rand.Rand) string {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		k := rng.Intn(8) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j := 0; j < n; j++ {
			a := rng.Intn(8) + 1
			b := rng.Intn(8) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		}
	}
	return sb.String()
}

func parseOutputs(out string) ([]int, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, fmt.Errorf("empty output")
	}
	lines := strings.Split(out, "\n")
	values := make([]int, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			return nil, fmt.Errorf("empty line at %d", i+1)
		}
		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q on line %d", line, i+1)
		}
		values[i] = val
	}
	return values, nil
}
