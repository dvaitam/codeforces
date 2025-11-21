package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "1000-1999/1500-1599/1530-1539/1532/1532F.go"

type testCase struct {
	n     int
	parts []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}
		refPattern, err := parsePattern(refOut, len(tc.parts))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}
		if err := validateSolution(tc, refPattern); err != nil {
			fmt.Fprintf(os.Stderr, "reference output failed validation on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candPattern, err := parsePattern(candOut, len(tc.parts))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}
		if err := validateSolution(tc, candPattern); err != nil {
			fmt.Fprintf(os.Stderr, "candidate answer rejected on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1532F-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1532F.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
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
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, part := range tc.parts {
		sb.WriteString(part)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parsePattern(output string, expected int) ([]byte, error) {
	var pattern []byte
	for _, ch := range output {
		if ch == 'P' || ch == 'S' {
			pattern = append(pattern, byte(ch))
		} else if ch == '\n' || ch == '\r' || ch == '\t' || ch == ' ' {
			continue
		} else {
			return nil, fmt.Errorf("invalid character %q in output", ch)
		}
	}
	if len(pattern) != expected {
		return nil, fmt.Errorf("expected %d labels, got %d", expected, len(pattern))
	}
	return pattern, nil
}

func validateSolution(tc testCase, pattern []byte) error {
	if len(pattern) != len(tc.parts) {
		return fmt.Errorf("pattern length mismatch")
	}
	lengthIdx := buildLengthIdx(tc)
	if err := checkLengthAssignments(lengthIdx, pattern); err != nil {
		return err
	}
	candidates := possibleStrings(tc)
	if len(candidates) == 0 {
		return fmt.Errorf("failed to derive candidate strings")
	}
	for _, cand := range candidates {
		if matchesPattern(cand, tc, pattern, lengthIdx) {
			return nil
		}
	}
	return fmt.Errorf("no consistent string exists for provided labeling")
}

func buildLengthIdx(tc testCase) [][]int {
	idx := make([][]int, tc.n)
	for i, part := range tc.parts {
		l := len(part)
		if l >= tc.n || l == 0 {
			continue
		}
		idx[l] = append(idx[l], i)
	}
	return idx
}

func checkLengthAssignments(lengthIdx [][]int, pattern []byte) error {
	for l := 1; l < len(lengthIdx); l++ {
		idxs := lengthIdx[l]
		if len(idxs) != 2 {
			return fmt.Errorf("length %d should have exactly 2 strings, got %d", l, len(idxs))
		}
		pCnt, sCnt := 0, 0
		for _, idx := range idxs {
			if idx < 0 || idx >= len(pattern) {
				return fmt.Errorf("invalid index %d for length %d", idx, l)
			}
			switch pattern[idx] {
			case 'P':
				pCnt++
			case 'S':
				sCnt++
			default:
				return fmt.Errorf("invalid character %q in pattern", pattern[idx])
			}
		}
		if pCnt != 1 || sCnt != 1 {
			return fmt.Errorf("length %d requires exactly one prefix and one suffix (got P=%d S=%d)", l, pCnt, sCnt)
		}
	}
	return nil
}

func possibleStrings(tc testCase) []string {
	var longest []string
	for _, part := range tc.parts {
		if len(part) == tc.n-1 {
			longest = append(longest, part)
		}
	}
	if len(longest) != 2 {
		return nil
	}
	makeCandidate := func(a, b string) string {
		if len(b) == 0 {
			return ""
		}
		return a + b[len(b)-1:]
	}
	s1 := makeCandidate(longest[0], longest[1])
	s2 := makeCandidate(longest[1], longest[0])
	var res []string
	if len(s1) == tc.n {
		res = append(res, s1)
	}
	if len(s2) == tc.n {
		if len(res) == 0 || s2 != res[0] {
			res = append(res, s2)
		}
	}
	return res
}

func matchesPattern(cand string, tc testCase, pattern []byte, lengthIdx [][]int) bool {
	if len(cand) != tc.n {
		return false
	}
	for l := 1; l < tc.n; l++ {
		prefix := cand[:l]
		suffix := cand[tc.n-l:]
		for _, idx := range lengthIdx[l] {
			part := tc.parts[idx]
			if pattern[idx] == 'P' {
				if part != prefix {
					return false
				}
			} else {
				if part != suffix {
					return false
				}
			}
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		{
			n: 5,
			parts: []string{
				"ba", "a", "abab", "a", "aba", "baba", "ab", "aba",
			},
		},
		{
			n:     3,
			parts: []string{"a", "aa", "aa", "a"},
		},
		{
			n:     2,
			parts: []string{"a", "c"},
		},
	}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(99) + 2
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte('a' + rng.Intn(26))
	}
	s := string(bytes)

	parts := make([]string, 0, 2*(n-1))
	for l := 1; l < n; l++ {
		parts = append(parts, s[:l])
		parts = append(parts, s[n-l:])
	}
	rng.Shuffle(len(parts), func(i, j int) {
		parts[i], parts[j] = parts[j], parts[i]
	})

	return testCase{
		n:     n,
		parts: parts,
	}
}
