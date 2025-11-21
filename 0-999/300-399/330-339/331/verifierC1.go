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

type testInput struct {
	text string
}

func buildReference() (string, error) {
	dir := filepath.Join("0-999", "300-399", "330-339", "331")
	tmp, err := os.CreateTemp("", "ref331C1")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "331C1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func fixedTests() []testInput {
	return []testInput{
		{text: "0\n"},
		{text: "1\n"},
		{text: "5\n"},
		{text: "24\n"},
		{text: "1000000\n"},
	}
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		n := rng.Intn(1_000_000 + 1) // n <= 1e6 in C1
		tests = append(tests, testInput{text: fmt.Sprintf("%d\n", n)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := randomTests()
	for idx, test := range tests {
		expect, err := runBinary(refBin, test.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput: %s\n", idx+1, err, strings.TrimSpace(test.text))
			os.Exit(1)
		}
		got, err := runBinary(candidate, test.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput: %s\n", idx+1, err, strings.TrimSpace(test.text))
			os.Exit(1)
		}
		if normalize(expect) != normalize(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\ninput: %s\nexpected: %s\ngot: %s\n", idx+1, strings.TrimSpace(test.text), normalize(expect), normalize(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
