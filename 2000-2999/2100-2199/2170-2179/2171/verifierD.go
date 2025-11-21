package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	g int
	c int
	l int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := buildTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d %d\n", tc.g, tc.c, tc.l)
		exp, err := runProgram(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if normalize(exp) != normalize(got) {
			fmt.Fprintf(os.Stderr, "test %d mismatch:\ninput: %s\nexpected: %s\ngot: %s", idx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2171D.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2171D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
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

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func buildTests() []testCase {
	tests := []testCase{
		{80, 80, 80},
		{80, 80, 89},
		{80, 80, 90},
		{80, 90, 100},
		{85, 86, 95},
		{95, 86, 85},
		{88, 94, 95},
		{100, 80, 81},
		{98, 99, 98},
		{100, 100, 100},
		{90, 99, 98},
		{100, 90, 89},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 250 {
		tests = append(tests, testCase{
			g: 80 + rng.Intn(21),
			c: 80 + rng.Intn(21),
			l: 80 + rng.Intn(21),
		})
	}
	return tests
}
