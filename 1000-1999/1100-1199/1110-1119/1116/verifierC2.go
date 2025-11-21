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
	"strings"
	"time"
)

type testCase struct {
	n int
	k int64
	p []int64
}

type testInput struct {
	cases []testCase
}

func buildReferenceBinary() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("unable to locate verifier")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-1116C2-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1116C2")
	cmd := exec.Command("go", "build", "-o", binPath, "1116C2.go")
	cmd.Dir = dir
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func inputString(ti testInput) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ti.cases)))
	for _, c := range ti.cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.k))
		for i, v := range c.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expectedCount int) ([]int, error) {
	reader := strings.NewReader(out)
	res := make([]int, 0, expectedCount)
	for len(res) < expectedCount {
		var x int
		_, err := fmt.Fscan(reader, &x)
		if err != nil {
			return nil, fmt.Errorf("expected %d integers, got %d (error: %v)\noutput:\n%s", expectedCount, len(res), err, out)
		}
		res = append(res, x)
	}
	return res, nil
}

func solveCase(tc testCase) int {
	arr := append([]int64(nil), tc.p...)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	left, right := 0, tc.n
	for left < right {
		mid := (left + right + 1) / 2
		var cost int64
		for i := 1; i <= mid; i++ {
			idx := tc.n - mid + i - 1
			need := int64(i) - arr[idx]
			if need > 0 {
				cost += need
				if cost > tc.k {
					break
				}
			}
		}
		if cost <= tc.k {
			left = mid
		} else {
			right = mid - 1
		}
	}
	return left
}

func solveInput(ti testInput) []int {
	ans := make([]int, len(ti.cases))
	for i, c := range ti.cases {
		ans[i] = solveCase(c)
	}
	return ans
}

func formatAnswers(ans []int) string {
	var sb strings.Builder
	for _, v := range ans {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String()
}

func deterministicTests() []testInput {
	return []testInput{
		{
			cases: []testCase{
				{n: 1, k: 0, p: []int64{1}},
			},
		},
		{
			cases: []testCase{
				{n: 3, k: 3, p: []int64{1, 2, 2}},
				{n: 5, k: 10, p: []int64{5, 4, 3, 2, 1}},
			},
		},
	}
}

func randomCase(rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN-1) + 1
	k := rng.Int63n(200)
	p := make([]int64, n)
	for i := range p {
		p[i] = rng.Int63n(20)
	}
	return testCase{n: n, k: k, p: p}
}

func randomLargeCase(rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN-1) + 1
	k := rng.Int63n(1_000_000_000)
	p := make([]int64, n)
	for i := range p {
		p[i] = rng.Int63n(1_000_000_000)
	}
	return testCase{n: n, k: k, p: p}
}

func genTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := deterministicTests()
	for i := 0; i < 60; i++ {
		tc := testInput{cases: []testCase{}}
		caseCount := rng.Intn(4) + 1
		for j := 0; j < caseCount; j++ {
			tc.cases = append(tc.cases, randomCase(rng, 8))
		}
		tests = append(tests, tc)
	}
	for i := 0; i < 60; i++ {
		tc := testInput{cases: []testCase{}}
		caseCount := rng.Intn(3) + 1
		for j := 0; j < caseCount; j++ {
			tc.cases = append(tc.cases, randomLargeCase(rng, 50))
		}
		tests = append(tests, tc)
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for idx, ti := range tests {
		input := inputString(ti)
		expected := solveInput(ti)
		expectedOutput := formatAnswers(expected)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutput(refOut, len(expected))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}
		for i, v := range refVals {
			if v != expected[i] {
				fmt.Fprintf(os.Stderr, "reference mismatch on test %d case %d: expected %d got %d\ninput:\n%sreference output:\n%scomputed:\n%s", idx+1, i+1, expected[i], v, input, refOut, expectedOutput)
				os.Exit(1)
			}
		}

		out, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(out, len(expected))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s", idx+1, err, input, out)
			os.Exit(1)
		}
		for i, v := range gotVals {
			if v != expected[i] {
				fmt.Fprintf(os.Stderr, "test %d failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%suser output:\n%s", idx+1, i+1, expected[i], v, input, refOut, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
