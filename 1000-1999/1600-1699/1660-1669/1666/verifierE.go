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

type segment struct {
	s int64
	f int64
}

type testCase struct {
	name  string
	input string
	l     int64
	n     int
	a     []int64
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
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("reference failed on case %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}
		refSegs, err := parseSegments(refOut, tc.n)
		if err != nil {
			fmt.Printf("reference output invalid on case %d (%s): %v\nraw output:\n%s\n", i+1, tc.name, err, refOut)
			os.Exit(1)
		}
		diff, err := evaluateSolution(refSegs, tc, -1)
		if err != nil {
			fmt.Printf("reference solution invalid on case %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Printf("case %d (%s): runtime error: %v\n", i+1, tc.name, err)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}
		candSegs, err := parseSegments(candOut, tc.n)
		if err != nil {
			fmt.Printf("case %d (%s): invalid output: %v\nraw output:\n%s\n", i+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if _, err := evaluateSolution(candSegs, tc, diff); err != nil {
			fmt.Printf("case %d (%s) failed: %v\n", i+1, tc.name, err)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func referencePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "1000-1999/1600-1699/1660-1669/1666/1666E.go"
	}
	return filepath.Join(filepath.Dir(file), "1666E.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier1666E")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref1666e")
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

func parseSegments(out string, n int) ([]segment, error) {
	tokens := strings.Fields(out)
	expected := n * 2
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(tokens))
	}
	segs := make([]segment, n)
	for i := 0; i < n; i++ {
		s, err := strconv.ParseInt(tokens[2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[2*i])
		}
		f, err := strconv.ParseInt(tokens[2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[2*i+1])
		}
		segs[i] = segment{s: s, f: f}
	}
	return segs, nil
}

func evaluateSolution(segs []segment, tc testCase, requiredDiff int64) (int64, error) {
	if len(segs) != tc.n {
		return 0, fmt.Errorf("expected %d segments, got %d", tc.n, len(segs))
	}
	if tc.n == 0 {
		return 0, fmt.Errorf("no citizens")
	}
	if segs[0].s != 0 {
		return 0, fmt.Errorf("first segment must start at 0")
	}
	if segs[tc.n-1].f != tc.l {
		return 0, fmt.Errorf("last segment must end at %d", tc.l)
	}
	minLen := int64(1 << 62)
	maxLen := int64(-1)
	for i := 0; i < tc.n; i++ {
		s := segs[i].s
		f := segs[i].f
		if s < 0 || f > tc.l {
			return 0, fmt.Errorf("segment %d out of bounds", i+1)
		}
		if s >= f {
			return 0, fmt.Errorf("segment %d has non-positive length", i+1)
		}
		if i > 0 && segs[i-1].f != s {
			return 0, fmt.Errorf("segments %d and %d are not contiguous", i, i+1)
		}
		length := f - s
		if length < minLen {
			minLen = length
		}
		if length > maxLen {
			maxLen = length
		}
		if tc.a[i] < s || tc.a[i] > f {
			return 0, fmt.Errorf("segment %d does not contain home coordinate %d", i+1, tc.a[i])
		}
	}
	if minLen <= 0 {
		return 0, fmt.Errorf("segments must have positive length")
	}
	diff := maxLen - minLen
	if requiredDiff >= 0 && diff != requiredDiff {
		return 0, fmt.Errorf("segment length difference %d, expected %d", diff, requiredDiff)
	}
	return diff, nil
}

func buildTests() []testCase {
	tests := []testCase{}
	tests = append(tests, sampleTest())
	tests = append(tests, customTest("single-citizen", 12, []int64{5}))
	tests = append(tests, customTest("two-citizens", 10, []int64{2, 7}))
	tests = append(tests, customTest("clustered", 25, []int64{1, 2, 4, 8, 15, 20, 24}))

	rng := rand.New(rand.NewSource(0x51c4e2a7))
	for i := 0; i < 5; i++ {
		n := 3 + rng.Intn(7)
		l := int64(n + 5 + rng.Intn(50))
		tests = append(tests, randomTest(fmt.Sprintf("random-small-%d", i+1), n, l, rng))
	}
	for i := 0; i < 3; i++ {
		n := 50 + rng.Intn(100)
		l := int64(n + 50 + rng.Intn(500))
		tests = append(tests, randomTest(fmt.Sprintf("random-medium-%d", i+1), n, l, rng))
	}
	tests = append(tests, randomTest("random-large-1", 1000, int64(5000+rng.Intn(5000)), rng))
	tests = append(tests, randomTest("random-large-2", 5000, int64(20000+rng.Intn(10000)), rng))
	tests = append(tests, randomTest("random-huge", 20000, int64(80000+rng.Intn(40000)), rng))

	return tests
}

func sampleTest() testCase {
	return customTest("sample", 6, []int64{1, 3, 5})
}

func customTest(name string, l int64, homes []int64) testCase {
	cp := append([]int64(nil), homes...)
	return testCase{
		name:  name,
		input: formatInput(l, cp),
		l:     l,
		n:     len(cp),
		a:     cp,
	}
}

func randomTest(name string, n int, l int64, rng *rand.Rand) testCase {
	if l <= int64(n) {
		l = int64(n) + 1
	}
	homes := make([]int64, n)
	prev := int64(0)
	for i := 0; i < n; i++ {
		minVal := prev + 1
		maxVal := l - int64(n-i)
		if maxVal < minVal {
			maxVal = minVal
		}
		val := minVal
		if diff := maxVal - minVal; diff > 0 {
			val = minVal + rng.Int63n(diff+1)
		}
		homes[i] = val
		prev = val
	}
	return customTest(name, l, homes)
}

func formatInput(l int64, homes []int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", l, len(homes)))
	for i, v := range homes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
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
	return "", fmt.Errorf("usage: go run verifierE.go /path/to/solution")
}
