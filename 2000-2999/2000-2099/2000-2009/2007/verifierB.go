package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const refSource = "2007B.go"

type op struct {
	typ string
	l   int64
	r   int64
}

type singleCase struct {
	a   []int64
	ops []op
}

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at token %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "oracle-2007B-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleB")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, expected int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		buildSampleStyle(),
		buildSingleElement(),
		buildNegativeDrift(),
		buildLargeValues(),
		buildLongBounces(),
	}

	rng := rand.New(rand.NewSource(20072007))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func buildSampleStyle() testCase {
	case1 := singleCase{
		a: []int64{1, 2, 3, 2, 1},
		ops: []op{
			{typ: "+", l: 1, r: 3},
			{typ: "-", l: 2, r: 3},
			{typ: "+", l: 1, r: 2},
			{typ: "+", l: 2, r: 4},
			{typ: "-", l: 6, r: 8},
		},
	}
	case2 := singleCase{
		a: []int64{1, 3, 3, 4, 5},
		ops: []op{
			{typ: "+", l: 1, r: 4},
			{typ: "+", l: 2, r: 3},
			{typ: "-", l: 4, r: 5},
			{typ: "-", l: 3, r: 3},
			{typ: "-", l: 2, r: 6},
		},
	}
	return buildTest("sample_style", []singleCase{case1, case2})
}

func buildSingleElement() testCase {
	ops := []op{
		{typ: "+", l: 1, r: 1},
		{typ: "-", l: 5, r: 10},
		{typ: "+", l: 2, r: 7},
		{typ: "-", l: 1, r: 3},
		{typ: "+", l: 100, r: 200},
		{typ: "-", l: 1, r: 150},
	}
	return buildTest("single_element", []singleCase{
		{
			a:   []int64{5},
			ops: ops,
		},
	})
}

func buildNegativeDrift() testCase {
	case1 := singleCase{
		a: []int64{5, 1, 3, 7},
		ops: []op{
			{typ: "-", l: 3, r: 7},
			{typ: "+", l: 1, r: 5},
			{typ: "-", l: 6, r: 6},
			{typ: "-", l: 1, r: 4},
			{typ: "+", l: 2, r: 10},
		},
	}
	case2 := singleCase{
		a: []int64{2, 2, 2, 2},
		ops: []op{
			{typ: "-", l: 2, r: 2},
			{typ: "-", l: 1, r: 3},
			{typ: "+", l: 1, r: 5},
			{typ: "+", l: 2, r: 4},
		},
	}
	return buildTest("negative_drift", []singleCase{case1, case2})
}

func buildLargeValues() testCase {
	case1 := singleCase{
		a: []int64{1_000_000_000, 1_000_000_000 - 1, 1_000_000_000 - 2},
		ops: []op{
			{typ: "+", l: 1_000_000_000, r: 1_000_000_000},
			{typ: "-", l: 999_999_999, r: 1_000_000_000},
			{typ: "+", l: 1, r: 1_000_000_000},
			{typ: "-", l: 1_000_000_000, r: 1_000_000_000},
			{typ: "+", l: 500_000_000, r: 1_000_000_000},
		},
	}
	return buildTest("large_values", []singleCase{case1})
}

func buildLongBounces() testCase {
	ops := make([]op, 0, 60)
	curr := int64(10)
	for i := 0; i < 30; i++ {
		ops = append(ops, op{typ: "+", l: curr, r: curr})
		curr++
		ops = append(ops, op{typ: "-", l: curr - 2, r: curr})
	}
	return buildTest("long_bounces", []singleCase{
		{
			a:   []int64{10, 9, 8},
			ops: ops,
		},
	})
}

func buildTest(name string, cases []singleCase) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	outputs := 0
	for _, cs := range cases {
		outputs += len(cs.ops)
		fmt.Fprintf(&b, "%d %d\n", len(cs.a), len(cs.ops))
		for i, v := range cs.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.FormatInt(v, 10))
		}
		b.WriteByte('\n')
		for _, op := range cs.ops {
			fmt.Fprintf(&b, "%s %d %d\n", op.typ, op.l, op.r)
		}
	}
	return testCase{
		name:    name,
		input:   b.String(),
		outputs: outputs,
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCount := rng.Intn(3) + 1
	cases := make([]singleCase, caseCount)
	for i := 0; i < caseCount; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(40) + 1
		orig := make([]int64, n)
		for j := 0; j < n; j++ {
			orig[j] = rng.Int63n(1_000_000_000) + 1
		}
		current := append([]int64(nil), orig...)
		ops := make([]op, m)
		for j := 0; j < m; j++ {
			l, r := randomRange(rng, current)
			t := "-"
			if rng.Intn(2) == 0 {
				t = "+"
			}
			ops[j] = op{typ: t, l: l, r: r}
			for k, v := range current {
				if l <= v && v <= r {
					if t == "+" {
						current[k]++
					} else {
						current[k]--
					}
				}
			}
		}
		cases[i] = singleCase{a: orig, ops: ops}
	}
	return buildTest(fmt.Sprintf("random_%d", idx), cases)
}

func randomRange(rng *rand.Rand, values []int64) (int64, int64) {
	if len(values) > 0 && rng.Intn(3) == 0 {
		base := values[rng.Intn(len(values))]
		left := base - int64(rng.Intn(3)+1)
		if left < 1 {
			left = 1
		}
		right := base + int64(rng.Intn(3)+1)
		if right < left {
			right = left
		}
		return left, right
	}
	l := rng.Int63n(1_000_000_000) + 1
	span := 1_000_000_000 - l + 1
	r := l + rng.Int63n(span)
	return l, r
}
