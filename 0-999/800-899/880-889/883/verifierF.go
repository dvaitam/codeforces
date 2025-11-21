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
	words []string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-883F-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", path, "883F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return path, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.words)))
	for _, w := range tc.words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer output, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 0 {
		return 0, fmt.Errorf("negative group count %d", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{words: []string{"a", "a"}},
		{words: []string{"mihail", "mikhail", "khun", "kkkhoon"}},
		{words: []string{"u", "oo", "uu"}},
		{words: []string{"h", "kh", "kkh", "kkkh"}},
		{words: []string{"aaaaaaaaaaaaaaaaaaaa"}},
	}
}

func randomWord(rng *rand.Rand) string {
	length := rng.Intn(20) + 1
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = byte('a' + rng.Intn(26))
	}
	return string(buf)
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(400) + 1
	words := make([]string, n)
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			// bias towards sequences of k and o to stress reductions
			choices := []string{"k", "h", "o", "u"}
			length := rng.Intn(12) + 1
			var sb strings.Builder
			for j := 0; j < length; j++ {
				sb.WriteString(choices[rng.Intn(len(choices))])
			}
			words[i] = sb.String()
		} else {
			words[i] = randomWord(rng)
		}
	}
	return testCase{words: words}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expStr)
			os.Exit(1)
		}

		gotStr, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotStr, input)
			os.Exit(1)
		}

		if got != exp {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
