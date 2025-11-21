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

type testFile struct {
	name  string
	input string
	cases []testQuery
}

type testQuery struct {
	n int64
	k int
}

const mod int64 = 1_000_000_007

func main() {
	candidate, err := candidatePath()
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
	for idx, tf := range tests {
		refOut, err := runProgram(refBin, tf.input)
		if err != nil {
			fmt.Printf("reference failed on file %d (%s): %v\n", idx+1, tf.name, err)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, len(tf.cases))
		if err != nil {
			fmt.Printf("reference output invalid on %s: %v\nraw output:\n%s\n", tf.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tf.input)
		if err != nil {
			fmt.Printf("case %d (%s): runtime error: %v\n", idx+1, tf.name, err)
			fmt.Println(previewInput(tf.input))
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, len(tf.cases))
		if err != nil {
			fmt.Printf("case %d (%s): invalid output format: %v\nraw output:\n%s\n", idx+1, tf.name, err, candOut)
			fmt.Println(previewInput(tf.input))
			os.Exit(1)
		}

		for caseIdx, q := range tf.cases {
			exp := refAns[caseIdx]
			got := candAns[caseIdx]
			if got != exp {
				fmt.Printf("file %d (%s) failed case %d: expected %d, got %d\n", idx+1, tf.name, caseIdx+1, exp, got)
				fmt.Println(previewInput(tf.input))
				os.Exit(1)
			}
			// Additional validation: ensure bingo result matches rank* n mod.
			if err := validateAnswer(got, q); err != nil {
				fmt.Printf("file %d (%s) case %d failed validation: %v\n", idx+1, tf.name, caseIdx+1, err)
				fmt.Println(previewInput(tf.input))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d test files passed\n", len(tests))
}

func candidatePath() (string, error) {
	if len(os.Args) == 2 {
		return os.Args[1], nil
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], nil
	}
	return "", fmt.Errorf("usage: go run verifierF.go /path/to/solution")
}

func referencePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "2000-2999/2000-2099/2030-2039/2033/2033F.go"
	}
	return filepath.Join(filepath.Dir(file), "2033F.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2033F")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2033f")
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

func parseAnswers(out string, count int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != count {
		return nil, fmt.Errorf("expected %d integers, got %d", count, len(tokens))
	}
	ans := make([]int64, count)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		val = ((val % mod) + mod) % mod
		ans[i] = val
	}
	return ans, nil
}

func validateAnswer(ans int64, q testQuery) error {
	rank := fibonacciRank(q.k)
	expected := (rank % mod) * ((q.n%mod + mod) % mod) % mod
	if ans != expected {
		return fmt.Errorf("answer %d does not match computed value %d", ans, expected)
	}
	return nil
}

func fibonacciRank(k int) int64 {
	if k == 1 {
		return 1
	}
	a, b := 0, 1%k
	for i := 1; i <= 6*k; i++ {
		if b == 0 {
			return int64(i)
		}
		a, b = b, (a+b)%k
	}
	return 0
}

func buildTests() []testFile {
	tests := []testFile{}
	tests = append(tests, sampleTest())
	tests = append(tests, customTest("edge-small", []testQuery{
		{n: 1, k: 1},
		{n: 1, k: 2},
		{n: 2, k: 2},
		{n: 5, k: 3},
		{n: 10, k: 5},
	}))
	rng := rand.New(rand.NewSource(0x5cc3e2a9))
	tests = append(tests, randomTest("random-small", rng, 40, 10, 1_000))
	tests = append(tests, randomTest("random-medium", rng, 60, 50, 100_000))
	tests = append(tests, randomTest("random-large", rng, 80, 100, 100_000))
	tests = append(tests, stressTest())
	return tests
}

func sampleTest() testFile {
	cases := []testQuery{
		{n: 3, k: 2},
		{n: 100, k: 1},
		{n: 10_000_000_000_000, k: 1},
	}
	return makeTestFile("sample", cases)
}

func customTest(name string, queries []testQuery) testFile {
	return makeTestFile(name, queries)
}

func randomTest(name string, rng *rand.Rand, count, maxK int, maxN int64) testFile {
	queries := make([]testQuery, count)
	for i := 0; i < count; i++ {
		n := randInt64(rng, 1, maxN)
		k := 1 + rng.Intn(maxK)
		queries[i] = testQuery{n: n, k: k}
	}
	return makeTestFile(name, queries)
}

func stressTest() testFile {
	queries := []testQuery{}
	// Large n and k combos
	for k := 1; k <= 100_000; k += 4999 {
		queries = append(queries, testQuery{n: 1_000_000_000_000_000_000, k: k})
	}
	return makeTestFile("stress", queries)
}

func makeTestFile(name string, cases []testQuery) testFile {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, q := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.n, q.k))
	}
	return testFile{name: name, input: sb.String(), cases: cases}
}

func randInt64(rng *rand.Rand, low, high int64) int64 {
	if low > high {
		low, high = high, low
	}
	span := high - low + 1
	if span <= 0 {
		return low
	}
	return low + rng.Int63n(span)
}

func previewInput(in string) string {
	const limit = 400
	if len(in) <= limit {
		return "Input:\n" + in
	}
	return fmt.Sprintf("Input (first %d chars):\n%s...\n", limit, in[:limit])
}
