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
	refSource = "1120B.go"
	refBinary = "ref1120B.bin"
	maxTests  = 80
	maxN      = 200
)

type testCase struct {
	n int
	a string
	b string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refPath)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	refOps, err := extractReferenceOps(refOut, len(tests))
	if err != nil {
		fmt.Printf("failed to parse reference output: %v\n", err)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	if err := verifyCandidateOutput(candOut, tests, refOps); err != nil {
		fmt.Println("candidate failed verification:", err)
		fmt.Println("Input used:")
		fmt.Println(string(input))
		return
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary), nil
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

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n%s\n%s\n", tc.n, tc.a, tc.b)
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for len(tests) < maxTests {
		n := rnd.Intn(maxN-1) + 2
		a := randomNumber(rnd, n)
		b := randomNumber(rnd, n)
		tests = append(tests, testCase{n: n, a: a, b: b})
	}
	return tests
}

func randomNumber(rnd *rand.Rand, n int) string {
	digits := make([]byte, n)
	digits[0] = byte(rnd.Intn(9)+1) + '0'
	for i := 1; i < n; i++ {
		digits[i] = byte(rnd.Intn(10)) + '0'
	}
	return string(digits)
}

func extractReferenceOps(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	pos := 0
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("insufficient tokens for test %d", i+1)
		}
		if tokens[pos] == "-1" {
			res[i] = -1
			pos++
			continue
		}
		val, err := strconv.ParseInt(tokens[pos], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[pos])
		}
		res[i] = val
		pos++
		// skip operations
		skip := int(val) * 2
		if pos+skip > len(tokens) {
			return nil, fmt.Errorf("reference output truncated for test %d", i+1)
		}
		pos += skip
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra tokens remaining in reference output")
	}
	return res, nil
}

func verifyCandidateOutput(out string, tests []testCase, refOps []int64) error {
	tokens := strings.Fields(out)
	pos := 0
	for i, tc := range tests {
		if pos >= len(tokens) {
			return fmt.Errorf("missing output for test %d", i+1)
		}
		refVal := refOps[i]
		if tokens[pos] == "-1" {
			if refVal != -1 {
				return fmt.Errorf("test %d: candidate claims impossible but reference has solution", i+1)
			}
			pos++
			continue
		}
		if refVal == -1 {
			return fmt.Errorf("test %d: candidate provides solution but reference says impossible", i+1)
		}
		cVal, err := strconv.ParseInt(tokens[pos], 10, 64)
		if err != nil {
			return fmt.Errorf("test %d: invalid integer %q", i+1, tokens[pos])
		}
		pos++
		if cVal != refVal {
			return fmt.Errorf("test %d: non-minimal operation count (got %d, expected %d)", i+1, cVal, refVal)
		}
		opCount := int(cVal)
		expectedTokens := pos + 2*opCount
		if expectedTokens > len(tokens) {
			return fmt.Errorf("test %d: insufficient operations provided", i+1)
		}
		ops := make([][2]int, opCount)
		for j := 0; j < opCount; j++ {
			d, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return fmt.Errorf("test %d: invalid position %q", i+1, tokens[pos])
			}
			s, err := strconv.Atoi(tokens[pos+1])
			if err != nil {
				return fmt.Errorf("test %d: invalid delta %q", i+1, tokens[pos+1])
			}
			ops[j] = [2]int{d, s}
			pos += 2
		}
		if err := simulate(tc, ops); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	if pos != len(tokens) {
		return fmt.Errorf("extra tokens in candidate output")
	}
	return nil
}

func simulate(tc testCase, ops [][2]int) error {
	n := tc.n
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = int(tc.a[i] - '0')
		b[i] = int(tc.b[i] - '0')
	}
	for _, op := range ops {
		pos := op[0]
		val := op[1]
		if pos < 1 || pos >= n {
			return fmt.Errorf("operation position %d out of range", pos)
		}
		if val != 1 && val != -1 {
			return fmt.Errorf("operation delta %d invalid", val)
		}
		if a[pos-1]+val < 0 || a[pos-1]+val > 9 {
			return fmt.Errorf("digit out of range after operation at %d", pos)
		}
		if pos-1 == 0 && a[0]+val == 0 {
			return fmt.Errorf("leading zero produced")
		}
		if a[pos]+val < 0 || a[pos]+val > 9 {
			return fmt.Errorf("digit out of range after operation at %d", pos)
		}
		a[pos-1] += val
		a[pos] += val
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return fmt.Errorf("numbers do not match target")
		}
	}
	return nil
}
