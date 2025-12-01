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

const refSourceD2 = "./2027D2.go"

type caseData struct {
	n int
	m int
	a []int64
	b []int64
}

type testInput struct {
	cases []caseData
}

func (ti testInput) buildInput() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(ti.cases))
	for _, cs := range ti.cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.m)
		for i, v := range cs.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.FormatInt(v, 10))
		}
		b.WriteByte('\n')
		for i, v := range cs.b {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.FormatInt(v, 10))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

type result struct {
	ok   bool
	cost int64
	ways int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/candidate")
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
		refRes, err := parseOutputs(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candRes, err := parseOutputs(candOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := range refRes {
			if refRes[i].ok != candRes[i].ok {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: feasibility mismatch\ninput:\n%sreference: %v\ncandidate: %v\n",
					idx+1, i+1, input, formatResult(refRes[i]), formatResult(candRes[i]))
				os.Exit(1)
			}
			if refRes[i].ok {
				if refRes[i].cost != candRes[i].cost || refRes[i].ways != candRes[i].ways {
					fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%sreference: %v\ncandidate: %v\n",
						idx+1, i+1, input, formatResult(refRes[i]), formatResult(candRes[i]))
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func formatResult(res result) string {
	if !res.ok {
		return "-1"
	}
	return fmt.Sprintf("%d %d", res.cost, res.ways)
}

func parseOutputs(out string, expected int) ([]result, error) {
	fields := strings.Fields(out)
	res := make([]result, 0, expected)
	idx := 0
	for i := 0; i < expected; i++ {
		if idx >= len(fields) {
			return nil, fmt.Errorf("not enough outputs for case %d", i+1)
		}
		if fields[idx] == "-1" {
			res = append(res, result{ok: false})
			idx++
			continue
		}
		if idx+1 >= len(fields) {
			return nil, fmt.Errorf("missing cost/ways for case %d", i+1)
		}
		cost, err1 := strconv.ParseInt(fields[idx], 10, 64)
		ways, err2 := strconv.ParseInt(fields[idx+1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid integers for case %d", i+1)
		}
		res = append(res, result{ok: true, cost: cost, ways: ways})
		idx += 2
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra output tokens")
	}
	return res, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2027D2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceD2))
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

	rng := rand.New(rand.NewSource(2027))
	for len(tests) < 30 {
		tests = append(tests, randomBundle(rng, 300000))
	}

	tests = append(tests, stressTest())
	return tests
}

func sampleTests() testInput {
	return testInput{cases: []caseData{
		{n: 4, m: 2, a: []int64{9, 3, 4, 3}, b: []int64{11, 7}},
		{n: 1, m: 2, a: []int64{20}, b: []int64{19, 18}},
		{n: 5, m: 2, a: []int64{10, 2, 2, 5, 2}, b: []int64{1, 1}},
	}}
}

func edgeTests() testInput {
	return testInput{cases: []caseData{
		{n: 1, m: 1, a: []int64{1}, b: []int64{1}},
		{n: 3, m: 3, a: []int64{1, 1, 1}, b: []int64{3, 2, 1}},
		{n: 2, m: 2, a: []int64{5, 1}, b: []int64{3, 2}},
	}}
}

func randomBundle(rng *rand.Rand, limit int) testInput {
	var cases []caseData
	used := 0
	for used < limit {
		n := rng.Intn(500) + 1
		m := rng.Intn(500) + 1
		if n*m > limit-used {
			n = maxInt(1, (limit-used)/m)
			if n == 0 {
				break
			}
		}
		if n*m == 0 {
			break
		}
		a := make([]int64, n)
		for i := range a {
			a[i] = rng.Int63n(1_000_000_000) + 1
		}
		b := make([]int64, m)
		prev := int64(1_000_000_000)
		for i := 0; i < m; i++ {
			decr := rng.Int63n(1_000_000_000/(int64(i)+1)) + 1
			if prev-decr < 1 {
				decr = prev - 1
				if decr < 1 {
					decr = 1
				}
			}
			b[i] = prev
			prev -= decr
			if prev < 1 {
				prev = 1
			}
		}
		// ensure strictly decreasing
		for i := 1; i < m; i++ {
			if b[i] >= b[i-1] {
				b[i] = maxInt64(1, b[i-1]-1)
			}
		}
		cases = append(cases, caseData{n: n, m: m, a: a, b: b})
		used += n * m
		if rng.Intn(4) == 0 {
			break
		}
	}
	if len(cases) == 0 {
		cases = append(cases, caseData{n: 1, m: 1, a: []int64{1}, b: []int64{1}})
	}
	return testInput{cases: cases}
}

func stressTest() testInput {
	// Create a case with n*m close to limit
	n := 300
	m := 1000
	if n*m > 300000 {
		m = 300000 / n
	}
	a := make([]int64, n)
	for i := range a {
		a[i] = 1_000_000_000
	}
	b := make([]int64, m)
	val := int64(1_000_000_000)
	for i := 0; i < m; i++ {
		b[i] = val
		val -= int64(i%5 + 1)
		if val < 1 {
			val = 1
		}
	}
	for i := 1; i < m; i++ {
		if b[i] >= b[i-1] {
			b[i] = maxInt64(1, b[i-1]-1)
		}
	}
	return testInput{cases: []caseData{
		{n: n, m: m, a: a, b: b},
	}}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
