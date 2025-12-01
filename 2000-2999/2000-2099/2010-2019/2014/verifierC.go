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
)

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2014C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
		}
		candOut, err := runProgramCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
		}
		if normalize(refOut) != normalize(candOut) {
			fail("wrong answer on test %d (%s)\nInput:\n%sExpected: %s\nGot: %s",
				idx+1, tc.name, tc.input, normalize(refOut), normalize(candOut))
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2014C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runProgramCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
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

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, cases []string) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(cases))
		for _, c := range cases {
			sb.WriteString(c)
		}
		tests = append(tests, testCase{name: name, input: sb.String()})
	}

	add("samples", []string{
		caseLine([]int{1, 2}),
		caseLine([]int{2, 19}),
		caseLine([]int{1, 3, 1, 3, 20}),
		caseLine([]int{1, 2, 3, 4}),
		caseLine([]int{1, 2, 3, 4, 5}),
		caseLine([]int{1, 2, 1, 1, 1, 2}),
	})

	add("small-enums", []string{
		intArrayCase([]int{2, 2, 2}),
		intArrayCase([]int{5, 5, 5}),
		intArrayCase([]int{10, 1, 1, 1, 1}),
	})

	rng := rand.New(rand.NewSource(2014))
	for i := 0; i < 5; i++ {
		tests = append(tests, randomCase(fmt.Sprintf("random-%d", i+1), rng, 5, 10, 30))
	}
	tests = append(tests, randomCase("large", rng, 10, 200000, 1000000))
	return tests
}

func caseLine(nums []int) string {
	n := len(nums)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func intArrayCase(nums []int) string {
	return caseLine(nums)
}

func randomCase(name string, rng *rand.Rand, t, maxN, maxA int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-3) + 3
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(maxA)+1)
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}
