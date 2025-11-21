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

const refSource = "2063B.go"

type singleCase struct {
	n   int
	l   int
	r   int
	arr []int64
}

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "oracle-2063B-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleB")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, expected int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	res := make([]int64, expected)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		sampleTest(),
		edgeSingleElement(),
		allInsideSegment(),
		allOutsideBetter(),
		mixedLargeSmall(),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTest() testCase {
	input := `6
2 1 1
2 1
3 2 3
1 2 3
3 1 3
3 1 3
4 2 3
1 2 2 2
5 2 5
3 3 2 3 5
6 1 3
3 6 6 4 3 2
`
	return testCase{name: "sample", input: input, outputs: 6}
}

func edgeSingleElement() testCase {
	cases := []singleCase{
		{n: 1, l: 1, r: 1, arr: []int64{10}},
		{n: 2, l: 2, r: 2, arr: []int64{5, 1}},
		{n: 3, l: 1, r: 1, arr: []int64{7, 9, 8}},
	}
	return buildTest("edge_single", cases)
}

func allInsideSegment() testCase {
	cases := []singleCase{
		{n: 4, l: 1, r: 4, arr: []int64{4, 3, 2, 1}},
		{n: 5, l: 1, r: 5, arr: []int64{5, 4, 3, 2, 1}},
	}
	return buildTest("all_inside", cases)
}

func allOutsideBetter() testCase {
	cases := []singleCase{
		{n: 6, l: 3, r: 4, arr: []int64{100, 100, 5, 6, 1, 2}},
		{n: 5, l: 2, r: 4, arr: []int64{10, 9, 8, 7, 11}},
	}
	return buildTest("outside_better", cases)
}

func mixedLargeSmall() testCase {
	cases := []singleCase{
		{n: 7, l: 2, r: 6, arr: []int64{1, 100, 2, 90, 3, 80, 4}},
		{n: 8, l: 3, r: 5, arr: []int64{50, 1, 2, 3, 4, 5, 6, 7}},
	}
	return buildTest("mixed_large_small", cases)
}

func buildTest(name string, cases []singleCase) testCase {
	var sb strings.Builder
	sb.Grow(len(cases) * 32)
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&sb, "%d %d %d\n", cs.n, cs.l, cs.r)
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    name,
		input:   sb.String(),
		outputs: len(cases),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(6) + 1
	cases := make([]singleCase, caseCnt)
	totalN := 0
	for i := 0; i < caseCnt; i++ {
		n := rng.Intn(80) + 1
		// keep total length moderate for fast checking
		if totalN+n > 50000 {
			n = 50000 - totalN
			if n == 0 {
				break
			}
		}
		totalN += n
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Int63n(1_000_000_000) + 1
		}
		// Build some structure: sometimes sorted, sometimes skewed
		switch rng.Intn(3) {
		case 0:
			// make inside large values to encourage swaps
			for j := l - 1; j < r; j++ {
				arr[j] += 500_000_000
			}
		case 1:
			// make outside large to discourage swaps
			for j := 0; j < n; j++ {
				if j < l-1 || j >= r {
					arr[j] += 500_000_000
				}
			}
		default:
		}
		cases[i] = singleCase{n: n, l: l, r: r, arr: arr}
	}
	return buildTest(fmt.Sprintf("random_%d", idx), cases)
}
