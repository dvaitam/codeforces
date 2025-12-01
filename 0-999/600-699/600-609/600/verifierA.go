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

const refSource = "./600A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expLines, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotLines, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}

		if gotLines[0] != expLines[0] || gotLines[1] != expLines[1] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s)\nexpected:\n%s\n%s\ngot:\n%s\n%s\ninput:\n%s\n",
				idx+1, tc.name, expLines[0], expLines[1], gotLines[0], gotLines[1], tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "600A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseOutput(out string) ([2]string, error) {
	var lines [2]string
	all := strings.Split(out, "\n")
	i := 0
	for _, line := range all {
		if i >= 2 {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" && len(all) < 2 {
			continue
		}
		lines[i] = line
		i++
	}
	if i < 2 {
		return lines, fmt.Errorf("expected two lines, got %d", i)
	}
	return lines, nil
}

func generateTests() []testCase {
	tests := []testCase{
		makeCase("simple-mix", "aba,123;1a;0"),
		makeCase("only-separators", ";;"),
		makeCase("leading-zero", "01;1.0;0;00"),
		makeCase("single-number", "987654321"),
		makeCase("only-letters", "abcde"),
		makeCase("dash-and-dot", "-.-;0,-1.;"),
		makeCase("empty-between-commas", ",,"),
		makeCase("upper-lower", "A1;B2,C3"),
		makeCase("ends-separator", ";123;"),
		makeCase("long-digits", strings.Repeat("9", 50)+",0"),
	}

	tests = append(tests, longRandomCase("max-length-random", 100000))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1)))
	}
	return tests
}

func makeCase(name, s string) testCase {
	return testCase{name: name, input: s + "\n"}
}

func longRandomCase(name string, length int) testCase {
	rng := rand.New(rand.NewSource(6000600))
	return randomCaseWithLength(rng, name, length)
}

var charset = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.,;")

func randomCase(rng *rand.Rand, name string) testCase {
	length := rng.Intn(200) + 1
	return randomCaseWithLength(rng, name, length)
}

func randomCaseWithLength(rng *rand.Rand, name string, length int) testCase {
	if length <= 0 {
		length = 1
	}
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rng.Intn(len(charset))])
	}
	// ensure no leading/trailing separators to check handling? not necessary
	return makeCase(name, sb.String())
}
