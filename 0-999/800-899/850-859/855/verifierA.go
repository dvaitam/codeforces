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

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"855A.go",
		filepath.Join("0-999", "800-899", "850-859", "855", "855A.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 855A.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref855A_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
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
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func normalize(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
		lines[i] = strings.ToUpper(lines[i])
	}
	return strings.Join(lines, "\n")
}

func generateTests() []testCase {
	tests := []testCase{
		newTest([]string{"a"}),
		newTest([]string{"a", "a"}),
		newTest([]string{"abc", "def", "abc", "ghi", "def"}),
		newTest([]string{"same", "same", "same", "same"}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		n := rng.Intn(100) + 1
		names := make([]string, n)
		for i := 0; i < n; i++ {
			names[i] = randomName(rng)
		}
		tests = append(tests, newTest(names))
	}
	return tests
}

func randomName(rng *rand.Rand) string {
	length := rng.Intn(100) + 1
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteByte(byte('a' + rng.Intn(26)))
	}
	return b.String()
}

func newTest(names []string) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(names)))
	for _, name := range names {
		b.WriteString(name)
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}
