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

const refSource = "./1270H.go"

type testCase struct {
	name  string
	input string
	q     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
		refVals, err := parseOutputs(refOut, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}
		if len(refVals) != tc.q {
			fmt.Fprintf(os.Stderr, "reference output length mismatch on test %d (%s): expected %d lines, got %d\n", idx+1, tc.name, tc.q, len(refVals))
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if len(candVals) != tc.q {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d lines, got %d\n", idx+1, tc.name, tc.q, len(candVals))
			os.Exit(1)
		}

		for i := range refVals {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) line %d: expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1270H-ref-*")
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

func parseOutputs(out string, input string) ([]int, error) {
	lines := filterNonEmpty(strings.Split(out, "\n"))
	values := make([]int, len(lines))
	for i, line := range lines {
		val, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", line)
		}
		values[i] = val
	}
	return values, nil
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
		buildCase("simple", []int{1, 2, 3}, [][]int{{1, 5}, {2, 4}, {3, 1}}),
		buildCase("reverse", []int{5, 4, 3, 2, 1}, [][]int{{5, 10}, {1, 6}, {3, 7}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(20) + 1
		q := rng.Intn(20) + 1
		arr := randomDistinctArray(rng, n, 1_000_000)
		initial := append([]int(nil), arr...)
		queries := make([][]int, q)
		// maintain set of current values
		current := make(map[int]struct{}, n)
		for _, v := range arr {
			current[v] = struct{}{}
		}
		for j := 0; j < q; j++ {
			pos := rng.Intn(n)
			oldVal := arr[pos]
			delete(current, oldVal)
			newVal := rng.Intn(1_000_000) + 1
			for _, ok := current[newVal]; ok || newVal == oldVal; {
				newVal = rng.Intn(1_000_000) + 1
				_, ok = current[newVal]
			}
			arr[pos] = newVal
			current[newVal] = struct{}{}
			queries[j] = []int{pos + 1, newVal}
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), initial, queries))
	}
	return tests
}

func randomDistinctArray(rng *rand.Rand, n int, maxVal int) []int {
	set := make(map[int]struct{}, n)
	res := make([]int, 0, n)
	for len(res) < n {
		val := rng.Intn(maxVal) + 1
		if _, ok := set[val]; ok {
			continue
		}
		set[val] = struct{}{}
		res = append(res, val)
	}
	return res
}

func buildCase(name string, arr []int, queries [][]int) testCase {
	n := len(arr)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(queries))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	return testCase{name: name, input: sb.String(), q: len(queries)}
}
