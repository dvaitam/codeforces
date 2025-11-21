package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2040-2049/2049/2049B.go"

type testCase struct {
	n int
	s string
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierB.go /path/to/candidate")
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInput(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refAnswers, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAnswers, err := parseAnswersWithWitness(candOut, tests)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	for i, tc := range tests {
		ref := refAnswers[i]
		got := strings.ToUpper(candAnswers[i].verdict)
		if got == "YES" {
			if ref != "YES" {
				fail("test %d: candidate says YES but reference says NO", i+1)
			}
			if err := checkWitness(tc, candAnswers[i].perm); err != nil {
				fail("test %d: %v", i+1, err)
			}
		} else if got == "NO" {
			if ref == "YES" {
				fail("test %d: candidate says NO but reference says YES", i+1)
			}
		} else {
			fail("test %d: invalid verdict %q", i+1, candAnswers[i].verdict)
		}
	}

	fmt.Println("OK")
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		var s string
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, err
		}
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return nil, err
		}
		tests[i] = testCase{n: n, s: s}
	}
	return tests, nil
}

func parseAnswers(out string, tests int) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]string, tests)
	for i := 0; i < tests; i++ {
		token, err := readToken(reader)
		if err != nil {
			return nil, fmt.Errorf("missing verdict for test %d: %v", i+1, err)
		}
		ans[i] = strings.ToUpper(token)
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("extra token %q after outputs", extra)
		}
		return nil, err
	}
	return ans, nil
}

type candidateAnswer struct {
	verdict string
	perm    []int
}

func parseAnswersWithWitness(out string, tests []testCase) ([]candidateAnswer, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]candidateAnswer, len(tests))
	for i, tc := range tests {
		token, err := readToken(reader)
		if err != nil {
			return nil, fmt.Errorf("missing verdict for test %d: %v", i+1, err)
		}
		ans[i].verdict = token
		if strings.ToUpper(token) == "YES" {
			perm := make([]int, tc.n)
			for j := 0; j < tc.n; j++ {
				tok, err := readToken(reader)
				if err != nil {
					return nil, fmt.Errorf("missing permutation value %d for test %d", j+1, i+1)
				}
				val, err := strconv.Atoi(tok)
				if err != nil {
					return nil, fmt.Errorf("invalid integer %q in permutation for test %d", tok, i+1)
				}
				perm[j] = val
			}
			ans[i].perm = perm
		}
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("extra token %q after outputs", extra)
		}
		return nil, err
	}
	return ans, nil
}

func checkWitness(tc testCase, perm []int) error {
	if len(perm) != tc.n {
		return fmt.Errorf("permutation length mismatch: expected %d got %d", tc.n, len(perm))
	}
	seen := make([]bool, tc.n+1)
	for idx, v := range perm {
		if v < 1 || v > tc.n {
			return fmt.Errorf("value %d at position %d out of range [1,%d]", v, idx+1, tc.n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
		if tc.s[idx] == 'p' {
			if err := checkPermutationPrefix(perm[:idx+1]); err != nil {
				return fmt.Errorf("prefix constraint at position %d violated: %v", idx+1, err)
			}
		}
		if tc.s[idx] == 's' {
			if err := checkPermutationPrefix(perm[idx:]); err != nil {
				return fmt.Errorf("suffix constraint at position %d violated: %v", idx+1, err)
			}
		}
	}
	return nil
}

func checkPermutationPrefix(segment []int) error {
	length := len(segment)
	if length == 0 {
		return fmt.Errorf("empty segment")
	}
	seen := make([]bool, length+1)
	for _, v := range segment {
		if v < 1 || v > length {
			return fmt.Errorf("value %d not in range [1,%d]", v, length)
		}
		if seen[v] {
			return fmt.Errorf("value %d repeats", v)
		}
		seen[v] = true
	}
	return nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2049B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
}

func readToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}
		if !isSpace(ch) {
			sb.WriteByte(ch)
			break
		}
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if isSpace(ch) {
			break
		}
		sb.WriteByte(ch)
	}
	return sb.String(), nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t' || b == '\v' || b == '\f'
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
