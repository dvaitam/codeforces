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

const refSource = "2000-2999/2000-2099/2030-2039/2034/2034E.go"

type testCase struct {
	n int
	k int
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierE.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInput(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), inputData)
	if err != nil {
		fail("reference solution failed: %v", err)
	}
	refVerdicts, err := parseOutput(refOut, tests, true)
	if err != nil {
		fail("invalid reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candVerdicts, err := parseOutput(candOut, tests, true)
	if err != nil {
		fail("invalid candidate output: %v", err)
	}

	for i := range tests {
		if refVerdicts[i] == "yes" && candVerdicts[i] == "no" {
			fail("test %d: reference found a solution but candidate printed NO", i+1)
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
		if _, err := fmt.Fscan(reader, &tests[i].n, &tests[i].k); err != nil {
			return nil, err
		}
	}
	return tests, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2034E-ref-*")
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

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
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

func parseOutput(out string, tests []testCase, validate bool) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	verdicts := make([]string, len(tests))
	for idx, tc := range tests {
		token, err := readToken(reader)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("not enough outputs for test %d", idx+1)
			}
			return nil, err
		}
		ver := strings.ToLower(token)
		if ver != "yes" && ver != "no" {
			return nil, fmt.Errorf("test %d: expected YES/NO, got %q", idx+1, token)
		}
		verdicts[idx] = ver
		if ver == "no" {
			continue
		}
		if !validate {
			// Still need to consume k*n integers to align parsing.
			for i := 0; i < tc.k*tc.n; i++ {
				if _, err := readToken(reader); err != nil {
					if err == io.EOF {
						return nil, fmt.Errorf("test %d: incomplete permutation data", idx+1)
					}
					return nil, err
				}
			}
			continue
		}

		perms := make([][]int, tc.k)
		sum := make([]int, tc.n)
		seenPerms := make(map[string]struct{}, tc.k)
		for i := 0; i < tc.k; i++ {
			p := make([]int, tc.n)
			for j := 0; j < tc.n; j++ {
				tok, err := readToken(reader)
				if err != nil {
					if err == io.EOF {
						return nil, fmt.Errorf("test %d: expected %d numbers for permutations, got %d", idx+1, tc.k*tc.n, i*tc.n+j)
					}
					return nil, err
				}
				val, convErr := strconv.Atoi(tok)
				if convErr != nil {
					return nil, fmt.Errorf("test %d: invalid integer %q", idx+1, tok)
				}
				p[j] = val
				sum[j] += val
			}
			if err := validatePermutation(p, tc.n); err != nil {
				return nil, fmt.Errorf("test %d: %v", idx+1, err)
			}
			key := permKey(p)
			if _, exists := seenPerms[key]; exists {
				return nil, fmt.Errorf("test %d: duplicate permutation", idx+1)
			}
			seenPerms[key] = struct{}{}
			perms[i] = p
		}
		if !allSumsEqual(sum) {
			return nil, fmt.Errorf("test %d: sums per position are not equal", idx+1)
		}
	}

	if extra, err := readToken(reader); err == nil {
		return nil, fmt.Errorf("unexpected extra output token %q", extra)
	} else if err != io.EOF {
		return nil, err
	}
	return verdicts, nil
}

func readToken(r *bufio.Reader) (string, error) {
	var b strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if ch > ' ' {
			b.WriteByte(ch)
			break
		}
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return b.String(), nil
			}
			return "", err
		}
		if ch <= ' ' {
			return b.String(), nil
		}
		b.WriteByte(ch)
	}
}

func permKey(p []int) string {
	var b strings.Builder
	for i, v := range p {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}

func validatePermutation(p []int, n int) error {
	seen := make([]bool, n+1)
	for _, v := range p {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range 1..%d", v, n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times in a permutation", v)
		}
		seen[v] = true
	}
	return nil
}

func allSumsEqual(sum []int) bool {
	if len(sum) == 0 {
		return true
	}
	target := sum[0]
	for _, v := range sum[1:] {
		if v != target {
			return false
		}
	}
	return true
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
