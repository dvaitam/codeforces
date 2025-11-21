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
	s string
	n int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-335A-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", outPath, "335A.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
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
	return stdout.String(), nil
}

func countFreq(str string) [26]int {
	var freq [26]int
	for _, ch := range str {
		freq[ch-'a']++
	}
	return freq
}

func maxSheets(freq [26]int, n int) int {
	distinct := 0
	for _, f := range freq {
		if f > 0 {
			distinct++
		}
	}
	if distinct > n {
		return -1
	}
	sheets := 1
	for {
		required := 0
		for _, f := range freq {
			if f > 0 {
				required += (f + sheets - 1) / sheets
			}
		}
		if required <= n {
			return sheets
		}
		sheets++
	}
}

func parseOutput(output string) (int, string, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return 0, "", fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, "", fmt.Errorf("first line not int: %v", err)
	}
	if k == -1 {
		return k, "", nil
	}
	if len(lines) < 2 {
		return 0, "", fmt.Errorf("missing second line for sheet string")
	}
	sheet := strings.TrimSpace(lines[1])
	return k, sheet, nil
}

func validate(tc testCase, k int, sheet string) error {
	freq := countFreq(tc.s)
	minSheets := maxSheets(freq, tc.n)
	if minSheets == -1 {
		if k != -1 {
			return fmt.Errorf("expected -1 but got %d", k)
		}
		return nil
	}
	if k == -1 {
		return fmt.Errorf("solution exists but returned -1")
	}
	if k != minSheets {
		return fmt.Errorf("expected %d sheets, got %d", minSheets, k)
	}
	if len(sheet) != tc.n {
		return fmt.Errorf("sheet length mismatch: expected %d, got %d", tc.n, len(sheet))
	}
	sheetFreq := countFreq(sheet)
	for i := 0; i < 26; i++ {
		if freq[i] > 0 {
			required := (freq[i] + k - 1) / k
			if sheetFreq[i] < required {
				return fmt.Errorf("character %c: need >= %d per sheet, got %d", 'a'+i, required, sheetFreq[i])
			}
		}
	}
	return nil
}

func generateTests() []testCase {
	tests := []testCase{
		{s: "banana", n: 4},
		{s: "banana", n: 3},
		{s: "banana", n: 2},
		{s: "a", n: 1},
		{s: "aaaa", n: 2},
		{s: "abc", n: 3},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		lenS := rng.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < lenS; j++ {
			sb.WriteByte(byte('a' + rng.Intn(3)))
		}
		n := rng.Intn(10) + 1
		tests = append(tests, testCase{s: sb.String(), n: n})
	}
	for i := 0; i < 10; i++ {
		lenS := rng.Intn(1000) + 1
		var sb strings.Builder
		for j := 0; j < lenS; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		n := rng.Intn(1000) + 1
		tests = append(tests, testCase{s: sb.String(), n: n})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%s\n%d\n", tc.s, tc.n)
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expK, _, err := parseOutput(expected)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expected)
			os.Exit(1)
		}

		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		actK, sheet, err := parseOutput(actual)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n%s", idx+1, err, actual)
			os.Exit(1)
		}
		if actK != expK {
			fmt.Fprintf(os.Stderr, "test %d: expected %d sheets, got %d\ninput: %s", idx+1, expK, actK, input)
			os.Exit(1)
		}
		if err := validate(tc, actK, sheet); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, actual)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
