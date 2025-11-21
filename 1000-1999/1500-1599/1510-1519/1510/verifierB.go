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

type testCase struct {
	input string
	d     int
	masks []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()

	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refSeq, err := parseSequence(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotSeq, err := parseSequence(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if len(gotSeq) < len(refSeq) {
			fmt.Fprintf(os.Stderr, "test %d invalid: participant sequence shorter than reference (expected %d got %d)\ninput:\n%s\n",
				idx+1, len(refSeq), len(gotSeq), tc.input)
			os.Exit(1)
		}
		if err := simulateSequence(tc.d, tc.masks, gotSeq); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
				idx+1, err, tc.input, refOut, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "1510B_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "1510B.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseSequence(out string) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("output should contain two lines")
	}
	var expected int
	if _, err := fmt.Sscan(lines[0], &expected); err != nil {
		return nil, fmt.Errorf("failed to parse sequence length: %v", err)
	}
	tokens := strings.Fields(lines[1])
	if len(tokens) != expected {
		return nil, fmt.Errorf("declared length %d but got %d tokens", expected, len(tokens))
	}
	return tokens, nil
}

func simulateSequence(d int, masks []int, seq []string) error {
	visited := make(map[int]bool)
	for _, mask := range masks {
		visited[mask] = false
	}
	cur := 0
	for _, token := range seq {
		if token == "R" {
			cur = 0
			continue
		}
		digit, err := strconv.Atoi(token)
		if err != nil {
			return fmt.Errorf("invalid token %q", token)
		}
		if digit < 0 || digit >= d {
			return fmt.Errorf("digit %d out of range", digit)
		}
		cur |= 1 << digit
		if _, exists := visited[cur]; exists {
			visited[cur] = true
		}
	}
	for mask, done := range visited {
		if !done {
			return fmt.Errorf("mask %b was never visited", mask)
		}
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 80)...)
	tests = append(tests, exhaustiveSmall()...)
	return tests
}

func manualTests() []testCase {
	return []testCase{
		makeTestCase(2, []int{0b10, 0b11}),
		makeTestCase(3, []int{0b001, 0b111, 0b101, 0b011}),
	}
}

func randomTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		d := rng.Intn(5) + 1
		limit := 1<<d - 1
		n := rng.Intn(limit) + 1
		maskSet := make(map[int]struct{})
		for len(maskSet) < n {
			mask := rng.Intn(limit) + 1
			maskSet[mask] = struct{}{}
		}
		masks := make([]int, 0, n)
		for mask := range maskSet {
			masks = append(masks, mask)
		}
		tests = append(tests, makeTestCase(d, masks))
	}
	return tests
}

func exhaustiveSmall() []testCase {
	var tests []testCase
	for d := 1; d <= 3; d++ {
		limit := 1<<d - 1
		for mask := 1; mask <= limit; mask++ {
			masks := []int{mask}
			tests = append(tests, makeTestCase(d, masks))
		}
	}
	return tests
}

func makeTestCase(d int, masks []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", d, len(masks)))
	for _, mask := range masks {
		line := make([]byte, d)
		for j := 0; j < d; j++ {
			if mask&(1<<j) != 0 {
				line[j] = '1'
			} else {
				line[j] = '0'
			}
		}
		sb.Write(line)
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), d: d, masks: masks}
}
