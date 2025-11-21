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
	refDir := filepath.Join("2000-2999", "2100-2199", "2160-2169", "2161")
	tmp, err := os.CreateTemp("", "ref2161G")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2161G.go")
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

type testInput struct {
	text string
}

func fixedTests() []testInput {
	return []testInput{
		{"5 4\n6 4 7 5 4\n0\n2\n4\n6\n"},
		{"2 3\n0 0\n0\n1\n3\n"},
		{"3 5\n1 2 3\n0\n1\n2\n3\n4\n"},
		{"4 4\n1048575 0 524287 262143\n0\n1048575\n1\n262143\n"},
	}
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1 << 20)
	}
	return arr
}

func randomQueries(rng *rand.Rand, q int) []int {
	qs := make([]int, q)
	for i := 0; i < q; i++ {
		qs[i] = rng.Intn(1 << 20)
	}
	return qs
}

func buildTest(n, q int, arr []int, qs []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, x := range qs {
		sb.WriteString(fmt.Sprintf("%d\n", x))
	}
	return sb.String()
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(2000-2+1) + 2
		q := rng.Intn(2000) + 1
		arr := randomArray(rng, n)
		qs := randomQueries(rng, q)
		tests = append(tests, testInput{text: buildTest(n, q, arr, qs)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
