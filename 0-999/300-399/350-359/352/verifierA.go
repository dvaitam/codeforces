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

const refSource = "0-999/300-399/350-359/352/352A.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "352A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	src, err := filepath.Abs(refSource)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
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
	absBin, err := filepath.Abs(bin)
	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", absBin)
	} else {
		cmd = exec.Command(absBin)
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

	tests = append(tests,
		makeTest([]int{0}),
		makeTest([]int{5}),
		makeTest([]int{0, 0, 0}),
		makeTest([]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 0}),
		makeTest([]int{5, 5, 5, 0, 0}),
		makeTest([]int{0, 5, 0, 5, 0, 5}),
	)

	// deterministic boundary tests
	tests = append(tests, makeCountsTest(1000, 0))
	tests = append(tests, makeCountsTest(0, 1000))
	tests = append(tests, makeCountsTest(992, 8))
	tests = append(tests, makeCountsTest(900, 100))

	for i := 0; i < 300; i++ {
		n := rng.Intn(1000) + 1
		arr := make([]int, n)
		for j := range arr {
			if rng.Intn(2) == 0 {
				arr[j] = 0
			} else {
				arr[j] = 5
			}
		}
		tests = append(tests, makeTest(arr))
	}

	// cases with exactly multiple-of-9 fives and at least one zero
	for m := 1; m <= 5; m++ {
		fives := 9 * m
		zeros := m
		tests = append(tests, makeCountsTest(zeros, fives))
	}

	return tests
}

func makeCountsTest(cnt0, cnt5 int) testCase {
	arr := make([]int, 0, cnt0+cnt5)
	for i := 0; i < cnt0; i++ {
		arr = append(arr, 0)
	}
	for i := 0; i < cnt5; i++ {
		arr = append(arr, 5)
	}
	return makeTest(arr)
}

func makeTest(numbers []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(numbers))
	for i, v := range numbers {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String()}
}
