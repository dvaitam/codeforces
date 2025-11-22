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
	name   string
	n, m   int
	prefs  []string
	slices []int64
}

const refSource = "2070F.go"

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candPath := os.Args[len(os.Args)-1]

	refBin, cleanupRef, err := buildBinary(referencePath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	if err := validateTests(tests); err != nil {
		fmt.Fprintln(os.Stderr, "generated invalid tests:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	if dumpPath := os.Getenv("VERIFIER_DUMP_INPUT"); dumpPath != "" {
		_ = os.WriteFile(dumpPath, []byte(input), 0644)
	}

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refParsed, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		fmt.Fprintln(os.Stderr, previewInput(input))
		os.Exit(1)
	}
	candParsed, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\n%s", err, candOut)
		fmt.Fprintln(os.Stderr, previewInput(input))
		os.Exit(1)
	}

	if len(refParsed) != len(candParsed) {
		fmt.Fprintf(os.Stderr, "mismatched test count: ref %d cand %d\n", len(refParsed), len(candParsed))
		os.Exit(1)
	}

	for i := range tests {
		exp := refParsed[i]
		got := candParsed[i]
		if len(exp) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d (%s) length mismatch: expected %d got %d\n", i+1, tests[i].name, len(exp), len(got))
			fmt.Fprintln(os.Stderr, previewInput(input))
			os.Exit(1)
		}
		for j := range exp {
			if exp[j] != got[j] {
				fmt.Fprintf(os.Stderr, "test %d (%s) answer %d mismatch: expected %d got %d\n", i+1, tests[i].name, j, exp[j], got[j])
				fmt.Fprintln(os.Stderr, previewInput(input))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), refSource)
	}
	return refSource
}

func buildBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "ref2070F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "ref")

	content, err := os.ReadFile(src)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	fixed := string(content)
	pat := "evenMask := (^oddMask) & (fullSize - 1)\n\n\t// Build"
	if strings.Contains(fixed, "evenMask :=") && strings.Contains(fixed, pat) {
		fixed = strings.Replace(fixed, pat, "evenMask := (^oddMask) & (fullSize - 1)\n\t_ = evenMask\n\n\t// Build", 1)
	}
	refFile := filepath.Join(tmpDir, "ref.go")
	if err := os.WriteFile(refFile, []byte(fixed), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}

	cmd := exec.Command("go", "build", "-o", outPath, refFile)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "cand2070F-")
	if err != nil {
		return "", nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	fixed := string(content)
	pat := "evenMask := (^oddMask) & (fullSize - 1)\n\n\t// Build"
	if strings.Contains(fixed, "evenMask :=") && strings.Contains(fixed, pat) {
		fixed = strings.Replace(fixed, pat, "evenMask := (^oddMask) & (fullSize - 1)\n\t_ = evenMask\n\n\t// Build", 1)
	}
	candSrc := filepath.Join(tmpDir, "candidate.go")
	if err := os.WriteFile(candSrc, []byte(fixed), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}

	outPath := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", outPath, candSrc)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutputs(out string, tests []testCase) ([][]int64, error) {
	tokens := strings.Fields(out)
	res := make([][]int64, 0, len(tests))
	idx := 0
	for ti, tc := range tests {
		need := int(sumSlices(tc.slices) + 1)
		if idx+need > len(tokens) {
			return nil, fmt.Errorf("not enough tokens for test %d", ti+1)
		}
		cur := make([]int64, need)
		for j := 0; j < need; j++ {
			val, err := strconv.ParseInt(tokens[idx+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer on test %d position %d: %v", ti+1, j, err)
			}
			cur[j] = val
		}
		idx += need
		res = append(res, cur)
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens at end of output")
	}
	return res, nil
}

func sumSlices(arr []int64) int64 {
	var s int64
	for _, v := range arr {
		s += v
	}
	return s
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 120)
	tests = append(tests, deterministicTests()...)
	tests = append(tests, randomTests(100)...)
	return tests
}

func validateTests(tests []testCase) error {
	for idx, tc := range tests {
		if len(tc.prefs) != tc.m {
			return fmt.Errorf("test %d (%s): prefs len %d != m %d", idx+1, tc.name, len(tc.prefs), tc.m)
		}
		if len(tc.slices) != tc.n {
			return fmt.Errorf("test %d (%s): slices len %d != n %d", idx+1, tc.name, len(tc.slices), tc.n)
		}
		for fi, s := range tc.prefs {
			if len(s) == 0 {
				return fmt.Errorf("test %d (%s): empty preference for friend %d", idx+1, tc.name, fi+1)
			}
			last := byte(0)
			seen := make(map[byte]bool)
			for i := 0; i < len(s); i++ {
				ch := s[i]
				if ch < 'A' || int(ch-'A') >= tc.n {
					return fmt.Errorf("test %d (%s): invalid char %q in friend %d", idx+1, tc.name, ch, fi+1)
				}
				if i > 0 && ch < last {
					return fmt.Errorf("test %d (%s): unsorted pref for friend %d", idx+1, tc.name, fi+1)
				}
				if seen[ch] {
					return fmt.Errorf("test %d (%s): duplicate char in friend %d", idx+1, tc.name, fi+1)
				}
				seen[ch] = true
				last = ch
			}
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name:   "single-pizza",
			n:      1,
			m:      2,
			prefs:  []string{"A", "A"},
			slices: []int64{3},
		},
		{
			name:   "small",
			n:      3,
			m:      4,
			prefs:  []string{"A", "AB", "BC", "AC"},
			slices: []int64{2, 3, 4},
		},
		{
			name:   "all-like",
			n:      2,
			m:      3,
			prefs:  []string{"AB", "AB", "AB"},
			slices: []int64{5, 6},
		},
		{
			name:   "disjoint",
			n:      4,
			m:      4,
			prefs:  []string{"AB", "CD", "A", "BCD"},
			slices: []int64{1, 2, 3, 4},
		},
	}
}

func randomTests(count int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		n := rng.Intn(16) + 1
		m := rng.Intn(100) + 2
		prefs := make([]string, m)
		for i := 0; i < m; i++ {
			mask := rng.Intn((1<<n)-1) + 1 // ensure non-empty
			prefs[i] = maskToString(mask, n)
		}
		slices := make([]int64, n)
		for i := 0; i < n; i++ {
			slices[i] = int64(rng.Intn(20) + 1)
		}
		tests = append(tests, testCase{
			name:   fmt.Sprintf("rnd-%d", len(tests)+1),
			n:      n,
			m:      m,
			prefs:  prefs,
			slices: slices,
		})
	}
	return tests
}

func maskToString(mask int, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 {
			sb.WriteByte(byte('A' + i))
		}
	}
	return sb.String()
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 256)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.m))
		sb.WriteByte('\n')
		for i, s := range tc.prefs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(s)
		}
		sb.WriteByte('\n')
		for i, v := range tc.slices {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func previewInput(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 25 {
		lines = lines[:25]
	}
	return strings.Join(lines, "\n")
}
