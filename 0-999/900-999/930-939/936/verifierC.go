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

const refSource = "0-999/900-999/930-939/936/936C.go"
const maxOps = 6100

type testCase struct {
	name  string
	input string
	n     int
	s     string
	t     string
}

type parsedAnswer struct {
	impossible bool
	k          int
	ops        []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		if err := validateAnswer(tc, refAns); err != nil {
			fmt.Fprintf(os.Stderr, "reference output failed validation on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswer(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if refAns.impossible {
			if !candAns.impossible {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected -1 (impossible) but candidate reported a solution\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, tc.input, refOut, candOut)
				os.Exit(1)
			}
			continue
		}

		if candAns.impossible {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected a sequence of operations but candidate printed -1\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
		if err := validateAnswer(tc, candAns); err != nil {
			fmt.Fprintf(os.Stderr, "candidate answer invalid on test %d (%s): %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-936C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref936C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswer(output string) (parsedAnswer, error) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return parsedAnswer{}, fmt.Errorf("empty output")
	}
	tokens := strings.Fields(trimmed)
	if len(tokens) == 0 {
		return parsedAnswer{}, fmt.Errorf("empty output")
	}
	if tokens[0] == "-1" {
		if len(tokens) > 1 {
			return parsedAnswer{}, fmt.Errorf("extra tokens after -1")
		}
		return parsedAnswer{impossible: true}, nil
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil {
		return parsedAnswer{}, fmt.Errorf("invalid number of operations %q", tokens[0])
	}
	if k < 0 {
		return parsedAnswer{}, fmt.Errorf("negative number of operations %d", k)
	}
	if len(tokens) != k+1 {
		return parsedAnswer{}, fmt.Errorf("expected %d shift values, got %d tokens", k, len(tokens)-1)
	}
	ops := make([]int, k)
	for i := 0; i < k; i++ {
		val, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return parsedAnswer{}, fmt.Errorf("invalid shift value %q", tokens[i+1])
		}
		ops[i] = val
	}
	return parsedAnswer{k: k, ops: ops}, nil
}

func validateAnswer(tc testCase, ans parsedAnswer) error {
	if ans.impossible {
		return nil
	}
	if ans.k != len(ans.ops) {
		return fmt.Errorf("operation count mismatch: declared %d, parsed %d", ans.k, len(ans.ops))
	}
	if ans.k > maxOps {
		return fmt.Errorf("too many operations: got %d (limit %d)", ans.k, maxOps)
	}
	cur := []byte(tc.s)
	for idx, x := range ans.ops {
		if x < 0 || x > tc.n {
			return fmt.Errorf("operation %d: shift value %d out of range [0, %d]", idx+1, x, tc.n)
		}
		cur = applyShift(cur, x)
	}
	if string(cur) != tc.t {
		return fmt.Errorf("final string mismatch: got %q expected %q", string(cur), tc.t)
	}
	return nil
}

func applyShift(p []byte, x int) []byte {
	n := len(p)
	if x == 0 || n == 0 {
		res := make([]byte, n)
		copy(res, p)
		return res
	}
	if x == n {
		res := make([]byte, n)
		for i := 0; i < n; i++ {
			res[i] = p[n-1-i]
		}
		return res
	}
	res := make([]byte, n)
	for i := 0; i < x; i++ {
		res[i] = p[n-1-i]
	}
	copy(res[x:], p[:n-x])
	return res
}

func buildTests() []testCase {
	tests := []testCase{
		makeManualTest("sample1", 6, "abacbb", "babcba"),
		makeManualTest("sample2_impossible", 3, "aba", "bba"),
		makeManualTest("already_equal", 4, "zzzz", "zzzz"),
		makeManualTest("single_char", 1, "a", "a"),
		makeManualTest("mismatch_counts", 5, "abcde", "fffff"),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManualTest(name string, n int, s, t string) testCase {
	return testCase{
		name:  name,
		n:     n,
		s:     s,
		t:     t,
		input: formatInput(n, s, t),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(20) + 1
	s := randomString(rng, n)
	var t string
	if rng.Intn(100) < 70 {
		t = shuffleString(rng, s)
	} else {
		t = randomString(rng, n)
	}
	name := fmt.Sprintf("random_%d_n%d", idx+1, n)
	return testCase{
		name:  name,
		n:     n,
		s:     s,
		t:     t,
		input: formatInput(n, s, t),
	}
}

func formatInput(n int, s, t string) string {
	return fmt.Sprintf("%d\n%s\n%s\n", n, s, t)
}

func randomString(rng *rand.Rand, length int) string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = byte('a' + rng.Intn(26))
	}
	return string(buf)
}

func shuffleString(rng *rand.Rand, s string) string {
	buf := []byte(s)
	for i := len(buf) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
