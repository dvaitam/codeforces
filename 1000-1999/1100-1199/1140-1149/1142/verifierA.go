package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
	desc  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		expOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, expOut)
			os.Exit(1)
		}
		expX, expY, err := parseAnswer(expOut)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, expOut)
			os.Exit(1)
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		x, y, err := parseAnswer(out)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		if x != expX || y != expY {
			fmt.Printf("Wrong answer on test %d (%s): expected (%d %d) got (%d %d)\nInput:\n%s", i+1, tc.desc, expX, expY, x, y, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1142A.bin"
	cmd := exec.Command("go", "build", "-o", path, "1142A.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseAnswer(out string) (int64, int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var x, y int64
	if _, err := fmt.Fscan(reader, &x, &y); err != nil {
		return 0, 0, fmt.Errorf("failed to read two integers: %v", err)
	}
	if extra := strings.TrimSpace(readRemaining(reader)); extra != "" {
		return 0, 0, fmt.Errorf("unexpected extra output: %q", extra)
	}
	return x, y, nil
}

func readRemaining(r *bufio.Reader) string {
	var sb strings.Builder
	for {
		line, err := r.ReadString('\n')
		sb.WriteString(line)
		if err != nil {
			break
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, n, k, a, b int) {
		input := fmt.Sprintf("%d %d\n%d %d\n", n, k, a, b)
		tests = append(tests, testCase{desc: desc, input: input})
	}

	// Manual edge cases
	add("tiny-1", 1, 1, 0, 0)
	add("tiny-2", 2, 1, 0, 0)
	add("max-a", 3, 10, 5, 5)
	add("diff-a-b", 4, 7, 1, 3)
	add("only-one-restaurant", 1, 100000, 0, 50000)

	rng := rand.New(rand.NewSource(123456789))
	// Random moderate tests
	for len(tests) < 120 {
		n := rng.Intn(100000) + 1
		k := rng.Intn(100000) + 1
		maxAB := k / 2
		a := 0
		b := 0
		if maxAB > 0 {
			a = rng.Intn(maxAB + 1)
			b = rng.Intn(maxAB + 1)
		}
		desc := fmt.Sprintf("random-%d", len(tests))
		add(desc, n, k, a, b)
	}

	// Stress edges at bounds
	add("max-values-0", 100000, 100000, 0, 0)
	add("max-values-mid", 100000, 100000, 50000, 12345)

	return tests
}
