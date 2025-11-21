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

type singleCase struct {
	n   int
	arr []int
}

type testCase struct {
	name    string
	input   string
	answers int
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
		expected, err := parseAnswers(refOut, tc.answers)
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
		got, err := parseAnswers(candOut, tc.answers)
		if err != nil {
			fmt.Printf("case %d (%s): invalid output: %v\nraw output:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}
		for i := 0; i < tc.answers; i++ {
			if got[i] != expected[i] {
				fmt.Printf("case %d (%s) failed at test %d: expected %s got %s\n", idx+1, tc.name, i+1, expected[i], got[i])
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
		return "2000-2999/2000-2099/2010-2019/2011/2011A.go"
	}
	return filepath.Join(filepath.Dir(file), "2011A.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2011A")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2011a")
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

func parseAnswers(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		if tok == "Ambiguous" {
			res[i] = tok
			continue
		}
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid output token %q", tok)
		}
		res[i] = strconv.Itoa(val)
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{}
	tests = append(tests, sampleTest())
	tests = append(tests, customTest("edge-cases", []singleCase{
		makeCase([]int{1, 2}),
		makeCase([]int{2, 2, 3}),
		makeCase([]int{4, 4, 4, 5}),
		makeCase([]int{1, 1, 1, 1, 3}),
		makeCase([]int{10, 20, 30, 50}),
	}))

	rng := rand.New(rand.NewSource(0x4a02b1cc))
	tests = append(tests, randomBatch("random-small", rng, 25))
	tests = append(tests, randomBatch("random-medium", rng, 120))
	tests = append(tests, randomBatch("random-large", rng, 400))

	tricky := []singleCase{
		makeCase([]int{49, 49, 49, 49, 50}),
		makeCase([]int{1, 49, 50}),
		makeCase([]int{25, 25, 25, 25, 26}),
		makeCase([]int{5, 6, 7, 9, 12, 20}),
		makeCase([]int{7, 7, 7, 7, 7, 8}),
		makeCase([]int{8, 8, 9, 9, 10, 15}),
	}
	tests = append(tests, customTest("tricky-values", tricky))

	return tests
}

func sampleTest() testCase {
	cases := []singleCase{
		makeCase([]int{1, 2, 3, 4, 5}),
		makeCase([]int{8, 8, 5, 3, 4, 6, 8, 12}),
		makeCase([]int{3, 3, 3, 4}),
	}
	return customTest("sample", cases)
}

func makeCase(arr []int) singleCase {
	n := len(arr)
	if n < 2 {
		panic("case must have at least two problems")
	}
	last := arr[n-1]
	maxBefore := arr[0]
	for i := 1; i < n-1; i++ {
		if arr[i] > maxBefore {
			maxBefore = arr[i]
		}
	}
	if maxBefore >= last {
		panic("last difficulty must be strictly greater than previous ones")
	}
	return singleCase{n: n, arr: append([]int(nil), arr...)}
}

func customTest(name string, cases []singleCase) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String(), answers: len(cases)}
}

func randomBatch(name string, rng *rand.Rand, caseCount int) testCase {
	cases := make([]singleCase, caseCount)
	for i := 0; i < caseCount; i++ {
		cases[i] = randomCase(rng)
	}
	return customTest(name, cases)
}

func randomCase(rng *rand.Rand) singleCase {
	n := 2 + rng.Intn(49)
	arr := make([]int, n)
	maxBefore := 0
	for i := 0; i < n-1; i++ {
		val := 1 + rng.Intn(49)
		arr[i] = val
		if val > maxBefore {
			maxBefore = val
		}
	}
	if maxBefore == 0 {
		maxBefore = 1
	}
	last := maxBefore + 1 + rng.Intn(50-maxBefore)
	arr[n-1] = last
	return singleCase{n: n, arr: arr}
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
