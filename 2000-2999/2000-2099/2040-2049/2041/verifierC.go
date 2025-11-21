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
	refDir := filepath.Join("2000-2999", "2000-2099", "2040-2049", "2041")
	tmp, err := os.CreateTemp("", "ref2041C")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2041C.go")
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

func normalizeOutput(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func fixedTests() []testInput {
	tests := []testInput{
		{"2\n1 2\n3 4\n5 6\n7 8\n"},
		{"3\n1 2 3\n4 5 6\n7 8 9\n1 1 1\n2 2 2\n3 3 3\n"},
		{"4\n0 0 0 0\n1 1 1 1\n2 2 2 2\n3 3 3 3\n4 4 4 4\n5 5 5 5\n6 6 6 6\n7 7 7 7\n16 16 16 16\n"},
	}
	// Add one near-limit case with n=12.
	var sb strings.Builder
	n := 12
	sb.WriteString(fmt.Sprintf("%d\n", n))
	rng := rand.New(rand.NewSource(42))
	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			for z := 0; z < n; z++ {
				val := rng.Intn(20000001)
				if z > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", val))
			}
			sb.WriteByte('\n')
		}
	}
	tests = append(tests, testInput{text: sb.String()})
	return tests
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 18 {
		n := rng.Intn(6) + 2 // keep small for speed
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				for z := 0; z < n; z++ {
					val := rng.Intn(20000001)
					if z > 0 {
						sb.WriteByte(' ')
					}
					sb.WriteString(fmt.Sprintf("%d", val))
				}
				sb.WriteByte('\n')
			}
		}
		tests = append(tests, testInput{text: sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := randomTests()
	for idx, input := range tests {
		expect, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		got, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		if normalizeOutput(expect) != normalizeOutput(got) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input.text, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
