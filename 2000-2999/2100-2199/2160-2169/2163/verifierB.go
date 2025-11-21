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
	"time"
)

type caseData struct {
	n int
	p []int
	x string
}

type testCase struct {
	cases []caseData
	input string
}

type refAnswer struct {
	possible bool
}

type pair struct {
	l int
	r int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseReferenceOutput(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", i+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		if err := verifyTest(tc, refAns, gotOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n", i+1, err, tc.input, refOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2163B_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2163B.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier path")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseReferenceOutput(out string, caseCount int) ([]refAnswer, error) {
	fields := strings.Fields(out)
	idx := 0
	ans := make([]refAnswer, caseCount)
	for i := 0; i < caseCount; i++ {
		if idx >= len(fields) {
			return nil, fmt.Errorf("case %d: incomplete reference output", i+1)
		}
		token := fields[idx]
		idx++
		if token == "-1" {
			ans[i] = refAnswer{possible: false}
			continue
		}
		k, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("case %d: invalid operation count %q", i+1, token)
		}
		if k < 0 || k > 5 {
			return nil, fmt.Errorf("case %d: reference used invalid number of operations %d", i+1, k)
		}
		if idx+2*k > len(fields) {
			return nil, fmt.Errorf("case %d: reference output truncated", i+1)
		}
		idx += 2 * k
		ans[i] = refAnswer{possible: true}
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("reference output has %d extra tokens", len(fields)-idx)
	}
	return ans, nil
}

func verifyTest(tc testCase, refs []refAnswer, out string) error {
	if len(refs) != len(tc.cases) {
		return fmt.Errorf("internal error: reference cases mismatch")
	}
	fields := strings.Fields(out)
	idx := 0
	for caseIdx, data := range tc.cases {
		ref := refs[caseIdx]
		if !ref.possible {
			if idx >= len(fields) {
				return fmt.Errorf("case %d: expected -1 but output ended early", caseIdx+1)
			}
			if fields[idx] != "-1" {
				return fmt.Errorf("case %d: expected -1 but got %s", caseIdx+1, fields[idx])
			}
			idx++
			continue
		}

		if idx >= len(fields) {
			return fmt.Errorf("case %d: missing operation count", caseIdx+1)
		}
		k, err := strconv.Atoi(fields[idx])
		if err != nil {
			return fmt.Errorf("case %d: invalid operation count %q", caseIdx+1, fields[idx])
		}
		idx++
		if k < 0 || k > 5 {
			return fmt.Errorf("case %d: number of operations must be between 0 and 5", caseIdx+1)
		}
		if idx+2*k > len(fields) {
			return fmt.Errorf("case %d: expected %d operation pairs but output ended early", caseIdx+1, k)
		}
		ops := make([]pair, k)
		for j := 0; j < k; j++ {
			l, err1 := strconv.Atoi(fields[idx])
			r, err2 := strconv.Atoi(fields[idx+1])
			if err1 != nil || err2 != nil {
				return fmt.Errorf("case %d: invalid operation bounds", caseIdx+1)
			}
			idx += 2
			ops[j] = pair{l: l, r: r}
		}
		if err := validateOperations(data, ops); err != nil {
			return fmt.Errorf("case %d: %v", caseIdx+1, err)
		}
	}
	if idx != len(fields) {
		return fmt.Errorf("output has %d extra tokens", len(fields)-idx)
	}
	return nil
}

func validateOperations(data caseData, ops []pair) error {
	n := data.n
	s := make([]bool, n)
	for idx, op := range ops {
		l, r := op.l, op.r
		if l < 1 || l > r || r > n {
			return fmt.Errorf("operation %d has invalid bounds %d %d", idx+1, l, r)
		}
		lo := data.p[l-1]
		hi := data.p[r-1]
		if lo > hi {
			lo, hi = hi, lo
		}
		for i := l + 1; i <= r-1; i++ {
			val := data.p[i-1]
			if lo < val && val < hi {
				s[i-1] = true
			}
		}
	}
	for i := 0; i < n; i++ {
		if data.x[i] == '1' && !s[i] {
			return fmt.Errorf("position %d remains 0 but x has 1", i+1)
		}
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 40, 4, 25)...)
	tests = append(tests, randomTests(rng, 40, 5, 200)...)
	tests = append(tests, randomTests(rng, 30, 6, 2000)...)
	tests = append(tests, largeStressTests(rng)...)
	return tests
}

