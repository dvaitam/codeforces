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

const refSource = "2065C1.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
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
		refTokens, err := parseOutputs(refOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candTokens, err := parseOutputs(candOut, tc.outputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for i := 0; i < tc.outputs; i++ {
			if !strings.EqualFold(refTokens[i], candTokens[i]) {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at case %d: expected %q got %q\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refTokens[i], candTokens[i], tc.input, refOut, candOut)
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

	tmpDir, err := os.MkdirTemp("", "oracle-2065C1-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleC1")

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
	cmd := commandFor(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, expected int) ([]string, error) {
	tokens := strings.Fields(output)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	return tokens, nil
}

func buildTests() []testCase {
	tests := []testCase{
		sampleTest(),
		trivialCases(),
		descendingNeedsOps(),
		allSameValue(),
		bMixedSigns(),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTest() testCase {
	input := `5
1 1
5
9
3 1
1 4 3
4
4 1
1 4 2 5
6
4 1
1 4 2 5
3
3 1
9 8 7
8
`
	return testCase{name: "sample", input: input, outputs: 5}
}

func trivialCases() testCase {
	input := `4
1 1
10
1
2 1
1 1
2
2 1
5 4
9
3 1
2 2 2
5
`
	return testCase{name: "trivial", input: input, outputs: 4}
}

func descendingNeedsOps() testCase {
	input := `3
4 1
8 7 6 5
10
5 1
5 4 3 2 1
6
6 1
6 5 4 3 2 1
7
`
	return testCase{name: "descending", input: input, outputs: 3}
}

func allSameValue() testCase {
	n := 8
	val := int64(5)
	arr := make([]string, n)
	for i := range arr {
		arr[i] = fmt.Sprint(val)
	}
	input := fmt.Sprintf("2\n%d 1\n%s\n5\n%d 1\n%s\n10\n", n, strings.Join(arr, " "), n, strings.Join(arr, " "))
	return testCase{name: "all_same", input: input, outputs: 2}
}

func bMixedSigns() testCase {
	input := `3
3 1
1 2 3
1
3 1
3 2 1
4
4 1
9 1 9 1
5
`
	return testCase{name: "b_mixed", input: input, outputs: 3}
}

func buildTestInput(cases []struct {
	n int
	a []int64
	b int64
}) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&sb, "%d 1\n", cs.n)
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", cs.b)
	}
	return sb.String()
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(6) + 1
	cases := make([]struct {
		n int
		a []int64
		b int64
	}, caseCnt)
	totalN := 0
	for i := 0; i < caseCnt; i++ {
		n := rng.Intn(15) + 1
		if totalN+n > 200000 {
			n = 200000 - totalN
			if n <= 0 {
				caseCnt = i
				break
			}
		}
		totalN += n
		a := make([]int64, n)
		base := rng.Int63n(20) - 10 // allow negative after ops
		for j := 0; j < n; j++ {
			a[j] = rng.Int63n(20) + base + 1_000_000_000/2
		}
		b := rng.Int63n(1_000_000_000) + 1
		// occasionally pick b that equals some sum to create symmetry
		if rng.Intn(4) == 0 {
			b = a[rng.Intn(n)] + a[rng.Intn(n)]
		}
		cases[i] = struct {
			n int
			a []int64
			b int64
		}{n: n, a: a, b: b}
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   buildTestInput(cases[:caseCnt]),
		outputs: caseCnt,
	}
}
