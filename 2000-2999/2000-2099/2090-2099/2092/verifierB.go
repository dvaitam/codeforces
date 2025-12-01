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
	"time"
)

const refSource = "2092B.go"

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

	tmpDir, err := os.MkdirTemp("", "oracle-2092B-")
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
		simpleNoCase(),
		edgeAllZeros(),
		largeAlternating(),
		smallMixed(),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTest() testCase {
	input := `5
3
000
000
6
010001
010111
5
10000
10101
2
01
10
4
1100
1100
`
	return testCase{name: "sample_like", input: input, outputs: 5}
}

func simpleNoCase() testCase {
	input := `2
2
01
10
2
10
01
`
	return testCase{name: "simple_no", input: input, outputs: 2}
}

func edgeAllZeros() testCase {
	input := `2
4
0000
1111
5
00000
00000
`
	return testCase{name: "edge_zeros", input: input, outputs: 2}
}

func largeAlternating() testCase {
	n := 200000
	var a strings.Builder
	var b strings.Builder
	a.Grow(n)
	b.Grow(n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			a.WriteByte('1')
			b.WriteByte('0')
		} else {
			a.WriteByte('0')
			b.WriteByte('1')
		}
	}
	input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, a.String(), b.String())
	return testCase{name: "large_alternating", input: input, outputs: 1}
}

func smallMixed() testCase {
	input := `3
3
101
010
4
1111
0000
4
1001
0110
`
	return testCase{name: "small_mixed", input: input, outputs: 3}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	caseCnt := rng.Intn(20) + 1
	remaining := 200000
	cases := make([]string, 0, caseCnt)
	for i := 0; i < caseCnt; i++ {
		if remaining < 2 {
			break
		}
		n := rng.Intn(2000) + 2
		if n > remaining {
			n = remaining
		}
		remaining -= n
		a := make([]byte, n)
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(5) == 0 {
				a[j] = '1'
			} else {
				a[j] = '0'
			}
			if rng.Intn(5) == 0 {
				b[j] = '1'
			} else {
				b[j] = '0'
			}
		}
		var part strings.Builder
		fmt.Fprintf(&part, "%d\n%s\n%s\n", n, string(a), string(b))
		cases = append(cases, part.String())
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, cs := range cases {
		sb.WriteString(cs)
	}

	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: len(cases),
	}
}
