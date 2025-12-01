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

const refSource = "256D.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		if normalize(got) != normalize(expect) {
			fail("mismatch on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s", i+1, tc.input, expect, got)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "256D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	powers := []int{1, 2, 4, 8, 16}
	// enumerate all small cases
	for _, n := range []int{1, 2, 4} {
		for k := 1; k <= n; k++ {
			tests = append(tests, makeTest(n, k))
		}
	}
	// stress edges
	tests = append(tests, makeTest(16, 1))
	tests = append(tests, makeTest(16, 16))
	tests = append(tests, makeTest(16, 8))
	tests = append(tests, makeTest(8, 3))
	tests = append(tests, makeTest(8, 7))
	tests = append(tests, makeTest(4, 4))
	// random cases
	for i := 0; i < 200; i++ {
		n := powers[rng.Intn(len(powers))]
		k := rng.Intn(n) + 1
		tests = append(tests, makeTest(n, k))
	}
	return tests
}

func makeTest(n, k int) testCase {
	return testCase{input: fmt.Sprintf("%d %d\n", n, k)}
}
