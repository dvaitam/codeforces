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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
	const refName = "./ref_2172A.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2172A.go")
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
	return strings.ToLower(strings.TrimSpace(s))
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 21*21*21)
	for g := 80; g <= 100; g++ {
		for c := 80; c <= 100; c++ {
			for l := 80; l <= 100; l++ {
				tests = append(tests, testCase{g: g, c: c, l: l})
			}
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, testCase{
			g: 80 + rng.Intn(21),
			c: 80 + rng.Intn(21),
			l: 80 + rng.Intn(21),
		})
	}
	return tests
}
