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
	"time"
)

type testInput struct {
	text     string
	ansCount int
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2000-2099", "2000-2009", "2002")
	tmp, err := os.CreateTemp("", "ref2002B")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2002B.go")
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

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(lines))
	}
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return lines, nil
}

func fixedTests() []testInput {
	return []testInput{
		{text: "2\n1\n1\n1\n2\n1 2\n2 1\n", ansCount: 2},
		{text: "1\n3\n1 2 3\n3 2 1\n", ansCount: 1},
	}
}

func randomPermutation(n int, rng *rand.Rand) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { p[i], p[j] = p[j], p[i] })
	return p
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(8) + 1
		if rng.Intn(5) == 0 {
			n = rng.Intn(50) + 1
		}
		a := randomPermutation(n, rng)
		b := randomPermutation(n, rng)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testInput{text: sb.String(), ansCount: t}
}

func largeStructuredTest() []testInput {
	n := 200000
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i + 1
		b[i] = n - i
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return []testInput{
		{text: sb.String(), ansCount: 1},
	}
}

func generateTests() []testInput {
	tests := fixedTests()
	tests = append(tests, largeStructuredTest()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		expectRaw, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput excerpt:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		expect, err := parseOutput(expectRaw, input.ansCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, expectRaw)
			os.Exit(1)
		}

		gotRaw, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput excerpt:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		got, err := parseOutput(gotRaw, input.ansCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, gotRaw)
			os.Exit(1)
		}

		for i := range expect {
			if expect[i] != got[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d case %d: expected %s got %s\ninput excerpt:\n%s\n", idx+1, i+1, expect[i], got[i], preview(input.text))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func preview(s string) string {
	if len(s) <= 400 {
		return s
	}
	return s[:400] + "...\n"
}
