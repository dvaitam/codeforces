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
	refDir := filepath.Join("2000-2999", "2100-2199", "2160-2169", "2166")
	tmp, err := os.CreateTemp("", "ref2166A")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2166A.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v\n%s", err, string(out))
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

func fixedTests() []string {
	return []string{
		"5\n3\nqwq\n2\naa\n4\ntest\n5\nabbac\n6\nabcabc\n",
		"2\n2\nab\n2\nba\n",
		"1\n100\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n",
		"1\n100\nabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmn\n",
	}
}

func randomString(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func randomTest(rng *rand.Rand) string {
	t := rng.Intn(20) + 1
	remaining := 100
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		casesLeft := t - i
		minLen := 2
		maxLen := remaining - (casesLeft-1)*2
		if maxLen < minLen {
			maxLen = minLen
		}
		length := minLen
		if maxLen > minLen {
			length += rng.Intn(maxLen - minLen + 1)
		}
		remaining -= length
		sb.WriteString(fmt.Sprintf("%d\n", length))
		sb.WriteString(randomString(rng, length))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateTests() []string {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		expect, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if normalizeOutput(expect) != normalizeOutput(got) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
