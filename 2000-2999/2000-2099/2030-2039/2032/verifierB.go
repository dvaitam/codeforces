package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const refSourceB = "./2032B.go"

type testCase struct {
	n int
	k int
}

type testInput struct {
	cases []testCase
}

func (ti testInput) buildInput() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(ti.cases))
	for _, cs := range ti.cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.k)
	}
	return b.String()
}

type parsedCase struct {
	ok      bool
	m       int
	p       []int
	rawLine string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := tc.buildInput()
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refCases, err := parseOutput(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candCases, err := parseOutput(candOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i, current := range tc.cases {
			if err := validateCase(current, refCases[i], candCases[i]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d case %d invalid: %v\ninput case: n=%d k=%d\nreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, i+1, err, current.n, current.k, formatCase(refCases[i]), formatCase(candCases[i]))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func validateCase(tc testCase, ref parsedCase, cand parsedCase) error {
	if !ref.ok {
		if cand.ok {
			return fmt.Errorf("expected -1 but got a partition")
		}
		return nil
	}
	if !cand.ok {
		return fmt.Errorf("partition exists but candidate output -1")
	}
	n, k := tc.n, tc.k
	if cand.m%2 == 0 || cand.m < 1 || cand.m > n {
		return fmt.Errorf("invalid m=%d (must be odd, 1<=m<=n)", cand.m)
	}
	if len(cand.p) != cand.m {
		return fmt.Errorf("expected %d partition points, got %d", cand.m, len(cand.p))
	}
	if cand.p[0] != 1 {
		return fmt.Errorf("first partition point must be 1, got %d", cand.p[0])
	}
	for i := 1; i < len(cand.p); i++ {
		if cand.p[i] <= cand.p[i-1] {
			return fmt.Errorf("partition points must be strictly increasing")
		}
	}
	if cand.p[len(cand.p)-1] > n {
		return fmt.Errorf("last partition point exceeds n")
	}
	// Validate segment lengths and compute medians.
	meds := make([]int, cand.m)
	for i := 0; i < cand.m; i++ {
		l := cand.p[i]
		var r int
		if i+1 < cand.m {
			r = cand.p[i+1] - 1
		} else {
			r = n
		}
		if r < l {
			return fmt.Errorf("empty segment between %d and %d", l, r)
		}
		length := r - l + 1
		if length%2 == 0 {
			return fmt.Errorf("segment [%d,%d] has even length", l, r)
		}
		meds[i] = l + length/2
	}
	if cand.p[len(cand.p)-1] != n && cand.m > 0 {
		// this actually cannot happen because last segment uses n, but guard
		lastR := n
		if lastR < cand.p[len(cand.p)-1] {
			return fmt.Errorf("segments do not cover array")
		}
	}
	// medians should be non-decreasing; sort copy to compute median if needed
	sortedMeds := append([]int(nil), meds...)
	sort.Ints(sortedMeds)
	target := sortedMeds[len(sortedMeds)/2]
	if target != k {
		return fmt.Errorf("median of medians is %d, expected %d", target, k)
	}
	return nil
}

func formatCase(pc parsedCase) string {
	if !pc.ok {
		return "-1"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", pc.m)
	for i, v := range pc.p {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}

func parseOutput(out string, t int) ([]parsedCase, error) {
	res := make([]parsedCase, 0, t)
	reader := strings.NewReader(out)
	for i := 0; i < t; i++ {
		token, err := nextToken(reader)
		if err != nil {
			return nil, fmt.Errorf("case %d: %v", i+1, err)
		}
		if token == "-1" {
			res = append(res, parsedCase{ok: false, rawLine: "-1"})
			continue
		}
		m, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("case %d: invalid m %q", i+1, token)
		}
		points := make([]int, m)
		for j := 0; j < m; j++ {
			valToken, err := nextToken(reader)
			if err != nil {
				return nil, fmt.Errorf("case %d: missing partition points (%v)", i+1, err)
			}
			val, err := strconv.Atoi(valToken)
			if err != nil {
				return nil, fmt.Errorf("case %d: invalid partition point %q", i+1, valToken)
			}
			points[j] = val
		}
		res = append(res, parsedCase{ok: true, m: m, p: points})
	}
	// ignore remaining tokens if any
	return res, nil
}

func nextToken(r *strings.Reader) (string, error) {
	var sb strings.Builder
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if sb.Len() == 0 {
				return "", fmt.Errorf("unexpected EOF")
			}
			return sb.String(), nil
		}
		if ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' {
			if sb.Len() > 0 {
				return sb.String(), nil
			}
			continue
		}
		sb.WriteRune(ch)
	}
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2032B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceB))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
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

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTests())
	tests = append(tests, edgeTests())

	rng := rand.New(rand.NewSource(2032))
	const limit = 200000
	used := 0
	var randomCases []testCase
	for used < limit {
		n := rng.Intn(199999) + 1
		if n%2 == 0 {
			n++
			if n > 199999 {
				n -= 2
			}
		}
		if n <= 0 {
			n = 1
		}
		k := rng.Intn(n) + 1
		randomCases = append(randomCases, testCase{n: n, k: k})
		used += n
		if len(randomCases) >= 5000 {
			break
		}
	}
	tests = append(tests, testInput{cases: randomCases})
	return tests
}

func sampleTests() testInput {
	return testInput{cases: []testCase{
		{n: 1, k: 1},
		{n: 3, k: 2},
		{n: 3, k: 3},
		{n: 15, k: 8},
	}}
}

func edgeTests() testInput {
	return testInput{cases: []testCase{
		{n: 1, k: 1},
		{n: 1, k: 1},
		{n: 199999, k: 1},
		{n: 199999, k: 199999},
		{n: 99999, k: 50001},
	}}
}
