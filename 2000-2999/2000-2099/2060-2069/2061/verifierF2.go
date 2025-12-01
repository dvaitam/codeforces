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

const (
	refSource     = "./2061F2.go"
	randomTests   = 140
	maxTotalLen   = 250000 // keep well under statement limit of 4e5
	maxCaseLength = 4000
)

type testCase struct {
	s string
	t string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}

	expect, err := parseOutputs(expectRaw, len(tests))
	if err != nil {
		fail("could not parse reference output: %v", err)
	}
	got, err := parseOutputs(gotRaw, len(tests))
	if err != nil {
		fail("could not parse candidate output: %v", err)
	}

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at test %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2061F2-ref-*")
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

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+10)

	// Deterministic coverage for edge and sample-like situations.
	tests = append(tests,
		testCase{s: "0", t: "0"},
		testCase{s: "1", t: "?"},
		testCase{s: "01", t: "??"},
		testCase{s: "10", t: "10"},
		testCase{s: "111000", t: "010101"},
		testCase{s: "000111", t: "111000"},
		testCase{s: "010101", t: "??????"},
		testCase{s: "0000", t: "1111"},
		testCase{s: "1010101", t: "1?0?0??"},
		testCase{s: "001100", t: "??11??"},
	)

	totalLen := 0
	for _, tc := range tests {
		totalLen += len(tc.s)
	}

	for i := 0; i < randomTests && totalLen < maxTotalLen; i++ {
		remain := maxTotalLen - totalLen
		maxLen := maxCaseLength
		if maxLen > remain {
			maxLen = remain
		}
		if maxLen < 1 {
			break
		}

		// Bias towards shorter cases but occasionally generate long strings.
		shortCap := maxLen
		if shortCap > 40 {
			shortCap = 40
		}
		n := rng.Intn(shortCap) + 1
		if rng.Intn(5) == 0 {
			n = rng.Intn(maxLen) + 1
		}

		var sbS strings.Builder
		var sbT strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sbS.WriteByte('0')
			} else {
				sbS.WriteByte('1')
			}

			x := rng.Intn(3)
			if x == 0 {
				sbT.WriteByte('0')
			} else if x == 1 {
				sbT.WriteByte('1')
			} else {
				sbT.WriteByte('?')
			}
		}

		tests = append(tests, testCase{s: sbS.String(), t: sbT.String()})
		totalLen += n
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(tc.t)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
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

func parseOutputs(out string, t int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int, t)
	for i, tok := range tokens {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
