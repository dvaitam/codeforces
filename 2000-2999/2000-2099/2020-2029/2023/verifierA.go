package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type pair struct {
	a int64
	b int64
}

type caseData struct {
	pairs []pair
}

type testCase struct {
	name  string
	input string
	cases []caseData
	ns    []int
}

func main() {
	candidate, err := candidatePathFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refSrc := referencePath()
	refBin, cleanup, err := buildReferenceBinary(refSrc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("reference failed on case %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		refSeqs, err := parseOutputs(refOut, tc.ns)
		if err != nil {
			fmt.Printf("reference output invalid on case %d (%s): %v\nraw output:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Printf("case %d (%s): runtime error: %v\n", idx+1, tc.name, err)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}
		candSeqs, err := parseOutputs(candOut, tc.ns)
		if err != nil {
			fmt.Printf("case %d (%s): invalid output: %v\nraw output:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}

		for caseIdx, data := range tc.cases {
			refSeq := refSeqs[caseIdx]
			candSeq := candSeqs[caseIdx]
			if err := verifyPermutation(refSeq, data.pairs); err != nil {
				fmt.Printf("reference solution invalid on %s case %d: %v\n", tc.name, caseIdx+1, err)
				os.Exit(1)
			}
			if err := verifyPermutation(candSeq, data.pairs); err != nil {
				fmt.Printf("case %d (%s) failed validation on test %d: %v\n", idx+1, tc.name, caseIdx+1, err)
				fmt.Println(previewInput(tc.input))
				os.Exit(1)
			}
			refInv := countInversions(refSeq)
			candInv := countInversions(candSeq)
			if candInv != refInv {
				fmt.Printf("case %d (%s) failed on test %d: inversion count %d expected %d\n", idx+1, tc.name, caseIdx+1, candInv, refInv)
				fmt.Println(previewInput(tc.input))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func referencePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "2000-2999/2000-2099/2020-2029/2023/2023A.go"
	}
	return filepath.Join(filepath.Dir(file), "2023A.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2023A")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2023a")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference solution: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseOutputs(out string, ns []int) ([][]int64, error) {
	tokens := strings.Fields(out)
	expected := 0
	for _, n := range ns {
		expected += 2 * n
	}
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	res := make([][]int64, len(ns))
	idx := 0
	for i, n := range ns {
		seq := make([]int64, 2*n)
		for j := 0; j < 2*n; j++ {
			val, err := strconv.ParseInt(tokens[idx], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", tokens[idx])
			}
			seq[j] = val
			idx++
		}
		res[i] = seq
	}
	return res, nil
}

func verifyPermutation(seq []int64, orig []pair) error {
	if len(seq) != 2*len(orig) {
		return fmt.Errorf("expected %d numbers, got %d", 2*len(orig), len(seq))
	}
	freq := make(map[pair]int)
	for _, p := range orig {
		freq[p]++
	}
	for i := 0; i < len(seq); i += 2 {
		p := pair{a: seq[i], b: seq[i+1]}
		if freq[p] == 0 {
			return fmt.Errorf("pair (%d,%d) not in input multiset", p.a, p.b)
		}
		freq[p]--
	}
	for _, v := range freq {
		if v != 0 {
			return fmt.Errorf("not all pairs used")
		}
	}
	return nil
}

func countInversions(arr []int64) int64 {
	if len(arr) == 0 {
		return 0
	}
	values := append([]int64(nil), arr...)
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	uniq := values[:0]
	for _, v := range values {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	comp := make(map[int64]int, len(uniq))
	for i, v := range uniq {
		comp[v] = i + 1
	}
	fen := newFenwick(len(uniq) + 2)
	inv := int64(0)
	for i := len(arr) - 1; i >= 0; i-- {
		idx := comp[arr[i]]
		inv += fen.sum(idx - 1)
		fen.add(idx, 1)
	}
	return inv
}

type fenwick struct {
	tree []int64
}

func newFenwick(n int) *fenwick {
	return &fenwick{tree: make([]int64, n+2)}
}

func (f *fenwick) add(pos int, delta int64) {
	for pos < len(f.tree) {
		f.tree[pos] += delta
		pos += pos & -pos
	}
}

func (f *fenwick) sum(pos int) int64 {
	res := int64(0)
	for pos > 0 {
		res += f.tree[pos]
		pos -= pos & -pos
	}
	return res
}

func buildTests() []testCase {
	tests := []testCase{}
	tests = append(tests, sampleTest())
	tests = append(tests, customTest("edgecases", []caseData{
		{pairs: []pair{{1, 1}}},
		{pairs: []pair{{2, 3}, {4, 5}}},
		{pairs: []pair{{5, 1}, {5, 1}}},
		{pairs: []pair{{7, 2}, {2, 7}, {3, 3}}},
	}))

	rng := rand.New(rand.NewSource(0x75a4b223))
	tests = append(tests, randomBatch("random-small", rng, 20, 5, 30))
	tests = append(tests, randomBatch("random-medium", rng, 10, 200, 1_000_000))
	tests = append(tests, randomBatch("random-large", rng, 4, 5000, 1_000_000_000))
	tests = append(tests, structuredBatch())
	return tests
}

func sampleTest() testCase {
	cases := []caseData{
		{pairs: []pair{{1, 4}, {2, 3}}},
		{pairs: []pair{{3, 2}, {4, 3}, {2, 1}}},
		{pairs: []pair{{5, 10}, {2, 3}, {9, 6}, {4, 1}, {8, 7}}},
		{pairs: []pair{{10, 20}}},
	}
	return makeTestCase("sample", cases)
}

func customTest(name string, cases []caseData) testCase {
	return makeTestCase(name, cases)
}

func randomBatch(name string, rng *rand.Rand, caseCount, maxN int, valueMax int64) testCase {
	cases := make([]caseData, caseCount)
	for i := 0; i < caseCount; i++ {
		n := 1 + rng.Intn(maxN)
		pairs := make([]pair, n)
		for j := 0; j < n; j++ {
			a := randValue(rng, valueMax)
			b := randValue(rng, valueMax)
			pairs[j] = pair{a: a, b: b}
		}
		cases[i] = caseData{pairs: pairs}
	}
	return makeTestCase(name, cases)
}

func structuredBatch() testCase {
	cases := []caseData{}
	// ascending sums
	asc := make([]pair, 0)
	for i := 1; i <= 20; i++ {
		asc = append(asc, pair{a: int64(i), b: int64(100 - i)})
	}
	cases = append(cases, caseData{pairs: asc})

	// identical sums, varying a
	eq := make([]pair, 0)
	for i := 1; i <= 30; i++ {
		eq = append(eq, pair{a: int64(i), b: int64(60 - i)})
	}
	cases = append(cases, caseData{pairs: eq})

	// repeated values
	rep := make([]pair, 0)
	for i := 0; i < 40; i++ {
		rep = append(rep, pair{a: 5, b: 5})
	}
	cases = append(cases, caseData{pairs: rep})

	// large n with random values but fixed pattern
	large := make([]pair, 6000)
	for i := 0; i < len(large); i++ {
		large[i] = pair{a: int64((i%100)+1) * 1000, b: int64((100 - i%100) + 1)}
	}
	cases = append(cases, caseData{pairs: large})
	return makeTestCase("structured", cases)
}

func makeTestCase(name string, cases []caseData) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	ns := make([]int, len(cases))
	for i, cs := range cases {
		n := len(cs.pairs)
		ns[i] = n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, p := range cs.pairs {
			if j > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(fmt.Sprintf("%d %d", p.a, p.b))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String(), cases: cases, ns: ns}
}

func randValue(rng *rand.Rand, maxVal int64) int64 {
	if maxVal <= 1 {
		return 1
	}
	return rng.Int63n(maxVal) + 1
}

func previewInput(in string) string {
	const limit = 400
	if len(in) <= limit {
		return "Input:\n" + in
	}
	return fmt.Sprintf("Input (first %d chars):\n%s...\n", limit, in[:limit])
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) == 2 {
		return os.Args[1], nil
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], nil
	}
	return "", fmt.Errorf("usage: go run verifierA.go /path/to/solution")
}
