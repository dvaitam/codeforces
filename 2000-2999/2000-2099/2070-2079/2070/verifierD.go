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

const refSource = "2070D.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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

	tmpDir, err := os.MkdirTemp("", "oracle-2070D-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleD")

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
		simpleChains(),
		stars(),
		mixedManual(),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTest() testCase {
	input := `3
4
1 2 1
3
1 2
7
1 2 2 1 4 5
`
	return testCase{name: "sample", input: input, outputs: 3}
}

func simpleChains() testCase {
	var sb strings.Builder
	sb.WriteString("3\n")
	// chain length 2
	sb.WriteString("2\n1\n")
	// chain length 4
	sb.WriteString("4\n1 2 3\n")
	// chain length 6
	sb.WriteString("6\n1 2 3 4 5\n")
	return testCase{name: "chains", input: sb.String(), outputs: 3}
}

func stars() testCase {
	var sb strings.Builder
	sb.WriteString("2\n")
	// star with 5 leaves
	sb.WriteString("6\n1 1 1 1 1\n")
	// star with 8 leaves
	sb.WriteString("9\n1 1 1 1 1 1 1 1\n")
	return testCase{name: "stars", input: sb.String(), outputs: 2}
}

func mixedManual() testCase {
	var sb strings.Builder
	sb.WriteString("3\n")
	// balanced-ish
	sb.WriteString("7\n1 1 2 2 3 3\n")
	// skewed with back-and-forth depths
	sb.WriteString("8\n1 2 2 3 3 4 4\n")
	// random small tree
	sb.WriteString("5\n1 1 2 4\n")
	return testCase{name: "manual_mix", input: sb.String(), outputs: 3}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(5) + 1
	remaining := 200_000
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", caseCnt)
	for i := 0; i < caseCnt; i++ {
		n := rng.Intn(2000) + 2
		if n > remaining {
			n = remaining
		}
		remaining -= n
		if remaining < 0 {
			remaining = 0
		}
		parents := make([]int, n+1)
		for v := 2; v <= n; v++ {
			// build random parent among predecessors
			parents[v] = rng.Intn(v-1) + 1
		}
		fmt.Fprintf(&sb, "%d\n", n)
		for v := 2; v <= n; v++ {
			if v > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(parents[v]))
		}
		sb.WriteByte('\n')
		if remaining == 0 {
			// fill rest trivial if we ran out
			for j := i + 1; j+1 <= caseCnt; j++ {
				fmt.Fprintln(&sb, "2")
				fmt.Fprintln(&sb, "1")
			}
			break
		}
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: caseCnt,
	}
}
