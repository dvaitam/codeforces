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
	"time"
)

const refSource = "1000-1999/1900-1999/1990-1999/1998/1998C.go"

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
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d cases, got %d\ninput:\n%s\n", idx+1, tc.name, len(refVals), len(candVals), tc.input)
			os.Exit(1)
		}

		for caseIdx := range refVals {
			if refVals[caseIdx] != candVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected %d, got %d\ninput:\n%s\n", idx+1, caseIdx+1, refVals[caseIdx], candVals[caseIdx], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1998C-ref-*")
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

func parseOutputs(input, out string) ([]int64, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}

	lines := filterNonEmpty(strings.Split(out, "\n"))
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(lines))
	}
	results := make([]int64, t)
	for i, line := range lines {
		val, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", line, err)
		}
		results[i] = val
	}
	return results, nil
}

func filterNonEmpty(lines []string) []string {
	var res []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

type caseData struct {
	arr   []int64
	flags []int
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("small", []caseData{
			{arr: []int64{1, 2, 3}, flags: []int{1, 0, 0}},
		}),
		buildCase("all-allowed", []caseData{
			{arr: []int64{2, 4, 6, 8}, flags: []int{1, 1, 1, 1}},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		t := rng.Intn(5) + 1
		cases := make([]caseData, 0, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(5) + 1
			arr := make([]int64, n)
			flags := make([]int, n)
			for idx := range arr {
				arr[idx] = randRange(rng, 1, 20)
				flags[idx] = rng.Intn(2)
			}
			cases = append(cases, caseData{arr: arr, flags: flags})
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), cases))
	}
	return tests
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}

func buildCase(name string, data []caseData) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(data))
	for _, d := range data {
		arr := d.arr
		flags := d.flags
		fmt.Fprintf(&sb, "%d %d\n", len(arr), len(flags))
		for idx, v := range arr {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for idx, v := range flags {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}
