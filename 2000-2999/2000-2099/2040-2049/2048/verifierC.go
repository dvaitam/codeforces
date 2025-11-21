package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	s string
}

type selection struct {
	l1 int
	r1 int
	l2 int
	r2 int
}

const maxTotalLen = 5000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	input := serializeInput(tests)
	lengths := make([]int, len(tests))
	for i, tc := range tests {
		lengths[i] = len(tc.s)
	}

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refSelections, err := parseSelections(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid reference output: %v\n", err)
		os.Exit(1)
	}

	expectedVals := make([]*big.Int, len(tests))
	for i := range tests {
		if err := validateSelection(refSelections[i], lengths[i]); err != nil {
			fmt.Fprintf(os.Stderr, "reference selection invalid for test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expectedVals[i] = xorValue(tests[i].s, refSelections[i])
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}
	candSelections, err := parseSelections(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if err := validateSelection(candSelections[i], lengths[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid selection: %v\nstring: %s\n", i+1, err, tests[i].s)
			os.Exit(1)
		}
		got := xorValue(tests[i].s, candSelections[i])
		if got.Cmp(expectedVals[i]) != 0 {
			fmt.Fprintf(os.Stderr, "test %d mismatch\nstring: %s\nexpected XOR: %s\ngot XOR: %s\n", i+1, tests[i].s, expectedVals[i].String(), got.String())
			fmt.Fprintf(os.Stderr, "candidate selection: %d %d %d %d\n", candSelections[i].l1, candSelections[i].r1, candSelections[i].l2, candSelections[i].r2)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2048C.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2048C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
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

func parseSelections(out string, tests int) ([]selection, error) {
	fields := strings.Fields(out)
	if len(fields) != tests*4 {
		return nil, fmt.Errorf("expected %d integers, got %d", tests*4, len(fields))
	}
	res := make([]selection, tests)
	for i := 0; i < tests; i++ {
		idx := i * 4
		vals := make([]int, 4)
		for j := 0; j < 4; j++ {
			v, err := strconv.Atoi(fields[idx+j])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q: %v", fields[idx+j], err)
			}
			vals[j] = v
		}
		res[i] = selection{l1: vals[0], r1: vals[1], l2: vals[2], r2: vals[3]}
	}
	return res, nil
}

func validateSelection(sel selection, length int) error {
	if sel.l1 < 1 || sel.l1 > sel.r1 || sel.r1 > length {
		return fmt.Errorf("invalid first substring: [%d, %d]", sel.l1, sel.r1)
	}
	if sel.l2 < 1 || sel.l2 > sel.r2 || sel.r2 > length {
		return fmt.Errorf("invalid second substring: [%d, %d]", sel.l2, sel.r2)
	}
	return nil
}

func xorValue(s string, sel selection) *big.Int {
	sub1 := s[sel.l1-1 : sel.r1]
	sub2 := s[sel.l2-1 : sel.r2]
	a := new(big.Int)
	b := new(big.Int)
	a.SetString(sub1, 2)
	b.SetString(sub2, 2)
	return new(big.Int).Xor(a, b)
}

func buildTests() []testCase {
	tests := []testCase{
		{s: "1"},
		{s: "11"},
		{s: "10"},
		{s: "101"},
		{s: "1010"},
		{s: "10001"},
		{s: "111111"},
		{s: "1011001"},
		{s: "1100110011"},
	}
	total := 0
	for _, tc := range tests {
		total += len(tc.s)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < maxTotalLen {
		remaining := maxTotalLen - total
		maxLen := remaining
		if maxLen > 300 {
			maxLen = 300
		}
		if maxLen == 0 {
			break
		}
		length := rng.Intn(maxLen) + 1
		str := randomBinaryString(rng, length)
		tests = append(tests, testCase{s: str})
		total += length
	}
	return tests
}

func randomBinaryString(rng *rand.Rand, length int) string {
	if length <= 0 {
		return "1"
	}
	b := make([]byte, length)
	b[0] = '1'
	for i := 1; i < length; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}
