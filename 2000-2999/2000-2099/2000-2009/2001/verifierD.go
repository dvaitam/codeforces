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

type caseMeta struct {
	n   int
	arr []int
}

type testCase struct {
	name  string
	input string
	meta  []caseMeta
}

const referenceSource = "./2001D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests, err := buildTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build tests:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		refOutRaw, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOutRaw, tc.meta)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOutRaw)
			os.Exit(1)
		}

		candOutRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOutRaw)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOutRaw, tc.meta)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOutRaw)
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d cases got %d\ninput:\n%sreference:\n%s\ncandidate:\n%s",
				idx+1, tc.name, len(refAns), len(candAns), tc.input, refOutRaw, candOutRaw)
			os.Exit(1)
		}
		for i := range refAns {
			if !equalSlices(refAns[i], candAns[i]) {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, tc.input, refOutRaw, candOutRaw)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2001D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2001D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() ([]testCase, error) {
	var tests []testCase
	manualInputs := []struct {
		name  string
		input string
	}{
		{
			name:  "sample-small",
			input: "3\n3\n3 2 1\n4\n1 1 1 1\n5\n2 1 2 1 2\n\n",
		},
		{
			name:  "single-element",
			input: "2\n1\n1\n1\n1\n\n",
		},
	}
	for _, item := range manualInputs {
		meta, err := parseInputMeta(item.input)
		if err != nil {
			return nil, fmt.Errorf("failed to parse manual test %s: %v", item.name, err)
		}
		tests = append(tests, testCase{name: item.name, input: item.input, meta: meta})
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc, err := randomTest(rng, i)
		if err != nil {
			return nil, err
		}
		tests = append(tests, tc)
	}
	return tests, nil
}

func parseInputMeta(input string) ([]caseMeta, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	metas := make([]caseMeta, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("failed to read n of case %d: %v", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &arr[j]); err != nil {
				return nil, fmt.Errorf("failed to read a[%d] of case %d: %v", j+1, i+1, err)
			}
		}
		metas[i] = caseMeta{n: n, arr: arr}
	}
	return metas, nil
}

func randomTest(rng *rand.Rand, idx int) (testCase, error) {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	metas := make([]caseMeta, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			val := rng.Intn(n) + 1
			arr[j] = val
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		metas[i] = caseMeta{n: n, arr: append([]int(nil), arr...)}
	}
	sb.WriteByte('\n')
	return testCase{name: fmt.Sprintf("random-%d", idx+1), input: sb.String(), meta: metas}, nil
}

func parseAnswers(out string, metas []caseMeta) ([][]int, error) {
	tokens := strings.Fields(out)
	idx := 0
	answers := make([][]int, len(metas))
	for caseIdx, meta := range metas {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("missing length for case %d", caseIdx+1)
		}
		length, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid length token %q for case %d", tokens[idx], caseIdx+1)
		}
		idx++
		if length < 1 || length > meta.n {
			return nil, fmt.Errorf("invalid length %d for case %d", length, caseIdx+1)
		}
		if idx+length > len(tokens) {
			return nil, fmt.Errorf("not enough numbers for case %d", caseIdx+1)
		}
		seq := make([]int, length)
		seen := make(map[int]bool, length)
		pos := 0
		for i := 0; i < length; i++ {
			val, err := strconv.Atoi(tokens[idx+i])
			if err != nil {
				return nil, fmt.Errorf("invalid value %q in case %d", tokens[idx+i], caseIdx+1)
			}
			if val < 1 || val > meta.n {
				return nil, fmt.Errorf("value %d out of range in case %d", val, caseIdx+1)
			}
			if seen[val] {
				return nil, fmt.Errorf("duplicate value %d in case %d", val, caseIdx+1)
			}
			for pos < meta.n && meta.arr[pos] != val {
				pos++
			}
			if pos == meta.n {
				return nil, fmt.Errorf("value %d not subsequence of input in case %d", val, caseIdx+1)
			}
			pos++
			seq[i] = val
			seen[val] = true
		}
		idx += length
		answers[caseIdx] = seq
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens at the end")
	}
	return answers, nil
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
