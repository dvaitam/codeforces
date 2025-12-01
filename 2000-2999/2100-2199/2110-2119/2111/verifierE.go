package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource2111E = "./2111E.go"

type testCase struct {
	name  string
	input string
	ns    []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		refVals, err := parseOutput(tc.ns, refOut)
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
		candVals, err := parseOutput(tc.ns, candOut)
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
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2111E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2111E.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2111E)
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

func parseOutput(ns []int, output string) ([]string, error) {
	fields := strings.Fields(output)
	if len(fields) != len(ns) {
		return nil, fmt.Errorf("expected %d answers, got %d", len(ns), len(fields))
	}
	for i, s := range fields {
		if len(s) != ns[i] {
			return nil, fmt.Errorf("answer %d length %d, expected %d", i+1, len(s), ns[i])
		}
		for _, ch := range s {
			if ch != 'a' && ch != 'b' && ch != 'c' {
				return nil, fmt.Errorf("answer %d contains invalid char %q", i+1, ch)
			}
		}
	}
	return fields, nil
}

func buildTests() []testCase {
	tests := []testCase{
		buildManualSample(),
		makeManual("small_no_ops", []manualCase{{s: "abc", ops: []op{}, n: 3, q: 0}}),
		makeManual("single_change", []manualCase{
			{s: "c", ops: []op{{'c', 'a'}}, n: 1, q: 1},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 120; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	tests = append(tests, largeStress())
	return tests
}

type op struct {
	x, y byte
}

type manualCase struct {
	s         string
	ops       []op
	n, q      int
	randomLen bool
}

func buildManualSample() testCase {
	const sampleInput = `2
2 2
cb
c b
b a
5 5
abcac
c a
b c
a b
c b
b a
`
	return testCase{
		name:  "statement_sample",
		input: sampleInput,
		ns:    []int{2, 5},
	}
}

func makeManual(name string, cases []manualCase) testCase {
	var sb strings.Builder
	ns := make([]int, 0, len(cases))
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		ns = append(ns, cs.n)
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.q))
		sb.WriteString(cs.s)
		sb.WriteByte('\n')
		for _, o := range cs.ops {
			sb.WriteByte(o.x)
			sb.WriteByte(' ')
			sb.WriteByte(o.y)
			sb.WriteByte('\n')
		}
	}
	return testCase{name: name, input: sb.String(), ns: ns}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	ns := make([]int, t)
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(50) + 1
		q := rng.Intn(50) + 1
		ns[i] = n
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		sb.WriteString(randomString(rng, n))
		sb.WriteByte('\n')
		for j := 0; j < q; j++ {
			x := byte('a' + rng.Intn(3))
			y := byte('a' + rng.Intn(3))
			sb.WriteByte(x)
			sb.WriteByte(' ')
			sb.WriteByte(y)
			sb.WriteByte('\n')
		}
	}
	return testCase{name: fmt.Sprintf("random_%d", idx+1), input: sb.String(), ns: ns}
}

func largeStress() testCase {
	n := 200000
	q := 200000
	var sb strings.Builder
	ns := []int{n}
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	sb.WriteString(strings.Repeat("abc", n/3) + strings.Repeat("a", n%3))
	sb.WriteByte('\n')
	chars := []byte{'a', 'b', 'c'}
	for i := 0; i < q; i++ {
		x := chars[i%3]
		y := chars[(i+1)%3]
		sb.WriteByte(x)
		sb.WriteByte(' ')
		sb.WriteByte(y)
		sb.WriteByte('\n')
	}
	return testCase{name: "large_stress", input: sb.String(), ns: ns}
}

func randomString(rng *rand.Rand, n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(byte('a' + rng.Intn(3)))
	}
	return b.String()
}
