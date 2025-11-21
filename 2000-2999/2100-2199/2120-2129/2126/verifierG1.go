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

const refSource = "2126G1.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/candidate")
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

	tmpDir, err := os.MkdirTemp("", "oracle-2126G1-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleG1")

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
		sampleTests(),
		smallEdge(),
		increasingDecreasing(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTests() testCase {
	input := `5
3
2 5 3
4
1 1 1 3
6
1 3 4 6 2 7
4
2 3 1 5
1
2
`
	return testCase{name: "sample_like", input: input, outputs: 5}
}

func smallEdge() testCase {
	input := `4
1
5
2
1 100
3
5 5 1
3
1 2 3
`
	return testCase{name: "small_edge", input: input, outputs: 4}
}

func increasingDecreasing() testCase {
	input := `3
5
1 2 3 4 5
5
5 4 3 2 1
6
1 3 2 4 3 5
`
	return testCase{name: "inc_dec", input: input, outputs: 3}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(40) + 1
		var arr strings.Builder
		for j := 0; j < n; j++ {
			if j > 0 {
				arr.WriteByte(' ')
			}
			val := rng.Intn(minInt(n, 100)) + 1
			arr.WriteString(strconv.Itoa(val))
		}
		fmt.Fprintf(&sb, "%d\n%s\n", n, arr.String())
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: t,
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
