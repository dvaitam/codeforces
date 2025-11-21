package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2100-2199/2130-2139/2134/2134C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		t, err := readT(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal parse error on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(candAns) != len(refAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d answers, got %d\ninput:\n%s\n", idx+1, tc.name, len(refAns), len(candAns), tc.input)
			os.Exit(1)
		}
		for i := range refAns {
			if candAns[i] != refAns[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) case %d: expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, i+1, refAns[i], candAns[i], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2134C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseAnswers(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) < t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	if len(fields) > t {
		return nil, fmt.Errorf("expected %d numbers, got %d (extra output)", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d is not an integer: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func readT(input string) (int, error) {
	var t int
	_, err := fmt.Fscan(strings.NewReader(input), &t)
	return t, err
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("already-good", [][]int64{{3, 8, 4, 4}}),
		buildCase("simple-fix", [][]int64{{0, 2, 3, 5}}),
		buildCase("two-elements", [][]int64{{5, 1}}),
		buildCase("all-zero", [][]int64{{0, 0, 0, 0, 0}}),
	}

	rng := rand.New(rand.NewSource(21340315))
	for i := 0; i < 60; i++ {
		t := rng.Intn(4) + 1
		cases := make([][]int64, 0, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(30) + 2
			// occasionally test larger sizes
			if rng.Intn(10) == 0 {
				n = rng.Intn(200) + 50
			}
			arr := make([]int64, n)
			for k := 0; k < n; k++ {
				arr[k] = int64(rng.Intn(1_000_000_000))
			}
			cases = append(cases, arr)
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), cases))
	}
	return tests
}

func buildCase(name string, arrays [][]int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arrays))
	for _, arr := range arrays {
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}
