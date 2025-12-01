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

const (
	refSource     = "./1775A2.go"
	totalTests    = 100
	maxLenPerCase = 200000
)

type testCase struct {
	name  string
	input string
}

type answer struct {
	ok      bool
	a, b, c string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		answerCount := countAnswers(tc.input)
		refAns, err := parseOutputs(refOut, answerCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutputs(candOut, answerCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		for caseIdx, cAns := range candAns {
			rAns := refAns[caseIdx]
			if cAns.ok != rAns.ok {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: expected ok=%v, got %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, caseIdx+1, rAns.ok, cAns.ok, tc.input, refOut, candOut)
				os.Exit(1)
			}
			if !cAns.ok {
				continue
			}
			if err := validateAnswer(cAns, tc.input, caseIdx); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, caseIdx+1, err, tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref1775A2-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1775A2.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func countAnswers(input string) int {
	fields := strings.Fields(input)
	if len(fields) > 0 {
		t, err := strconv.Atoi(fields[0])
		if err == nil {
			return t
		}
	}
	return 0
}

func parseOutputs(out string, t int) ([]answer, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	idx := 0
	res := make([]answer, 0, t)
	for idx < len(lines) {
		line := strings.TrimSpace(lines[idx])
		if line == "" {
			idx++
			continue
		}
		if line == ":(" {
			res = append(res, answer{ok: false})
			idx++
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid answer line: %q", line)
		}
		res = append(res, answer{ok: true, a: parts[0], b: parts[1], c: parts[2]})
		idx++
	}
	if len(res) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(res))
	}
	return res, nil
}

func validateAnswer(ans answer, input string, caseIdx int) error {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if caseIdx+1 >= len(lines) {
		return fmt.Errorf("missing case %d in input", caseIdx+1)
	}
	s := lines[caseIdx+1]
	if ans.a+ans.b+ans.c != s {
		return fmt.Errorf("concatenation mismatch, expected %q, got %q+%q+%q", s, ans.a, ans.b, ans.c)
	}
	if len(ans.a) == 0 || len(ans.b) == 0 || len(ans.c) == 0 {
		return fmt.Errorf("empty segment")
	}
	if !(lexLE(ans.a, ans.b) && lexLE(ans.c, ans.b) || lexLE(ans.b, ans.a) && lexLE(ans.b, ans.c)) {
		return fmt.Errorf("lexicographic condition failed for %q %q %q", ans.a, ans.b, ans.c)
	}
	return nil
}

func lexLE(x, y string) bool {
	if x == y {
		return true
	}
	if len(x) == len(y) {
		return x <= y
	}
	for i := 0; i < len(x) && i < len(y); i++ {
		if x[i] == y[i] {
			continue
		}
		return x[i] < y[i]
	}
	return len(x) < len(y)
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "simple_a_b", input: "2\naba\n0\nbbb\n1\n"},
		{name: "all_a", input: "1\naaa\n"},
		{name: "all_b", input: "1\nbbb\n"},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests {
		tests = append(tests, randomTest(rng, len(tests)+1))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		length := rng.Intn(8) + 3
		if rng.Intn(5) == 0 {
			length = rng.Intn(5000) + 3
		}
		var b strings.Builder
		for j := 0; j < length; j++ {
			if rng.Intn(2) == 0 {
				b.WriteByte('a')
			} else {
				b.WriteByte('b')
			}
		}
		sb.WriteString(b.String())
		sb.WriteByte('\n')
	}
	return testCase{name: fmt.Sprintf("rand_%d", idx), input: sb.String()}
}
