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

const refSource2063D = "./2063D.go"

type testCase struct {
	name  string
	input string
}

type caseOutput struct {
	kmax int
	vals []int64
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
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d test outputs got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, refOut, candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i].kmax != candVals[i].kmax {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d failed: expected kmax=%d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i].kmax, candVals[i].kmax, tc.input, refOut, candOut)
				os.Exit(1)
			}
			if len(refVals[i].vals) != len(candVals[i].vals) {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d failed: expected %d values got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, len(refVals[i].vals), len(candVals[i].vals), tc.input, refOut, candOut)
				os.Exit(1)
			}
			for j := range refVals[i].vals {
				if refVals[i].vals[j] != candVals[i].vals[j] {
					fmt.Fprintf(os.Stderr, "test %d (%s) case %d value %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
						idx+1, tc.name, i+1, j+1, refVals[i].vals[j], candVals[i].vals[j], tc.input, refOut, candOut)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2063D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2063D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2063D)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
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

func parseOutput(input, output string) ([]caseOutput, error) {
	inFields := strings.Fields(input)
	if len(inFields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inFields[0])
	if err != nil || t < 1 || t > 30000 {
		return nil, fmt.Errorf("invalid test count %q", inFields[0])
	}

	outFields := strings.Fields(output)
	pos := 0
	results := make([]caseOutput, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(outFields) {
			return nil, fmt.Errorf("missing kmax for case %d", i+1)
		}
		kmax, err := strconv.Atoi(outFields[pos])
		if err != nil || kmax < 0 {
			return nil, fmt.Errorf("invalid kmax %q for case %d", outFields[pos], i+1)
		}
		pos++
		need := kmax
		vals := make([]int64, 0, need)
		for j := 0; j < need; j++ {
			if pos >= len(outFields) {
				return nil, fmt.Errorf("missing value %d for case %d", j+1, i+1)
			}
			v, err := strconv.ParseInt(outFields[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value %q for case %d", outFields[pos], i+1)
			}
			vals = append(vals, v)
			pos++
		}
		results = append(results, caseOutput{kmax: kmax, vals: vals})
	}
	if pos != len(outFields) {
		return nil, fmt.Errorf("extra output tokens: %v", outFields[pos:])
	}
	return results, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("kmax_zero", []nmCase{{n: 1, m: 1, a: []int{0}, b: []int{5}}}),
		makeManual("single_op", []nmCase{{n: 2, m: 1, a: []int{-3, 4}, b: []int{9}}}),
		makeManual("multi_small", []nmCase{
			{n: 3, m: 3, a: []int{0, 5, -4}, b: []int{2, -3, 9}},
			{n: 2, m: 4, a: []int{10, -10}, b: []int{1, 3, 7, 9}},
			{n: 4, m: 2, a: []int{-8, 5, 12, -1}, b: []int{20, -20}},
		}),
		makeManual("different_parity", []nmCase{{n: 5, m: 3, a: []int{1, 4, 9, 16, 25}, b: []int{-2, -3, -5}}}),
		makeManual("spread_coords", []nmCase{{n: 3, m: 5, a: []int{-1000000000, 0, 1000000000}, b: []int{-3, -2, -1, 1, 2}}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i))
	}

	tests = append(tests, largeTest(rng))

	return tests
}

type nmCase struct {
	n, m int
	a    []int
	b    []int
}

func makeManual(name string, cases []nmCase) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.m))
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range cs.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCount := rng.Intn(3) + 1
	totalN, totalM := 0, 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", caseCount))
	for i := 0; i < caseCount; i++ {
		n := rng.Intn(30) + 1
		m := rng.Intn(30) + 1
		if n+m < 3 {
			n = 2
			m = 1
		}
		totalN += n
		totalM += m
		a := makeRandomArray(n, rng)
		b := makeRandomArray(m, rng)
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		writeArray(&sb, a)
		sb.WriteByte('\n')
		writeArray(&sb, b)
		sb.WriteByte('\n')
	}
	name := fmt.Sprintf("random_%d_n%d_m%d", idx+1, totalN, totalM)
	return testCase{name: name, input: sb.String()}
}

func makeRandomArray(n int, rng *rand.Rand) []int {
	res := make([]int, 0, n)
	seen := make(map[int]struct{}, n*2)
	for len(res) < n {
		v := rng.Intn(2_000_000_001) - 1_000_000_000
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		res = append(res, v)
	}
	return res
}

func largeTest(rng *rand.Rand) testCase {
	n := 100000
	m := 90000
	a := make([]int, n)
	b := make([]int, m)
	for i := 0; i < n; i++ {
		a[i] = i * 2
	}
	for i := 0; i < m; i++ {
		b[i] = -1 - i*2
	}

	if rng.Intn(2) == 0 {
		reverseInts(a)
	}
	if rng.Intn(2) == 0 {
		reverseInts(b)
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	writeArray(&sb, a)
	sb.WriteByte('\n')
	writeArray(&sb, b)
	sb.WriteByte('\n')
	return testCase{name: "large_balanced", input: sb.String()}
}

func reverseInts(arr []int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
}
