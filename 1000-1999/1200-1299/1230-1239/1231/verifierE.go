package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseInt(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if candAns != refAns {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, refAns, candAns, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1231E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	refPath, err := referencePath()
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Base(refPath))
	cmd.Dir = filepath.Dir(refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller path")
	}
	return filepath.Abs(filepath.Join(filepath.Dir(thisFile), "1231E.go"))
}

func runBinary(path, input string) (string, error) {
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseInt(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("sample", []string{
			"3",
			"9",
			"iredppipe",
			"piedpiper",
			"4",
			"estt",
			"test",
			"4",
			"tste",
			"test",
		}),
		buildCase("single-char", []string{
			"3",
			"1",
			"a",
			"a",
			"1",
			"a",
			"b",
			"1",
			"b",
			"a",
		}),
		buildCase("simple", []string{
			"2",
			"4",
			"test",
			"test",
			"4",
			"test",
			"tset",
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1), rng.Intn(5)+1, rng.Intn(20)+1))
	}

	return tests
}

func buildCase(name string, blocks []string) testCase {
	var sb strings.Builder
	for i, block := range blocks {
		sb.WriteString(block)
		if i+1 < len(blocks) {
			sb.WriteByte('\n')
		}
	}
	return testCase{name: name, input: sb.String()}
}

func randomCase(rng *rand.Rand, name string, q, maxN int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		n := rng.Intn(maxN) + 1
		s := randomString(rng, n)
		t := randomPermutation(rng, s)
		if rng.Intn(5) == 0 {
			t = randomString(rng, n)
		}
		fmt.Fprintf(&sb, "%d\n%s\n%s\n", n, s, t)
	}
	return testCase{name: name, input: sb.String()}
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func randomPermutation(rng *rand.Rand, s string) string {
	runes := []byte(s)
	for i := range runes {
		j := rng.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
