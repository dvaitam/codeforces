package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "0-999/800-899/880-889/883/883J.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "883J.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n, m int
	a    []int64
	b    []int64
	p    []int64
}

func inputString(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseAnswer(out string) (int64, error) {
	reader := strings.NewReader(out)
	var val int64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer answer: %v\nfull output:\n%s", err, out)
	}
	return val, nil
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
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "883J-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_883J")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func randomArrayInt64(rng *rand.Rand, length int, lo, hi int64) []int64 {
	arr := make([]int64, length)
	for i := range arr {
		arr[i] = lo + rng.Int63n(hi-lo+1)
	}
	return arr
}

func randomCase(rng *rand.Rand, maxN, maxM int, maxA, maxB, maxP int64) testCase {
	n := rng.Intn(maxN) + 1
	m := rng.Intn(maxM) + 1
	a := randomArrayInt64(rng, n, 1, maxA)
	b := randomArrayInt64(rng, m, 1, maxB)
	p := randomArrayInt64(rng, m, 1, maxP)
	return testCase{n: n, m: m, a: a, b: b, p: p}
}

func structuredCase() []testCase {
	return []testCase{
		{
			n: 1, m: 1,
			a: []int64{5},
			b: []int64{3},
			p: []int64{2},
		},
		{
			n: 2, m: 2,
			a: []int64{1, 2},
			b: []int64{3, 3},
			p: []int64{5, 6},
		},
		{
			n: 3, m: 4,
			a: []int64{3, 2, 5},
			b: []int64{1, 2, 3, 4},
			p: []int64{4, 1, 2, 3},
		},
		{
			n: 5, m: 6,
			a: []int64{6, 3, 2, 4, 3},
			b: []int64{3, 6, 4, 5, 4, 2},
			p: []int64{1, 4, 3, 2, 5, 3},
		},
	}
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(20250307))
	tests := structuredCase()
	for i := 0; i < 80; i++ {
		tests = append(tests, randomCase(rng, 5, 6, 20, 20, 20))
	}
	for i := 0; i < 80; i++ {
		tests = append(tests, randomCase(rng, 50, 60, 1000, 1000, 1000))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomCase(rng, 200, 200, 1_000_000_000, 1_000_000_000, 1_000_000_000))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	if bin == "--" {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for i, tc := range tests {
		in := inputString(tc)
		refOut, err := runProgram(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		expected, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}

		out, runErr := runProgram(bin, in)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, in, out)
			os.Exit(1)
		}
		got, err := parseAnswer(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, expected, got, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
