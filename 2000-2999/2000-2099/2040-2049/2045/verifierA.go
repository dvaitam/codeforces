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

const refSource = "2000-2999/2000-2099/2040-2049/2045/2045A.go"

type testCase struct {
	name  string
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
		if strings.TrimSpace(refOut) != strings.TrimSpace(candOut) {
			fail("wrong answer on test %d (%s)\nInput:\n%sExpected: %s\nGot: %s", idx+1, tc.name, tc.input, strings.TrimSpace(refOut), strings.TrimSpace(candOut))
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2045A-ref-*")
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

func buildTests() []testCase {
	var tests []testCase
	add := func(name, s string) {
		tests = append(tests, testCase{name: name, input: s + "\n"})
	}

	add("sample1", "ICPCJAKARTA")
	add("sample2", "NGENG")
	add("sample3", "YYY")
	add("sample4", "DANGAN")
	add("sample5", "AEIOUY")
	add("minimal", "A")
	add("allY", strings.Repeat("Y", 10))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		add(fmt.Sprintf("random-%d", i+1), randomString(rng, 50))
	}
	add("max-length", randomString(rng, 5000))

	return tests
}

func randomString(rng *rand.Rand, length int) string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
	}
	return sb.String()
}
