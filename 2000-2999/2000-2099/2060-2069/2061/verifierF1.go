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
	randomTests   = 120
	maxTotalLen   = 180000
	maxCaseLength = 2000
)

type testCase struct {
	s string
	t string
}

// solveCase implements the correct greedy algorithm (translated from the C++ reference).
// It returns the minimum number of adjacent block swaps, or -1 if impossible.
func solveCase(s, t string) int {
	n := len(s)
	ans := 0
	var ch byte
	c := 0

	for i := 0; i < n; i++ {
		x := s[i]
		if c > 0 {
			x = ch
			c--
		}
		if x != t[i] {
			if i > 0 && t[i-1] != t[i] {
				return -1
			}
			p := i + c + 1
			for p < n && s[p] == x {
				p++
			}
			if p >= n {
				return -1
			}
			q := p + 1
			for q < n && s[q] == s[p] {
				q++
			}
			l1 := p - i
			l2 := q - p
			ans++
			for k := i + 1; k < i+l2; k++ {
				if t[k] != t[k-1] {
					return -1
				}
			}
			i = i + l2 - 1
			c = l1
			ch = x
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	input := buildInput(tests)

	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}

	got, err := parseOutputs(gotRaw, len(tests))
	if err != nil {
		fail("could not parse candidate output: %v", err)
	}

	for i, tc := range tests {
		expected := solveCase(tc.s, tc.t)
		if expected != got[i] {
			fail("mismatch at test %d (s=%q t=%q): expected %d, got %d",
				i+1, tc.s, tc.t, expected, got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+10)

	// Deterministic coverage.
	tests = append(tests,
		testCase{s: "0", t: "0"},
		testCase{s: "1", t: "1"},
		testCase{s: "01", t: "10"},
		testCase{s: "10", t: "01"},
		testCase{s: "111000", t: "000111"},
		testCase{s: "000111", t: "111000"},
		testCase{s: "000110", t: "111000"},
		testCase{s: "00011", t: "00011"},
		testCase{s: "0000011111", t: "1111100000"},
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
		n := rng.Intn(maxLen) + 1

		var sbS strings.Builder
		var sbT strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sbS.WriteByte('0')
			} else {
				sbS.WriteByte('1')
			}
			if rng.Intn(2) == 0 {
				sbT.WriteByte('0')
			} else {
				sbT.WriteByte('1')
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
