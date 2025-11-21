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

func buildReference() (string, error) {
	dir := filepath.Join("0-999", "500-599", "540-549", "542")
	src := "542D.go"
	tmp, err := os.CreateTemp("", "ref542D")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, src)
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

type testInput struct {
	text string
}

func fixedTests() []testInput {
	return []testInput{
		{text: "1\n"},
		{text: "3\n"},
		{text: "24\n"},
		{text: "1000000000000\n"},
	}
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		val := rng.Int63n(1_000_000_000_000) + 1 // up to 1e12
		tests = append(tests, testInput{text: fmt.Sprintf("%d\n", val)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput: %s\n", idx+1, err, test.text)
			os.Exit(1)
		}
		got, err := runBinary(candidate, test.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput: %s\n", idx+1, err, test.text)
			os.Exit(1)
		}
		if normalize(expect) != normalize(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\ninput: %s\nexpected: %s\ngot: %s\n", idx+1, strings.TrimSpace(test.text), normalize(expect), normalize(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
