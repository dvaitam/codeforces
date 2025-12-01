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

const refSource = "./2141B.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		want, err := parseAnswers(refOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		got, err := parseAnswers(candOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(got) != len(want) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d answers, got %d\ninput:\n%s\n", idx+1, tc.name, len(want), len(got), tc.input)
			os.Exit(1)
		}
		for i := range want {
			if got[i] != want[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) case %d: expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, i+1, want[i], got[i], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2141B-ref-*")
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

func parseAnswers(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) < t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	if len(fields) > t {
		return nil, fmt.Errorf("expected %d numbers, got %d (extra output)", t, len(fields))
	}
	res := make([]int, t)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
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
		buildCase("tiny", []caseData{{n: 1, m: 1, a: []int{42}, b: []int{42}}}),
		buildCase("sample-like", []caseData{
			{n: 2, m: 3, a: []int{1, 2}, b: []int{3, 5}},
			{n: 5, m: 4, a: []int{1, 3, 4, 7, 10}, b: []int{1, 3, 4, 7}},
		}),
	}

	rng := rand.New(rand.NewSource(21410220))
	for i := 0; i < 80; i++ {
		t := rng.Intn(5) + 1
		cases := make([]caseData, 0, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(100) + 1
			m := rng.Intn(100) + 1
			common := rng.Intn(100) + 1
			a := randomSortedSetWithCommon(rng, n, common)
			b := randomSortedSetWithCommon(rng, m, common)
			cases = append(cases, caseData{n: n, m: m, a: a, b: b})
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), cases))
	}
	return tests
}

type caseData struct {
	n, m int
	a    []int
	b    []int
}

func buildCase(name string, data []caseData) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(data))
	for _, c := range data {
		fmt.Fprintf(&sb, "%d %d\n", c.n, c.m)
		for i, v := range c.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i, v := range c.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomSortedSetWithCommon(rng *rand.Rand, size int, common int) []int {
	set := make(map[int]struct{})
	set[common] = struct{}{}
	for len(set) < size {
		set[rng.Intn(100)+1] = struct{}{}
	}
	arr := make([]int, 0, size)
	for v := range set {
		arr = append(arr, v)
	}
	sortInts(arr)
	return arr
}

func sortInts(a []int) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
