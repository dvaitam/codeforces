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

const (
	globalLimit = 200000
	mod         = 998244353
)

type testCase struct {
	n     int
	q     int
	s     string
	flips []int
}

func buildReference() (string, error) {
	path := "./2077C_ref.bin"
	cmd := exec.Command("go", "build", "-o", path, "2077C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.q))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		for _, idx := range tc.flips {
			sb.WriteString(strconv.Itoa(idx))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func sampleInput() string {
	cases := []testCase{
		{n: 3, q: 2, s: "010", flips: []int{1, 3}},
		{n: 4, q: 3, s: "1010", flips: []int{1, 2, 4}},
		{n: 5, q: 4, s: "11001", flips: []int{5, 4, 3, 2}},
	}
	return buildInput(cases)
}

func randomBinaryString(rng *rand.Rand, n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	return string(bytes)
}

func randomTestInput(rng *rand.Rand, maxCases, maxN, maxQ int) string {
	targetCases := rng.Intn(maxCases) + 1
	cases := []testCase{}
	sumN, sumQ := 0, 0
	for len(cases) < targetCases {
		remainN := globalLimit - sumN
		remainQ := globalLimit - sumQ
		if remainN <= 0 || remainQ <= 0 {
			break
		}
		nUpper := maxN
		if nUpper > remainN {
			nUpper = remainN
		}
		if nUpper < 1 {
			nUpper = 1
		}
		n := rng.Intn(nUpper) + 1

		qUpper := maxQ
		if qUpper > remainQ {
			qUpper = remainQ
		}
		if qUpper < 1 {
			qUpper = 1
		}
		q := rng.Intn(qUpper) + 1

		s := randomBinaryString(rng, n)
		flips := make([]int, q)
		for i := 0; i < q; i++ {
			flips[i] = rng.Intn(n) + 1
		}
		cases = append(cases, testCase{n: n, q: q, s: s, flips: flips})
		sumN += n
		sumQ += q
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 1, q: 1, s: "0", flips: []int{1}})
	}
	return buildInput(cases)
}

func alternatingString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	return string(bytes)
}

func buildLargeStress() string {
	n := 200000
	q := 200000
	cases := []testCase{
		{
			n:     n,
			q:     q,
			s:     alternatingString(n),
			flips: make([]int, q),
		},
	}
	for i := 0; i < q; i++ {
		cases[0].flips[i] = (i*37)%n + 1
	}
	return buildInput(cases)
}

func smallEdgeCases() string {
	cases := []testCase{
		{n: 1, q: 5, s: "0", flips: []int{1, 1, 1, 1, 1}},
		{n: 2, q: 4, s: "10", flips: []int{1, 2, 1, 2}},
		{n: 5, q: 5, s: "00000", flips: []int{1, 3, 5, 2, 4}},
	}
	return buildInput(cases)
}

func compareOutputs(exp, got string) error {
	expVals := strings.Fields(exp)
	gotVals := strings.Fields(got)
	if len(expVals) != len(gotVals) {
		return fmt.Errorf("expected %d outputs but got %d", len(expVals), len(gotVals))
	}
	for i := range expVals {
		if expVals[i] != gotVals[i] {
			return fmt.Errorf("output mismatch at line %d: expected %s, got %s", i+1, expVals[i], gotVals[i])
		}
	}
	return nil
}

func buildTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{}
	tests = append(tests, sampleInput())
	tests = append(tests, smallEdgeCases())
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTestInput(rng, 5, 50, 50))
	}
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestInput(rng, 5, 500, 500))
	}
	tests = append(tests, randomTestInput(rng, 10, 2000, 2000))
	tests = append(tests, buildLargeStress())
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var refPath string
	fail := func(format string, args ...interface{}) {
		if refPath != "" {
			_ = os.Remove(refPath)
		}
		fmt.Fprintf(os.Stderr, format+"\n", args...)
		os.Exit(1)
	}

	var err error
	refPath, err = buildReference()
	if err != nil {
		fail("%v", err)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, input := range tests {
		expect, err := runProgram(refPath, input)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		got, err := runProgram(bin, input)
		if err != nil {
			fail("test %d: runtime error: %v\ninput:\n%s", idx+1, err, input)
		}
		if err := compareOutputs(expect, got); err != nil {
			fail("test %d failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s", idx+1, err, input, expect, got)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
