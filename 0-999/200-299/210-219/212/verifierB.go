package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	s    string
	sets []string
}

const maxStringLen = 100000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeTests(tests)

	expected, err := runAndParse(refBin, input, totalQueries(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, totalQueries(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	idx := 0
	for ti, tc := range tests {
		for qi := range tc.sets {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "Mismatch on test %d query %d: expected %d got %d\ns=%s set=%s\n",
					ti+1, qi+1, expected[idx], got[idx], tc.s, tc.sets[qi])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_212B.bin"
	cmd := exec.Command("go", "build", "-o", refName, "212B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, ansCount int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != ansCount {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", ansCount, len(fields), out)
	}
	res := make([]int64, ansCount)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	totalLen := totalStringLen(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for totalLen < maxStringLen {
		length := rng.Intn(1000) + 1
		if totalLen+length > maxStringLen {
			length = maxStringLen - totalLen
		}
		var sb strings.Builder
		for i := 0; i < length; i++ {
			sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
		}
		s := sb.String()
		q := rng.Intn(20) + 1
		setList := make([]string, q)
		for i := 0; i < q; i++ {
			setSize := rng.Intn(5) + 1
			mask := make(map[byte]struct{})
			var str strings.Builder
			for str.Len() < setSize {
				ch := alphabet[rng.Intn(len(alphabet))]
				if _, ok := mask[ch]; !ok {
					mask[ch] = struct{}{}
					str.WriteByte(ch)
				}
			}
			setList[i] = str.String()
		}
		tests = append(tests, testCase{s: s, sets: setList})
		totalLen += len(s)
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{s: "aaaaa", sets: []string{"a", "ab"}},
		{s: "abacaba", sets: []string{"ac", "ba", "abc"}},
		{s: "xyz", sets: []string{"x", "yz"}},
	}
}

func serializeTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(len(tc.sets)))
		sb.WriteByte('\n')
		for _, c := range tc.sets {
			sb.WriteString(c)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func totalQueries(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.sets)
	}
	return total
}

func totalStringLen(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.s)
	}
	return total
}
