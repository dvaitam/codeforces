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

const refSource = "1998B.go"

type testCase struct {
	input string
}

type testInstance struct {
	p []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for ti, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", ti+1, err, tc.input)
			os.Exit(1)
		}

		expect, err := parseReference(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", ti+1, err, tc.input, got)
			os.Exit(1)
		}

		if err := validateCandidate(tc.input, got, expect); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%soutput:\n%s\n", ti+1, err, tc.input, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1998B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type refResult struct {
	t        int
	shift    int
	permutes []string
}

func parseReference(input, output string) (refResult, error) {
	inLines := strings.Fields(input)
	outLines := strings.Fields(output)
	if len(inLines) == 0 {
		return refResult{}, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inLines[0])
	if err != nil {
		return refResult{}, fmt.Errorf("bad test count: %v", err)
	}
	permIdx := 0
	expect := make([][]int, t)
	pos := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(inLines) {
			return refResult{}, fmt.Errorf("input truncated")
		}
		n, _ := strconv.Atoi(inLines[pos])
		pos++
		if pos+n > len(inLines) {
			return refResult{}, fmt.Errorf("input permutation truncated")
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i], _ = strconv.Atoi(inLines[pos+i])
		}
		pos += n

		if permIdx+n > len(outLines) {
			return refResult{}, fmt.Errorf("reference output truncated")
		}
		q := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(outLines[permIdx+i])
			if err != nil {
				return refResult{}, fmt.Errorf("bad number in reference output: %v", err)
			}
			q[i] = val
		}
		permIdx += n
		expect[caseIdx] = p

		shift := (matchShift(p, q, n))
		if shift == -1 {
			return refResult{}, fmt.Errorf("reference q is not rotation")
		}
	}

	return refResult{t: t, shift: -1}, nil
}

func matchShift(p, q []int, n int) int {
	if n == 0 {
		return 0
	}
	pos := make(map[int]int, n)
	for i, val := range q {
		pos[val] = i
	}
	shift := -1
	for i, val := range p {
		j := pos[val]
		cur := (j - i + n) % n
		if shift == -1 {
			shift = cur
		} else if shift != cur {
			return -1
		}
	}
	return shift
}

func validateCandidate(input, output string, expect refResult) error {
	in := strings.Fields(input)
	out := strings.Fields(output)
	if len(in) == 0 {
		return fmt.Errorf("empty input")
	}
	t, _ := strconv.Atoi(in[0])
	posOut := 0
	posIn := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n, _ := strconv.Atoi(in[posIn])
		posIn++
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i], _ = strconv.Atoi(in[posIn+i])
		}
		posIn += n

		if posOut+n > len(out) {
			return fmt.Errorf("output truncated for case %d", caseIdx+1)
		}
		q := make([]int, n)
		seen := make([]bool, n+1)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(out[posOut+i])
			if err != nil {
				return fmt.Errorf("case %d: invalid integer %q", caseIdx+1, out[posOut+i])
			}
			if val < 1 || val > n {
				return fmt.Errorf("case %d: value %d out of range", caseIdx+1, val)
			}
			if seen[val] {
				return fmt.Errorf("case %d: value %d repeated", caseIdx+1, val)
			}
			seen[val] = true
			q[i] = val
		}
		posOut += n

		if matchShift(p, q, n) == -1 {
			return fmt.Errorf("case %d: output is not rotation of input", caseIdx+1)
		}
	}
	if posOut != len(out) {
		return fmt.Errorf("extra tokens in output")
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(19981998))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]testInstance{
		{p: []int{1}},
		{p: []int{1, 2}},
	}))

	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(5)+1))
	}

	tests = append(tests, limitCase())

	return tests
}

func sampleTest() testCase {
	return makeTest([]testInstance{
		{p: []int{2, 1}},
		{p: []int{1, 2, 3, 4, 5}},
		{p: []int{4, 7, 5, 1, 2, 6, 3}},
	})
}

func makeTest(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(instances))
	for _, inst := range instances {
		fmt.Fprintln(&b, len(inst.p))
		for i, v := range inst.p {
			if i > 0 {
				fmt.Fprint(&b, " ")
			}
			fmt.Fprint(&b, v)
		}
		fmt.Fprintln(&b)
	}
	return testCase{input: b.String()}
}

func randomCase(rng *rand.Rand, maxCases int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	totalN := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 1
		totalN += n
		if totalN > 200000 {
			break
		}
		p := randPerm(rng, n)
		instances = append(instances, testInstance{p: p})
	}
	if len(instances) == 0 {
		instances = append(instances, testInstance{p: []int{1}})
	}
	return makeTest(instances)
}

func randPerm(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func limitCase() testCase {
	n := 200000
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	return makeTest([]testInstance{{p: p}})
}
