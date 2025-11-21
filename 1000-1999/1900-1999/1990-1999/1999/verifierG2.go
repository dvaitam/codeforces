package main

import (
	"bufio"
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

const refSource = "1000-1999/1900-1999/1990-1999/1999/1999G2.go"

type testCase struct {
	input string
	info  testInfo
}

type testInfo struct {
	answers []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests, err := buildTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build tests:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refOut, err := runExecutable(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(refOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(candOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\nInput:\n%sCandidate output:\n%s\n", i+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1999G2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutable(path, input string) (string, error) {
	cmd := exec.Command(path)
	return execute(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return execute(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func execute(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildTests() ([]testCase, error) {
	base := []string{
		"1\n2\n",
		"3\n2\n500\n999\n",
		"5\n3\n4\n5\n6\n7\n",
	}
	var tests []testCase
	for _, input := range base {
		info, err := prepareTest(input)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: input, info: info})
	}

	randomConfigs := []struct {
		t    int
		seed int64
	}{
		{10, 1},
		{25, 2},
		{50, 3},
		{75, 4},
		{100, 5},
		{200, 6},
		{500, 7},
		{1000, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		input := randomTest(cfg.t, cfg.seed)
		info, err := prepareTest(input)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: input, info: info})
	}
	return tests, nil
}

func randomTest(t int, seed int64) string {
	if t < 1 {
		t = 1
	}
	if t > 1000 {
		t = 1000
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		val := r.Intn(998) + 2
		sb.WriteString(fmt.Sprintf("%d\n", val))
	}
	return sb.String()
}

func prepareTest(input string) (testInfo, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return testInfo{}, fmt.Errorf("failed to read t: %v", err)
	}
	if t < 1 || t > 1000 {
		return testInfo{}, fmt.Errorf("invalid t %d", t)
	}
	ans := make([]int, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &ans[i]); err != nil {
			return testInfo{}, fmt.Errorf("failed to read x[%d]: %v", i, err)
		}
		if ans[i] < 2 || ans[i] > 999 {
			return testInfo{}, fmt.Errorf("x[%d]=%d out of range", i, ans[i])
		}
	}
	return testInfo{answers: ans}, nil
}

func checkOutput(output string, info testInfo) error {
	tokens := strings.Fields(output)
	if len(tokens) != len(info.answers) {
		return fmt.Errorf("expected %d answers, got %d", len(info.answers), len(tokens))
	}
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("answer %d is not an integer: %q", i+1, tok)
		}
		if val != info.answers[i] {
			return fmt.Errorf("answer %d mismatch: expected %d got %d", i+1, info.answers[i], val)
		}
	}
	return nil
}
