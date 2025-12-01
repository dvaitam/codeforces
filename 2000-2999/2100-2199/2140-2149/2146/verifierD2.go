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

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2146D2.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
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
		refAns, err := parseOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d test cases, got %d\ninput:\n%s\n", idx+1, tc.name, len(refAns), len(candAns), tc.input)
			os.Exit(1)
		}

		reader := strings.NewReader(tc.input)
		var t int
		fmt.Fscan(reader, &t)
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			n := r - l + 1
			vals := make([]int, n)
			for i := 0; i < n; i++ {
				vals[i] = l + i
			}
			refCase := refAns[caseIdx]
			candCase := candAns[caseIdx]

			if refCase.total != candCase.total {
				fmt.Fprintf(os.Stderr, "wrong total on test %d case %d: expected %d, got %d\n", idx+1, caseIdx+1, refCase.total, candCase.total)
				os.Exit(1)
			}
			if len(candCase.sequence) != n {
				fmt.Fprintf(os.Stderr, "wrong sequence length on test %d case %d: expected %d numbers, got %d\n", idx+1, caseIdx+1, n, len(candCase.sequence))
				os.Exit(1)
			}

			used := make(map[int]bool, n)
			for i, x := range candCase.sequence {
				if x < 0 {
					fmt.Fprintf(os.Stderr, "invalid sequence on test %d case %d: values must be non-negative, got %d at position %d\n", idx+1, caseIdx+1, x, i+1)
					os.Exit(1)
				}
				if used[x] {
					fmt.Fprintf(os.Stderr, "invalid sequence on test %d case %d: duplicate value %d\n", idx+1, caseIdx+1, x)
					os.Exit(1)
				}
				used[x] = true
			}

			total := int64(0)
			for i, x := range candCase.sequence {
				total += int64(vals[i] | x)
			}
			if total != refCase.total {
				fmt.Fprintf(os.Stderr, "wrong OR sum on test %d case %d: expected %d, got %d\n", idx+1, caseIdx+1, refCase.total, total)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2146D2-ref-*")
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

type caseAnswer struct {
	total    int64
	sequence []int
}

func parseOutput(input, out string) ([]caseAnswer, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	cases := make([]struct {
		l int
		r int
	}, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &cases[i].l, &cases[i].r)
	}
	lines := filterNonEmpty(strings.Split(out, "\n"))
	pos := 0
	ans := make([]caseAnswer, t)
	for i := 0; i < t; i++ {
		if pos >= len(lines) {
			return nil, fmt.Errorf("missing output for case %d", i+1)
		}
		total, err := strconv.ParseInt(strings.TrimSpace(lines[pos]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid total on case %d: %v", i+1, err)
		}
		pos++
		if pos >= len(lines) {
			return nil, fmt.Errorf("missing sequence on case %d", i+1)
		}
		fields := strings.Fields(lines[pos])
		pos++
		seq := make([]int, len(fields))
		for j, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q in sequence of case %d: %v", f, i+1, err)
			}
			seq[j] = val
		}
		ans[i] = caseAnswer{total: total, sequence: seq}
	}
	return ans, nil
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

func generateTests() []testCase {
	tests := []testCase{
		buildCase("small", [][2]int{{1, 2}, {3, 5}}),
		buildCase("single", [][2]int{{7, 7}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		t := rng.Intn(5) + 1
		qs := make([][2]int, t)
		for j := 0; j < t; j++ {
			l := rng.Intn(50)
			r := l + rng.Intn(10)
			qs[j] = [2]int{l, r}
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), qs))
	}
	return tests
}

func buildCase(name string, qs [][2]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(qs))
	for _, q := range qs {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	return testCase{name: name, input: sb.String()}
}
