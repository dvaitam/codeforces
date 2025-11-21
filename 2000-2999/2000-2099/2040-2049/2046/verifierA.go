package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceA = "2000-2999/2000-2099/2040-2049/2046/2046A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		userOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if normalize(userOut) != normalize(refOut) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, refOut, userOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2046A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(source))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
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

func generateTests() []string {
	var tests []string
	tests = append(tests, sampleInput())
	tests = append(tests, singleColumnCases())
	tests = append(tests, randomCase(5, 10, 42))
	tests = append(tests, randomCase(20, 100, 99))
	tests = append(tests, randomCase(1000, 2000, 1234))
	tests = append(tests, randomCase(3000, 5000, 98765))
	return tests
}

func sampleInput() string {
	return strings.TrimSpace(`3
1
-10
5
3
1 2 3
10 -5 -3
4
2 8 5 3
1 10 3 4
`) + "\n"
}

func singleColumnCases() string {
	var sb strings.Builder
	sb.WriteString("4\n")
	sb.WriteString("1\n0\n0\n")
	sb.WriteString("1\n100000\n-100000\n")
	sb.WriteString("1\n-100000\n100000\n")
	sb.WriteString("1\n12345\n67890\n")
	return sb.String()
}

func randomCase(t, maxSumN int, seed int64) string {
	rng := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	remaining := maxSumN
	for i := 0; i < t; i++ {
		maxN := remaining - (t - i - 1)
		if maxN < 1 {
			maxN = 1
		}
		n := 1 + rng.Intn(maxN)
		remaining -= n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			val := rng.Intn(200001) - 100000
			sb.WriteString(fmt.Sprintf("%d", val))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			val := rng.Intn(200001) - 100000
			sb.WriteString(fmt.Sprintf("%d", val))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
