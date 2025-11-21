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
	refDir := filepath.Join("1000-1999", "1900-1999", "1910-1919", "1912")
	tmp, err := os.CreateTemp("", "ref1912D")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "1912D.go")
	cmd.Dir = refDir
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

func fixedTests() []testInput {
	return []testInput{
		{text: "3\n10 3\n10 11\n10 4\n"},
		{text: "2\n8 2\n7 5\n"},
	}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		b := rng.Intn(2000) + 2
		n := rng.Intn(2000) + 2
		sb.WriteString(fmt.Sprintf("%d %d\n", b, n))
	}
	return testInput{text: sb.String()}
}

func edgeTests() []testInput {
	cases := []struct {
		b int
		n int
	}{
		{2, 2},
		{1000000, 1000000},
		{999983, 999979},
		{1000000, 2},
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.b, c.n))
	}
	return []testInput{
		{text: sb.String()},
	}
}

func generateTests() []testInput {
	tests := fixedTests()
	tests = append(tests, edgeTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, input := range tests {
		expect, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		got, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		if normalize(expect) != normalize(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input.text, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func normalize(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func preview(s string) string {
	if len(s) <= 400 {
		return s
	}
	return s[:400] + "...\n"
}
