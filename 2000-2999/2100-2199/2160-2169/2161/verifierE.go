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
)

const (
	refSource2161E = "2161E.go"
	refBinary2161E = "ref2161E.bin"
	maxTests       = 140
	maxTotalN      = 100000
)

type testCase struct {
	n int
	k int
	s string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on case %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2161E, refSource2161E)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2161E), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	ans := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		ans[i] = val
	}
	return ans, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n%s\n", tc.n, tc.k, tc.s)
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2161))
	var tests []testCase
	totalN := 0

	add := func(n, k int, s string) {
		tests = append(tests, testCase{n: n, k: k, s: s})
		totalN += n
	}

	add(5, 3, "0??0?")
	add(7, 7, "1??1??1")
	add(9, 5, "?????????")
	add(6, 3, "101010")
	add(6, 5, "??????")

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain < 3 {
			break
		}
		maxLen := 2000
		if remain < maxLen {
			maxLen = remain
		}
		n := rnd.Intn(maxLen-2) + 3
		k := rnd.Intn(n/2)*2 + 3 // ensure odd <= n
		if k > n {
			k = n
			if k%2 == 0 {
				k--
			}
			if k < 3 {
				k = 3
			}
		}
		var sb strings.Builder
		for i := 0; i < n; i++ {
			switch rnd.Intn(5) {
			case 0:
				sb.WriteByte('0')
			case 1:
				sb.WriteByte('1')
			default:
				sb.WriteByte('?')
			}
		}
		add(n, k, sb.String())
	}
	return tests
}
