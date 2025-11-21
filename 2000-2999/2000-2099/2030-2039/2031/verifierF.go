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

type caseData struct {
	n    int
	perm []int
}

type testFile struct {
	name  string
	input string
	cases []caseData
}

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
		refAnswers, err := parseAnswers(refOut, tf.cases)
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
		candAnswers, err := parseAnswers(candOut, tf.cases)
		if err != nil {
			fmt.Printf("case %d (%s): invalid output format: %v\nraw output:\n%s\n", idx+1, tf.name, err, candOut)
			fmt.Println(previewInput(tf.input))
			os.Exit(1)
		}
		for caseIdx, cs := range tf.cases {
			refAns := refAnswers[caseIdx]
			candAns := candAnswers[caseIdx]
			if err := validateAnswer(refAns, cs); err != nil {
				fmt.Printf("reference answer invalid on %s case %d: %v\n", tf.name, caseIdx+1, err)
				os.Exit(1)
			}
			if err := validateAnswer(candAns, cs); err != nil {
				fmt.Printf("file %d (%s) failed case %d: %v\n", idx+1, tf.name, caseIdx+1, err)
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
		return "2000-2999/2000-2099/2030-2039/2031/2031F.go"
	}
	return filepath.Join(filepath.Dir(file), "2031F.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2031F")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2031f")
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

func parseAnswers(out string, cases []caseData) ([][2]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != len(cases)*2 {
		return nil, fmt.Errorf("expected %d integers, got %d", len(cases)*2, len(tokens))
	}
	ans := make([][2]int, len(cases))
	idx := 0
	for i := range cases {
		x, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[idx])
		}
		y, err := strconv.Atoi(tokens[idx+1])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[idx+1])
		}
		ans[i] = [2]int{x, y}
		idx += 2
	}
	return ans, nil
}

func validateAnswer(ans [2]int, cs caseData) error {
	n := cs.n
	i1, i2 := ans[0], ans[1]
	if i1 < 1 || i1 > n || i2 < 1 || i2 > n {
		return fmt.Errorf("indices must be between 1 and %d, got %d %d", n, i1, i2)
	}
	target1 := n / 2
	target2 := target1 + 1
	v1 := cs.perm[i1-1]
	v2 := cs.perm[i2-1]
	if (v1 == target1 && v2 == target2) || (v1 == target2 && v2 == target1) {
		return nil
	}
	return fmt.Errorf("indices %d %d do not point to required values", i1, i2)
}

func buildTests() []testFile {
	tests := []testFile{}
	tests = append(tests, sampleTest())
	tests = append(tests, customTest("small-cases", []caseData{
		{n: 6, perm: []int{6, 2, 3, 5, 1, 4}},
		{n: 6, perm: []int{1, 2, 3, 4, 5, 6}},
		{n: 8, perm: []int{8, 7, 6, 5, 4, 3, 2, 1}},
	}))
	rng := rand.New(rand.NewSource(0x51d34ac2))
	tests = append(tests, randomTest("random-small", rng, 20, 10))
	tests = append(tests, randomTest("random-medium", rng, 15, 40))
	tests = append(tests, randomTest("random-large", rng, 10, 100))
	return tests
}

func sampleTest() testFile {
	cases := []caseData{
		{n: 6, perm: []int{6, 2, 3, 5, 1, 4}},
		{n: 10, perm: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}},
	}
	return makeTestFile("sample", cases)
}

func customTest(name string, cases []caseData) testFile {
	return makeTestFile(name, cases)
}

func randomTest(name string, rng *rand.Rand, caseCount, maxN int) testFile {
	cases := make([]caseData, caseCount)
	for i := 0; i < caseCount; i++ {
		n := 6 + 2*rng.Intn((maxN-4)/2)
		perm := randPermutation(rng, n)
		cases[i] = caseData{n: n, perm: perm}
	}
	return makeTestFile(name, cases)
}

func randPermutation(rng *rand.Rand, n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}

func makeTestFile(name string, cases []caseData) testFile {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for i, v := range cs.perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testFile{name: name, input: sb.String(), cases: cases}
}

func previewInput(in string) string {
	const limit = 400
	if len(in) <= limit {
		return "Input:\n" + in
	}
	return fmt.Sprintf("Input (first %d chars):\n%s...\n", limit, in[:limit])
}
