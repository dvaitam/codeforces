package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const referenceSolutionRel = "1000-1999/1200-1299/1270-1279/1277/1277D.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "1277D.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	name  string
	words []string
}

func makeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.words)))
		for _, w := range tc.words {
			sb.WriteString(w)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func copyWords(words []string) []string {
	cp := make([]string, len(words))
	copy(cp, words)
	return cp
}

func reverseWord(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func canArrange(words []string) bool {
	if len(words) == 0 {
		return true
	}
	out := [2]int{}
	in := [2]int{}
	used := [2]bool{}
	hasBridge := false
	for _, w := range words {
		s := int(w[0] - '0')
		e := int(w[len(w)-1] - '0')
		out[s]++
		in[e]++
		used[s] = true
		used[e] = true
		if s != e {
			hasBridge = true
		}
	}
	if used[0] && used[1] && !hasBridge {
		return false
	}
	diffPlus := 0
	diffMinus := 0
	for v := 0; v < 2; v++ {
		diff := out[v] - in[v]
		switch {
		case diff == 1:
			diffPlus++
		case diff == -1:
			diffMinus++
		case diff == 0:
		default:
			return false
		}
		if diffPlus > 1 || diffMinus > 1 {
			return false
		}
	}
	return (diffPlus == 0 && diffMinus == 0) || (diffPlus == 1 && diffMinus == 1)
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name: "sample1",
			words: []string{
				"0001",
				"1000",
				"0011",
				"0111",
			},
		},
		{
			name: "sample2",
			words: []string{
				"010",
				"101",
				"0",
			},
		},
		{
			name: "sample3",
			words: []string{
				"00000",
				"00001",
			},
		},
		{
			name: "sample4",
			words: []string{
				"01",
				"001",
				"0001",
				"00001",
			},
		},
		{
			name: "only_loops_zero",
			words: []string{
				"0",
				"00",
				"000",
			},
		},
		{
			name: "only_loops_both",
			words: []string{
				"0",
				"00",
				"1",
				"11",
			},
		},
		{
			name: "needs_balancing",
			words: []string{
				"01",
				"001",
				"100",
				"0101",
			},
		},
		{
			name: "already_balanced",
			words: []string{
				"01",
				"10",
				"101",
				"010",
			},
		},
	}
}

func randomWord(rng *rand.Rand, maxLen int) string {
	length := rng.Intn(maxLen) + 1
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	for i := 0; i < 60; i++ {
		n := rng.Intn(20) + 1
		maxLen := rng.Intn(10) + 1
		set := make(map[string]struct{})
		words := make([]string, 0, n)
		for len(words) < n {
			w := randomWord(rng, maxLen)
			if _, ok := set[w]; ok {
				continue
			}
			set[w] = struct{}{}
			words = append(words, w)
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_small_%d", i+1),
			words: words,
		})
	}
	// A larger stress-like test
	stressWords := make([]string, 0, 200)
	set := make(map[string]struct{})
	for len(stressWords) < 200 {
		w := randomWord(rng, 30)
		if _, ok := set[w]; ok {
			continue
		}
		set[w] = struct{}{}
		stressWords = append(stressWords, w)
	}
	tests = append(tests, testCase{
		name:  "random_stress",
		words: stressWords,
	})
	return tests
}

type parsedAnswer struct {
	possible bool
	k        int
	indexes  []int
}

type tokenScanner struct {
	r *bufio.Reader
}

func newTokenScanner(s string) *tokenScanner {
	return &tokenScanner{r: bufio.NewReader(strings.NewReader(s))}
}

func (ts *tokenScanner) next() (string, error) {
	var sb strings.Builder
	for {
		b, err := ts.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				if sb.Len() == 0 {
					return "", io.EOF
				}
				return sb.String(), nil
			}
			return "", err
		}
		if b > ' ' {
			sb.WriteByte(b)
			break
		}
	}
	for {
		b, err := ts.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if b <= ' ' {
			return sb.String(), nil
		}
		sb.WriteByte(b)
	}
}

func parseAnswers(output string, tests int) ([]parsedAnswer, error) {
	ts := newTokenScanner(output)
	res := make([]parsedAnswer, 0, tests)
	for i := 0; i < tests; i++ {
		token, err := ts.next()
		if err != nil {
			return nil, fmt.Errorf("failed to read answer for test %d: %v", i+1, err)
		}
		if token == "-1" {
			res = append(res, parsedAnswer{possible: false})
			continue
		}
		k, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("invalid number of reversals on test %d: %v", i+1, err)
		}
		ans := parsedAnswer{possible: true, k: k}
		if k > 0 {
			ans.indexes = make([]int, k)
			for j := 0; j < k; j++ {
				token, err = ts.next()
				if err != nil {
					return nil, fmt.Errorf("not enough indexes for test %d", i+1)
				}
				val, err := strconv.Atoi(token)
				if err != nil {
					return nil, fmt.Errorf("invalid index on test %d: %v", i+1, err)
				}
				ans.indexes[j] = val
			}
		}
		res = append(res, ans)
	}
	// Ensure no extra tokens
	if extra, err := ts.next(); err == nil {
		return nil, fmt.Errorf("extra output detected after processing tests: %s", extra)
	}
	return res, nil
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "1277D-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1277D")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validateAnswer(tc testCase, ans parsedAnswer) error {
	n := len(tc.words)
	if ans.k != len(ans.indexes) {
		return fmt.Errorf("reported k=%d but provided %d indexes", ans.k, len(ans.indexes))
	}
	indices := append([]int(nil), ans.indexes...)
	sort.Ints(indices)
	for i := 0; i < len(indices); i++ {
		if indices[i] < 1 || indices[i] > n {
			return fmt.Errorf("index %d out of range", indices[i])
		}
		if i > 0 && indices[i] == indices[i-1] {
			return fmt.Errorf("duplicate index %d", indices[i])
		}
	}
	final := copyWords(tc.words)
	for _, idx := range ans.indexes {
		final[idx-1] = reverseWord(final[idx-1])
	}
	seen := make(map[string]struct{}, n)
	for _, w := range final {
		if _, ok := seen[w]; ok {
			return fmt.Errorf("word %q appears multiple times after reversals", w)
		}
		seen[w] = struct{}{}
	}
	if !canArrange(final) {
		return fmt.Errorf("final set cannot be ordered according to rules")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := makeInput(tests)

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAnswers, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	userAnswers, err := parseAnswers(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		refAns := refAnswers[i]
		userAns := userAnswers[i]
		if !refAns.possible {
			if userAns.possible {
				fmt.Fprintf(os.Stderr, "test %s (%d): expected -1, but participant provided a solution\n", tc.name, i+1)
				os.Exit(1)
			}
			continue
		}
		if !userAns.possible {
			fmt.Fprintf(os.Stderr, "test %s (%d): expected a valid solution with k=%d, but participant output -1\n", tc.name, i+1, refAns.k)
			os.Exit(1)
		}
		if userAns.k != refAns.k {
			fmt.Fprintf(os.Stderr, "test %s (%d): minimal reversals is %d, participant reported %d\n", tc.name, i+1, refAns.k, userAns.k)
			os.Exit(1)
		}
		if err := validateAnswer(tc, userAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %s (%d) invalid solution: %v\n", tc.name, i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
