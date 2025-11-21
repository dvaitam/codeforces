package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2120-2129/2124/2124A.go"

type testCase struct {
	n   int
	arr []int
}

type testSuite struct {
	name  string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	suites := buildTests()

	for idx, suite := range suites {
		input := buildInput(suite.cases)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, refOut)
			os.Exit(1)
		}
		expectYes, err := parseYesNo(refOut, len(suite.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, candOut)
			os.Exit(1)
		}
		if err := validateCandidate(candOut, suite.cases, expectYes); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d test suites passed\n", len(suites))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2124A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2124A.bin")
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

func parseYesNo(out string, expected int) ([]bool, error) {
	tokens := strings.Fields(out)
	res := make([]bool, expected)
	pos := 0
	for i := 0; i < expected; i++ {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("expected more output for case %d", i+1)
		}
		switch strings.ToUpper(tokens[pos]) {
		case "YES":
			res[i] = true
		case "NO":
			res[i] = false
		default:
			return nil, fmt.Errorf("token %q is not YES/NO", tokens[pos])
		}
		pos++
		if res[i] {
			if pos >= len(tokens) {
				return nil, fmt.Errorf("case %d: missing k after YES", i+1)
			}
			k, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: invalid k %q", i+1, tokens[pos])
			}
			pos++
			pos += k
			if pos > len(tokens) {
				return nil, fmt.Errorf("case %d: not enough elements for k=%d", i+1, k)
			}
		}
	}
	return res, nil
}

func validateCandidate(out string, cases []testCase, expectYes []bool) error {
	tokens := strings.Fields(out)
	pos := 0
	nextToken := func() (string, error) {
		if pos >= len(tokens) {
			return "", fmt.Errorf("unexpected end of output at token %d", pos+1)
		}
		t := tokens[pos]
		pos++
		return t, nil
	}

	for i, tc := range cases {
		tok, err := nextToken()
		if err != nil {
			return fmt.Errorf("case %d: %v", i+1, err)
		}
		ans := strings.ToUpper(tok)
		if ans != "YES" && ans != "NO" {
			return fmt.Errorf("case %d: expected YES/NO, got %q", i+1, tok)
		}
		if ans == "NO" {
			if expectYes[i] {
				return fmt.Errorf("case %d: expected YES solution exists", i+1)
			}
			continue
		}
		// YES
		if !expectYes[i] {
			return fmt.Errorf("case %d: reported YES but reference says impossible", i+1)
		}
		kTok, err := nextToken()
		if err != nil {
			return fmt.Errorf("case %d: missing k: %v", i+1, err)
		}
		k, err := strconv.Atoi(kTok)
		if err != nil {
			return fmt.Errorf("case %d: invalid k %q", i+1, kTok)
		}
		if k < 1 || k > tc.n {
			return fmt.Errorf("case %d: k out of bounds %d", i+1, k)
		}
		if pos+k > len(tokens) {
			return fmt.Errorf("case %d: insufficient elements for k=%d", i+1, k)
		}
		seq := make([]int, k)
		for j := 0; j < k; j++ {
			val, err := strconv.Atoi(tokens[pos+j])
			if err != nil {
				return fmt.Errorf("case %d: invalid element %q", i+1, tokens[pos+j])
			}
			seq[j] = val
		}
		pos += k

		if !isSubsequence(tc.arr, seq) {
			return fmt.Errorf("case %d: provided sequence is not a subsequence of input", i+1)
		}
		if !isDerangement(seq) {
			return fmt.Errorf("case %d: provided sequence is not a derangement", i+1)
		}
	}
	if pos < len(tokens) {
		return fmt.Errorf("extra tokens after all cases starting at %q", tokens[pos])
	}
	return nil
}

func isSubsequence(a, seq []int) bool {
	if len(seq) == 0 {
		return false
	}
	p := 0
	for _, v := range a {
		if v == seq[p] {
			p++
			if p == len(seq) {
				return true
			}
		}
	}
	return false
}

func isDerangement(b []int) bool {
	c := make([]int, len(b))
	copy(c, b)
	sort.Ints(c)
	for i := range b {
		if b[i] == c[i] {
			return false
		}
	}
	return true
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			sb.WriteString(strconv.Itoa(v))
			if i+1 < len(tc.arr) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testSuite {
	sample := testSuite{
		name: "samples",
		cases: []testCase{
			{n: 3, arr: []int{2, 2, 3}},
			{n: 5, arr: []int{4, 5, 5, 2, 4}},
			{n: 3, arr: []int{1, 1, 1}},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomCases []testCase
	for i := 0; i < 10; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(n) + 1
		}
		randomCases = append(randomCases, testCase{n: n, arr: arr})
	}
	randomSuite := testSuite{name: "random", cases: randomCases}

	return []testSuite{sample, randomSuite}
}
