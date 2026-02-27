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

// solveCase is the embedded reference implementation.
// For each query position a, the answer is:
//   - a is left of all teachers:  b[0] - 1
//   - a is right of all teachers: n - b[m-1]
//   - a is between two teachers:  (b[right] - b[left]) / 2
func solveCase(n int64, teachers []int64, queries []int64) []int64 {
	b := make([]int64, len(teachers))
	copy(b, teachers)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	m := len(b)
	results := make([]int64, len(queries))
	for qi, a := range queries {
		idx := sort.Search(m, func(i int) bool { return b[i] >= a })
		var ans int64
		if idx == 0 {
			ans = b[0] - 1
		} else if idx == m {
			ans = n - b[m-1]
		} else {
			ans = (b[idx] - b[idx-1]) / 2
		}
		results[qi] = ans
	}
	return results
}

type testInput struct {
	text     string
	expected []int64
}

func buildTest(n int64, teachers []int64, queries []int64) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, len(teachers), len(queries))
	for i, v := range teachers {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return testInput{text: sb.String(), expected: solveCase(n, teachers, queries)}
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(out string, count int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != count {
		return nil, fmt.Errorf("expected %d integers, got %d", count, len(fields))
	}
	ans := make([]int64, count)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		ans[i] = val
	}
	return ans, nil
}

func fixedTests() []testInput {
	return []testInput{
		buildTest(8, []int64{6}, []int64{1}),
		buildTest(10, []int64{1, 4, 8}, []int64{2, 3, 10}),
	}
}

func edgeTests() []testInput {
	return []testInput{
		buildTest(1_000_000_000, []int64{1, 500_000_000, 1_000_000_000}, []int64{2, 999_999_999, 500_000_001}),
	}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	var allExpected []int64
	for range t {
		n := int64(rng.Intn(999_999_998) + 3)
		m := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)

		teacherSet := make(map[int64]struct{})
		for len(teacherSet) < m {
			teacherSet[int64(rng.Intn(int(n)))+1] = struct{}{}
		}
		teachers := make([]int64, 0, m)
		for pos := range teacherSet {
			teachers = append(teachers, pos)
		}
		for i, pos := range teachers {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(pos, 10))
		}
		sb.WriteByte('\n')

		queries := make([]int64, 0, q)
		for len(queries) < q {
			pos := int64(rng.Intn(int(n))) + 1
			if _, exists := teacherSet[pos]; !exists {
				queries = append(queries, pos)
			}
		}
		for i, pos := range queries {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(pos, 10))
		}
		sb.WriteByte('\n')

		allExpected = append(allExpected, solveCase(n, teachers, queries)...)
	}
	return testInput{text: sb.String(), expected: allExpected}
}

func generateTests(rng *rand.Rand) []testInput {
	tests := fixedTests()
	tests = append(tests, edgeTests()...)
	for len(tests) < 60 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func preview(s string) string {
	if len(s) <= 400 {
		return s
	}
	return s[:400] + "...\n"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests(rng)
	for idx, tc := range tests {
		gotRaw, err := runBinary(bin, tc.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, preview(tc.text))
			os.Exit(1)
		}
		got, err := parseOutput(gotRaw, len(tc.expected))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, gotRaw)
			os.Exit(1)
		}
		for i := range tc.expected {
			if tc.expected[i] != got[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d answer %d: expected %d got %d\ninput:\n%s\n", idx+1, i+1, tc.expected[i], got[i], preview(tc.text))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