func manualTests() []testCase {
	var tests []testCase
	tests = append(tests, buildTestCase([]caseData{
		newCase([]int{1, 2, 3}, "010"),
	}))
	tests = append(tests, buildTestCase([]caseData{
		newCase([]int{1, 2, 3, 4, 5}, "00000"),
		newCase([]int{5, 4, 3, 2, 1}, "11111"),
	}))
	tests = append(tests, buildTestCase([]caseData{
		newCase([]int{1, 3, 2, 4, 6, 5}, "001100"),
		newCase([]int{1, 6, 2, 5, 3, 4}, "101010"),
		newCase([]int{2, 4, 1, 5, 3}, "01010"),
	}))
	tests = append(tests, buildTestCase([]caseData{
		newCase([]int{3, 1, 4, 2, 5, 6, 7}, "0011100"),
	}))
	return tests
}

func randomTests(rng *rand.Rand, batches, maxCases, maxN int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCount := rng.Intn(maxCases) + 1
		var cases []caseData
		sumN := 0
		for len(cases) < caseCount {
			if sumN+3 > 200000 {
				break
			}
			limit := maxN
			if limit > 200000-sumN {
				limit = 200000 - sumN
			}
			if limit < 3 {
				break
			}
			n := rng.Intn(limit-3+1) + 3
			cases = append(cases, randomCase(rng, n))
			sumN += n
		}
		if len(cases) == 0 {
			cases = append(cases, randomCase(rng, 3))
		}
		tests = append(tests, buildTestCase(cases))
	}
	return tests
}

func largeStressTests(rng *rand.Rand) []testCase {
	var tests []testCase
	tests = append(tests, buildTestCase([]caseData{identityCase(200000, rng)}))
	tests = append(tests, buildTestCase([]caseData{randomCase(rng, 200000)}))
	tests = append(tests, buildTestCase([]caseData{reverseCase(200000, rng)}))
	return tests
}

func newCase(p []int, x string) caseData {
	if len(p) != len(x) {
		panic("length mismatch between p and x")
	}
	cp := append([]int(nil), p...)
	return caseData{
		n: len(p),
		p: cp,
		x: x,
	}
}

func buildTestCase(cases []caseData) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	stored := make([]caseData, len(cases))
	for i, c := range cases {
		sb.WriteString(strconv.Itoa(c.n))
		sb.WriteByte('\n')
		for j, val := range c.p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		sb.WriteString(c.x)
		sb.WriteByte('\n')
		stored[i] = caseData{
			n: c.n,
			p: append([]int(nil), c.p...),
			x: c.x,
		}
	}
	return testCase{cases: stored, input: sb.String()}
}

func randomCase(rng *rand.Rand, n int) caseData {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})
	var x string
	mode := rng.Intn(4)
	switch mode {
	case 0:
		x = randomBinaryString(rng, n, 30)
	case 1:
		l := rng.Intn(n)
		r := rng.Intn(n-l) + l
		builder := make([]byte, n)
		for i := 0; i < n; i++ {
			if i >= l && i <= r && r-l >= 2 {
				builder[i] = '1'
			} else {
				builder[i] = '0'
			}
		}
		x = string(builder)
	case 2:
		builder := make([]byte, n)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				builder[i] = '1'
			} else {
				builder[i] = '0'
			}
		}
		x = string(builder)
	default:
		x = randomBinaryString(rng, n, 10)
	}
	return caseData{n: n, p: p, x: x}
}

func identityCase(n int, rng *rand.Rand) caseData {
	p := make([]int, n)
	for i := range p {
		p[i] = i + 1
	}
	x := randomBinaryString(rng, n, 20)
	return caseData{n: n, p: p, x: x}
}

func reverseCase(n int, rng *rand.Rand) caseData {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = n - i
	}
	x := randomBinaryString(rng, n, 35)
	return caseData{n: n, p: p, x: x}
}

func randomBinaryString(rng *rand.Rand, n int, percent int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	bytes := make([]byte, n)
	ones := 0
	for i := 0; i < n; i++ {
		if rng.Intn(100) < percent {
			bytes[i] = '1'
			ones++
		} else {
			bytes[i] = '0'
		}
	}
	if ones == 0 && n > 0 && percent > 0 {
		pos := rng.Intn(n)
		bytes[pos] = '1'
	}
	return string(bytes)
}
