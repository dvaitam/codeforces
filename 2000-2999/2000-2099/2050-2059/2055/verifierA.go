package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	a int
	b int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %s got %s\n", i+1, expected[i], got[i])
			fmt.Fprintf(os.Stderr, "n=%d a=%d b=%d\n", tests[i].n, tests[i].a, tests[i].b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2055A.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2055A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, cases int) ([]string, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	lines := strings.Fields(out)
	if len(lines) != cases {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", cases, len(lines), out)
	}
	for i := range lines {
		lines[i] = strings.ToUpper(lines[i])
	}
	return lines, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 500 {
		n := rng.Intn(99) + 2
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			b = (b % n) + 1
		}
		tests = append(tests, testCase{n: n, a: a, b: b})
	}
	return tests
}

func deterministicTests() []testCase {
	tests := []testCase{
		{n: 2, a: 1, b: 2},
		{n: 3, a: 3, b: 1},
		{n: 4, a: 2, b: 3},
		{n: 5, a: 2, b: 4},
		{n: 7, a: 6, b: 2},
	}
	for n := 2; n <= 20; n++ {
		for a := 1; a <= n; a++ {
			for b := 1; b <= n; b++ {
				if a != b {
					tests = append(tests, testCase{n: n, a: a, b: b})
				}
			}
		}
	}
	return tests
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.a, tc.b))
	}
	return sb.String()
}
