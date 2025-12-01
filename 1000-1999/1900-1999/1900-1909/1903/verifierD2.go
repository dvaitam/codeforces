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
const refSource = "1903D2.go"

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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d answers, got %d\ninput:\n%s\n", idx+1, tc.name, len(refVals), len(candVals), tc.input)
			os.Exit(1)
		}

		for caseIdx := range refVals {
			if refVals[caseIdx] != candVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) query %d: expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, caseIdx+1, refVals[caseIdx], candVals[caseIdx], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1903D2-ref-*")
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
	var n, q int
	fmt.Fscan(reader, &n, &q)
	for i := 0; i < n; i++ {
		var tmp int64
		fmt.Fscan(reader, &tmp)
	}

	lines := filterNonEmpty(strings.Split(out, "\n"))
	if len(lines) != q {
		return nil, fmt.Errorf("expected %d outputs, got %d", q, len(lines))
	}
	res := make([]int64, q)
	for i, line := range lines {
		val, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", line, err)
		}
		res[i] = val
	}
	return res, nil
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
		buildCase("small", []int64{1, 2, 3}, []int64{0, 5, 10}),
		buildCase("single", []int64{7}, []int64{0, 1, 2, 3}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = randRange(rng, 0, 100)
		}
		qs := make([]int64, q)
		for j := range qs {
			qs[j] = randRange(rng, 0, 100)
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), arr, qs))
	}
	return tests
}

func buildCase(name string, arr []int64, queries []int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(arr), len(queries))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for _, k := range queries {
		fmt.Fprintf(&sb, "%d\n", k)
	}
	return testCase{name: name, input: sb.String()}
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}
