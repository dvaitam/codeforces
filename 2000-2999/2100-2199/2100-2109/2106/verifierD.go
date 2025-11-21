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

const refSource2106D = "2000-2999/2100-2199/2100-2109/2106/2106D.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refVals, err := parseOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d answers got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, refOut, candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2106D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2106D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2106D)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(input, output string) ([]int64, error) {
	inFields := strings.Fields(input)
	if len(inFields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inFields[0])
	if err != nil || t < 1 || t > 10000 {
		return nil, fmt.Errorf("invalid test count %q", inFields[0])
	}

	fields := strings.Fields(output)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("simple_possible_no_add", []manCase{
			{n: 3, m: 2, a: []int{5, 4, 3}, b: []int{3, 3}, expect: 0},
		}),
		makeManual("simple_need_add", []manCase{
			{n: 3, m: 2, a: []int{1, 2, 3}, b: []int{4, 4}, expect: -1}, // even add cannot, but reference decides
			{n: 3, m: 2, a: []int{1, 2, 3}, b: []int{2, 4}, expect: 4},
		}),
		makeManual("already_matches", []manCase{
			{n: 5, m: 3, a: []int{2, 5, 6, 1, 9}, b: []int{2, 5, 6}, expect: 0},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	tests = append(tests, largeTest())
	return tests
}

type manCase struct {
	n, m   int
	a, b   []int
	expect int64 // unused in logic; kept for readability
}

func makeManual(name string, cases []manCase) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.m))
		writeArray(&sb, cs.a)
		sb.WriteByte('\n')
		writeArray(&sb, cs.b)
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(40) + 1
		m := rng.Intn(n) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(20) + 1
		}
		b := make([]int, m)
		for j := 0; j < m; j++ {
			b[j] = rng.Intn(20) + 1
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		writeArray(&sb, a)
		sb.WriteByte('\n')
		writeArray(&sb, b)
		sb.WriteByte('\n')
	}
	return testCase{name: fmt.Sprintf("random_%d", idx+1), input: sb.String()}
}

func largeTest() testCase {
	n := 200000
	m := 200000
	a := make([]int, n)
	b := make([]int, m)
	for i := 0; i < n; i++ {
		a[i] = (i % 5) + 1
	}
	for i := 0; i < m; i++ {
		b[i] = (i % 5) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	writeArray(&sb, a)
	sb.WriteByte('\n')
	writeArray(&sb, b)
	sb.WriteByte('\n')
	return testCase{name: "large_equal", input: sb.String()}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
}
