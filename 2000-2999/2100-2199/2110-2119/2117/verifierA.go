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

const refSource = "2117A.go"

type testCase struct {
	name    string
	input   string
	outputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
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

	tmpDir, err := os.MkdirTemp("", "oracle-2117A-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleA")

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
		sampleTests(),
		exhaustiveSmall(),
		simpleExtreme(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func sampleTests() testCase {
	input := `7
4 2
0 1 1 0
6 3
1 0 1 1 0 0
8 8
1 1 1 0 0 1 1 1
1 2
1
5 1
1 0 1 0 1
7 4
0 0 0 1 1 0 1
10 3
0 1 0 0 1 0 0 1 0 0
`
	return testCase{name: "sample_like", input: input, outputs: 7}
}

func exhaustiveSmall() testCase {
	var cases []string
	for n := 1; n <= 4; n++ {
		total := 1 << n
		for mask := 1; mask < total; mask++ { // need at least one closed door
			var arr []byte
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					arr = append(arr, '1')
				} else {
					arr = append(arr, '0')
				}
			}
			for _, x := range []int{1, 2, 3, 4, 5} {
				var sb strings.Builder
				fmt.Fprintf(&sb, "%d %d\n", n, x)
				for i, v := range arr {
					if i > 0 {
						sb.WriteByte(' ')
					}
					sb.WriteByte(v)
				}
				sb.WriteByte('\n')
				cases = append(cases, sb.String())
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, cs := range cases {
		sb.WriteString(cs)
	}
	return testCase{name: "exhaustive_small", input: sb.String(), outputs: len(cases)}
}

func simpleExtreme() testCase {
	cases := []struct {
		n int
		x int
		a []int
	}{
		{n: 10, x: 10, a: []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		{n: 10, x: 1, a: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{n: 10, x: 5, a: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&sb, "%d %d\n", cs.n, cs.x)
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: "extreme", input: sb.String(), outputs: len(cases)}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(25) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		x := rng.Intn(10) + 1
		arr := make([]int, n)
		hasClosed := false
		for j := 0; j < n; j++ {
			if rng.Intn(3) == 0 {
				arr[j] = 1
				hasClosed = true
			} else {
				arr[j] = 0
			}
		}
		if !hasClosed {
			pos := rng.Intn(n)
			arr[pos] = 1
		}
		fmt.Fprintf(&sb, "%d %d\n", n, x)
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		outputs: t,
	}
}
