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

const refSource = "2082B.go"

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

		for i := 0; i < len(refVals); i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) at token %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
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

	tmpDir, err := os.MkdirTemp("", "oracle-2082B-")
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
	cmd := commandFor(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, cases int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != cases*2 {
		return nil, fmt.Errorf("expected %d integers, got %d", cases*2, len(tokens))
	}
	res := make([]int64, len(tokens))
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", tok, err)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		sampleTest(),
		edgeZeros(),
		allCeilsOrFloors(),
		largeCounts(),
		smallMixed(),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTest() testCase {
	input := `5
12 1 2
12 1 1
12 0 0
12 1000000000 1000000000
706636307 0 3
`
	return testCase{name: "sample", input: input, outputs: 5}
}

func edgeZeros() testCase {
	input := `5
0 0 0
0 5 7
1 0 0
1 5 0
1 0 5
`
	return testCase{name: "zeros_and_ones", input: input, outputs: 5}
}

func allCeilsOrFloors() testCase {
	input := `4
9 0 5
9 5 0
1000000000 0 30
1000000000 30 0
`
	return testCase{name: "all_ceils_floors", input: input, outputs: 4}
}

func largeCounts() testCase {
	input := `4
1000000000 1000000000 0
1000000000 0 1000000000
1000000000 1000000000 1000000000
2 1000000000 1000000000
`
	return testCase{name: "large_counts", input: input, outputs: 4}
}

func smallMixed() testCase {
	input := `5
5 3 3
6 2 3
7 3 2
8 4 1
15 1 4
`
	return testCase{name: "small_mixed", input: input, outputs: 5}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(200) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", caseCnt)
	for i := 0; i < caseCnt; i++ {
		x := rng.Int63n(1_000_000_001)
		n := rng.Int63n(1_000_000_001)
		m := rng.Int63n(1_000_000_001)
		if rng.Intn(5) == 0 {
			// make small x to force early stop
			x = int64(rng.Intn(3))
		}
		fmt.Fprintf(&sb, "%d %d %d\n", x, n, m)
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: caseCnt,
	}
}
