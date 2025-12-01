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

const refSource = "./2046F1.go"

type testCase struct {
	input string
	cases []string
}

type pair struct {
	ch  byte
	pos int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}

	tests := generateTests()

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]

	for ti, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}
		expect, err := parseReference(refOut, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", ti+1, err, tc.input, got)
			os.Exit(1)
		}

		if err := validateCandidate(got, tc.cases, expect); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%soutput:\n%s\n", ti+1, err, tc.input, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2046F1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseReference(output string, cases []string) ([]bool, error) {
	tokens := strings.Fields(output)
	idx := 0
	expect := make([]bool, len(cases))
	for i, s := range cases {
		n := len(s)
		if idx >= len(tokens) {
			return nil, fmt.Errorf("missing verdict for case %d", i+1)
		}
		verdict := strings.ToUpper(tokens[idx])
		idx++
		if verdict != "YES" && verdict != "NO" {
			return nil, fmt.Errorf("invalid verdict %q on case %d", verdict, i+1)
		}
		expect[i] = verdict == "YES"
		if expect[i] {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("missing string for case %d", i+1)
			}
			idx++ // skip string
			need := 2 * n
			if idx+need > len(tokens) {
				return nil, fmt.Errorf("reference operations too short for case %d", i+1)
			}
			idx += need
		}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("reference output has extra tokens")
	}
	return expect, nil
}

func validateCandidate(output string, cases []string, expect []bool) error {
	tokens := strings.Fields(output)
	idx := 0
	for i, tmpl := range cases {
		n := len(tmpl)
		if idx >= len(tokens) {
			return fmt.Errorf("missing verdict for case %d", i+1)
		}
		verdict := strings.ToUpper(tokens[idx])
		idx++
		if verdict != "YES" && verdict != "NO" {
			return fmt.Errorf("invalid verdict %q on case %d", verdict, i+1)
		}
		ansYes := verdict == "YES"
		if ansYes != expect[i] {
			if expect[i] {
				return fmt.Errorf("case %d should be YES", i+1)
			}
			return fmt.Errorf("case %d should be NO", i+1)
		}
		if !ansYes {
			continue
		}
		if idx >= len(tokens) {
			return fmt.Errorf("missing cuneiform for case %d", i+1)
		}
		cuneiform := tokens[idx]
		idx++
		if cuneiform != tmpl {
			return fmt.Errorf("case %d: expected string %s, got %s", i+1, tmpl, cuneiform)
		}
		totalOps := n
		ops := make([]pair, 0, totalOps)
		for k := 0; k < totalOps; k++ {
			if idx >= len(tokens) {
				return fmt.Errorf("case %d: missing letter for operation %d", i+1, k+1)
			}
			letter := tokens[idx]
			idx++
			if len(letter) != 1 || !strings.Contains("YDX", letter) {
				return fmt.Errorf("case %d: invalid letter %q", i+1, letter)
			}
			if idx >= len(tokens) {
				return fmt.Errorf("case %d: missing position for operation %d", i+1, k+1)
			}
			posStr := tokens[idx]
			idx++
			pos, err := strconv.Atoi(posStr)
			if err != nil || pos < 0 {
				return fmt.Errorf("case %d: invalid position %q", i+1, posStr)
			}
			ops = append(ops, pair{ch: letter[0], pos: pos})
		}
		if err := simulateOperations(ops, tmpl); err != nil {
			return fmt.Errorf("case %d: %v", i+1, err)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("extra tokens in output")
	}
	return nil
}

func simulateOperations(ops []pair, target string) error {
	cur := make([]byte, 0, len(target))
	counts := [3]int{}
	step := 0
	for i, op := range ops {
		if op.pos > len(cur) {
			return fmt.Errorf("operation %d: position %d out of range %d", i+1, op.pos, len(cur))
		}
		cur = append(cur, 0)
		copy(cur[op.pos+1:], cur[op.pos:])
		cur[op.pos] = op.ch
		if op.pos > 0 && cur[op.pos-1] == cur[op.pos] {
			return fmt.Errorf("operation %d: creates equal adjacent letters", i+1)
		}
		if op.pos+1 < len(cur) && cur[op.pos+1] == cur[op.pos] {
			return fmt.Errorf("operation %d: creates equal adjacent letters", i+1)
		}
		switch op.ch {
		case 'Y':
			counts[0]++
		case 'D':
			counts[1]++
		case 'X':
			counts[2]++
		default:
			return fmt.Errorf("operation %d: invalid letter %c", i+1, op.ch)
		}
		step++
		if step == 3 {
			if counts[0] != 1 || counts[1] != 1 || counts[2] != 1 {
				return fmt.Errorf("operations %d-%d must insert each letter once", i-1, i+1)
			}
			counts = [3]int{}
			step = 0
		}
	}
	if step != 0 {
		return fmt.Errorf("incomplete triple of operations")
	}
	if string(cur) != target {
		return fmt.Errorf("operations result in %s, expected %s", string(cur), target)
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20462046))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, makeTest([]string{"YDX", "YDXDYX", "YDXYDX", "YYYXXX"}))
	tests = append(tests, randomCase(rng, rng.Intn(4)+1, 600))
	tests = append(tests, randomCase(rng, rng.Intn(5)+1, 2000))
	for i := 0; i < 20; i++ {
		tests = append(tests, randomCase(rng, rng.Intn(5)+1, 3000))
	}
	tests = append(tests, maxCase())

	return tests
}

func sampleTest() testCase {
	cases := []string{"YDX", "YDXDYX", "YDY", "YDXDYY"}
	return makeTest(cases)
}

func randomCase(rng *rand.Rand, maxCases, maxTotal int) testCase {
	if maxCases < 1 {
		maxCases = 1
	}
	t := rng.Intn(maxCases) + 1
	var stringsList []string
	remaining := maxTotal
	for i := 0; i < t && remaining >= 3; i++ {
		maxLen := min(remaining/3, 200)
		if maxLen == 0 {
			break
		}
		blocks := rng.Intn(maxLen) + 1
		length := blocks * 3
		stringsList = append(stringsList, randomString(rng, length))
		remaining -= length
	}
	if len(stringsList) == 0 {
		stringsList = append(stringsList, "YDX")
	}
	return makeTest(stringsList)
}

func randomString(rng *rand.Rand, length int) string {
	bytes := make([]byte, length)
	chars := []byte{'Y', 'D', 'X'}
	for i := 0; i < length; i++ {
		bytes[i] = chars[rng.Intn(len(chars))]
	}
	return string(bytes)
}

func maxCase() testCase {
	length := 198000
	bytes := make([]byte, length)
	pattern := []byte{'Y', 'D', 'X'}
	for i := 0; i < length; i++ {
		bytes[i] = pattern[i%3]
	}
	return makeTest([]string{string(bytes)})
}

func makeTest(cases []string) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(cases))
	for _, s := range cases {
		fmt.Fprintln(&b, s)
	}
	return testCase{
		input: b.String(),
		cases: cases,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
