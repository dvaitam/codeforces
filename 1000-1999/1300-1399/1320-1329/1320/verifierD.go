package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type testCase struct {
	name    string
	input   string
	answers int
}

type query struct {
	l1, l2, length int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/solution")
		os.Exit(1)
	}
	candidate := os.Args[1]
	refSrc := referencePath()
	refBin, cleanup, err := buildReferenceBinary(refSrc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("reference failed on case %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}
		exp, err := normalizeAnswers(refOut, tc.answers)
		if err != nil {
			fmt.Printf("reference produced invalid output on case %d (%s): %v\nraw output:\n%s\n", i+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Printf("case %d (%s): runtime error: %v\n", i+1, tc.name, err)
			fmt.Println(previewInput(tc.input))
			os.Exit(1)
		}
		got, err := normalizeAnswers(candOut, tc.answers)
		if err != nil {
			fmt.Printf("case %d (%s): invalid output: %v\nraw output:\n%s\n", i+1, tc.name, err, candOut)
			os.Exit(1)
		}
		for j := 0; j < tc.answers; j++ {
			if got[j] != exp[j] {
				fmt.Printf("case %d (%s) failed at query %d: expected %s got %s\n", i+1, tc.name, j+1, exp[j], got[j])
				fmt.Println(previewInput(tc.input))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func referencePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "1000-1999/1300-1399/1320-1329/1320/1320D.go"
	}
	return filepath.Join(filepath.Dir(file), "1320D.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier1320D")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref1320d")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference solution: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
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
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, out.String())
	}
	return out.String(), nil
}

func normalizeAnswers(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(tokens))
	}
	res := make([]string, expected)
	for i, tok := range tokens {
		up := strings.ToUpper(tok)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid answer %q", tok)
		}
		res[i] = up
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{}
	tests = append(tests, sampleTest())
	tests = append(tests, uniformTest("all ones", "1111111", []query{{1, 1, 7}, {2, 3, 3}, {1, 5, 3}}))
	tests = append(tests, uniformTest("all zeros", "0000000", []query{{1, 1, 7}, {1, 4, 4}, {2, 5, 2}}))
	tests = append(tests, customStringTest("alternating", "0101010101", []query{
		{1, 2, 9}, {1, 1, 10}, {2, 5, 3}, {3, 4, 4}, {1, 6, 2},
	}))
	tests = append(tests, customStringTest("blocks", "111000111000111", []query{
		{1, 4, 6}, {4, 7, 5}, {7, 10, 4}, {2, 9, 5}, {5, 9, 3},
	}))

	rng := rand.New(rand.NewSource(0x7a12c9d3))
	for i := 0; i < 8; i++ {
		n := 50 + rng.Intn(100)
		q := 60 + rng.Intn(60)
		tests = append(tests, randomTest(fmt.Sprintf("random-mid-%d", i+1), n, q, rng))
	}
	tests = append(tests, randomTest("random-compact", 40, 120, rng))
	tests = append(tests, randomTest("random-large", 2000, 1500, rng))
	tests = append(tests, randomTest("random-wide", 8000, 4000, rng))

	return tests
}

func sampleTest() testCase {
	queries := []query{
		{1, 3, 3},
		{1, 4, 2},
		{1, 2, 3},
	}
	return customStringTest("sample", "11011", queries)
}

func uniformTest(name, s string, queries []query) testCase {
	return customStringTest(name, s, queries)
}

func customStringTest(name, s string, queries []query) testCase {
	return testCase{
		name:    name,
		input:   formatInput(len(s), s, queries),
		answers: len(queries),
	}
}

func randomTest(name string, n, q int, rng *rand.Rand) testCase {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			if i > 0 {
				s[i] = s[i-1]
			} else {
				s[i] = byte('0' + rng.Intn(2))
			}
			continue
		}
		s[i] = byte('0' + rng.Intn(2))
	}
	queries := make([]query, 0, q)
	if q > 0 {
		queries = append(queries, query{1, 1, n})
	}
	if q > 1 {
		mid := n / 2
		if mid == 0 {
			mid = 1
		}
		maxLen := min(n-mid+1, n-mid+1)
		if maxLen < 1 {
			maxLen = 1
		}
		queries = append(queries, query{mid, mid, maxLen})
	}
	for len(queries) < q {
		l1 := rng.Intn(n) + 1
		l2 := rng.Intn(n) + 1
		maxLen := min(n-l1+1, n-l2+1)
		if maxLen <= 0 {
			continue
		}
		length := rng.Intn(maxLen) + 1
		queries = append(queries, query{l1, l2, length})
	}
	return testCase{
		name:    name,
		input:   formatInput(n, string(s), queries),
		answers: len(queries),
	}
}

func formatInput(n int, s string, queries []query) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	sb.WriteString(s)
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", q.l1, q.l2, q.length))
	}
	return sb.String()
}

func previewInput(in string) string {
	const limit = 400
	if len(in) <= limit {
		return "Input:\n" + in
	}
	return fmt.Sprintf("Input (first %d chars):\n%s...\n", limit, in[:limit])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
